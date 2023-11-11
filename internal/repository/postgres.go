package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"wb-tech-l0/internal/model"

	_ "github.com/lib/pq"
)

type PgDatabase struct {
	database *sql.DB
}

func NewPgDatabase(db *sql.DB) *PgDatabase {
	return &PgDatabase{database: db}
}

func (pgdb PgDatabase) getOrderTx(ctx context.Context, tx *sql.Tx, orderUid model.OrderId) (int, model.Order, error) {
	const query = `
	SELECT buy_id, name_receiver, phone_number, zip_code, name_city, direction, region, email, transaction_id, request_id, name_currency, name_provider, total_amount, payment_dt, name_bank, delivery_cost, total_goods, custom_fee, order_uid, track_number, name_entry, locale, internal_signature, customer_slug, delivery_service, shardkey, sm_id, date_created, oof_shard
	FROM buy
		INNER JOIN delivery USING(delivery_id)
		INNER JOIN payment USING(payment_id)
	WHERE order_uid = $1;
	`

	var buyId int
	var foundOrder model.Order
	err := tx.QueryRowContext(ctx, query, orderUid).Scan(
		&buyId, &foundOrder.Name, &foundOrder.Phone, &foundOrder.Zip, &foundOrder.City, &foundOrder.Address, &foundOrder.Region, &foundOrder.Email,
		&foundOrder.Transaction, &foundOrder.RequestId, &foundOrder.Currency, &foundOrder.Provider, &foundOrder.Amount, &foundOrder.PaymentDt, &foundOrder.Bank, &foundOrder.DeliveryCost, &foundOrder.GoodsTotal, &foundOrder.CustomFee,
		&foundOrder.Uid, &foundOrder.TrackNumber, &foundOrder.Entry, &foundOrder.Locale, &foundOrder.InternalSignature, &foundOrder.CustomerId, &foundOrder.DeliveryService, &foundOrder.Shardkey, &foundOrder.SmId, &foundOrder.DateCreated, &foundOrder.OofShard,
	)
	return buyId, foundOrder, err
}

func (pgdb PgDatabase) getItemsTx(ctx context.Context, tx *sql.Tx, buyId int) ([]model.Item, error) {
	const query = `
	SELECT chart_id, track_number, price, rid, name_item, sale, size_item, total_price, nm_id, brand, status_id
	FROM item
		INNER JOIN buy_item USING(item_id)
	WHERE buy_id = $1;
	`

	rows, err := tx.QueryContext(ctx, query, buyId)
	if err != nil {
		return []model.Item{}, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		rows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status)
		items = append(items, item)
	}
	return items, nil
}

func (pgdb PgDatabase) getOrderInfoUnwrapped(ctx context.Context, orderUid model.OrderId) (model.Order, error) {
	tx, err := pgdb.database.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return model.Order{}, err
	}

	buyId, order, err := pgdb.getOrderTx(ctx, tx, orderUid)
	if err != nil {
		tx.Rollback()
		return model.Order{}, err
	}

	items, err := pgdb.getItemsTx(ctx, tx, buyId)
	if err != nil {
		tx.Rollback()
		return model.Order{}, err
	}

	if err = tx.Commit(); err == nil {
		order.Items = items
	}
	return order, err
}

func (pgdb PgDatabase) GetOrderInfo(ctx context.Context, orderUid model.OrderId) (model.Order, error) {
	order, err := pgdb.getOrderInfoUnwrapped(ctx, orderUid)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return model.Order{}, model.ErrOrderBadParam
		}
	}
	return order, err
}

func (pgdb PgDatabase) addDeliveryTx(ctx context.Context, tx *sql.Tx, delivery model.Delivery) (int, error) {
	const query = `
	INSERT INTO delivery(name_receiver, phone_number, zip_code, name_city, direction, region, email)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING delivery_id;
	`

	var insertedDeliveryId int
	err := tx.QueryRowContext(ctx, query, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email).Scan(&insertedDeliveryId)
	return insertedDeliveryId, err
}

func (pgdb PgDatabase) addPaymentTx(ctx context.Context, tx *sql.Tx, payment model.Payment) (int, error) {
	const query = `
	INSERT INTO payment(transaction_id, request_id, name_currency, name_provider, total_amount, payment_dt, name_bank, delivery_cost, total_goods, custom_fee)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING payment_id;
	`

	var insertedPaymentId int
	err := tx.QueryRowContext(ctx, query, payment.Transaction, payment.RequestId, payment.Currency, payment.Provider, payment.Amount, payment.PaymentDt, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee).Scan(&insertedPaymentId)
	return insertedPaymentId, err
}

func (pgdb PgDatabase) addBuyTx(ctx context.Context, tx *sql.Tx, order model.Order, deliveryId, paymentId int) (int, error) {
	const query = `
	INSERT INTO buy(order_uid, track_number, name_entry, delivery_id, payment_id, locale, internal_signature, customer_slug, delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	RETURNING buy_id;
	`

	var insertedBuyId int
	err := tx.QueryRowContext(ctx, query, order.Uid, order.TrackNumber, order.Entry, deliveryId, paymentId, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.Shardkey, order.SmId, order.DateCreated, order.OofShard).Scan(&insertedBuyId)
	return insertedBuyId, err
}

func generatePlaceholdersTemplate(offset, quantity int) string {
	placeholdersSequence := make([]string, quantity)
	for i := range placeholdersSequence {
		placeholdersSequence[i] = fmt.Sprintf("$%d", offset+i)
	}
	return fmt.Sprintf("(%s)", strings.Join(placeholdersSequence, ", "))
}

func (pgdb PgDatabase) addItemsTx(ctx context.Context, tx *sql.Tx, items []model.Item) ([]int, error) {
	const (
		argsQuantity  = 11
		queryTemplate = `
	INSERT INTO item(chart_id, track_number, price, rid, name_item, sale, size_item, total_price, nm_id, brand, status_id)
	VALUES %s
	RETURNING item_id;
	`
	)

	placeholders := make([]string, 0, len(items))
	queryArguments := make([]any, 0, len(items)*argsQuantity)
	for i, item := range items {
		placeholders = append(placeholders, generatePlaceholdersTemplate(i*argsQuantity+1, argsQuantity))
		queryArguments = append(queryArguments, item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status)
	}

	rows, err := tx.QueryContext(ctx, fmt.Sprintf(queryTemplate, strings.Join(placeholders, ", ")), queryArguments...)
	if err != nil {
		return []int{}, err
	}
	defer rows.Close()

	insertedItemIds := make([]int, len(items))
	for i := 0; rows.Next(); i++ {
		if err = rows.Scan(&insertedItemIds[i]); err != nil {
			return []int{}, err
		}
	}
	return insertedItemIds, nil
}

func (pgdb PgDatabase) addBuyItemsTx(ctx context.Context, tx *sql.Tx, buyId int, itemIds []int) error {
	const (
		argsQuantity  = 2
		queryTemplate = `
	INSERT INTO buy_item(buy_id, item_id)
	VALUES %s;
	`
	)

	placeholders := make([]string, 0, len(itemIds))
	queryArguments := make([]any, 0, len(itemIds)*argsQuantity)
	for i, itemId := range itemIds {
		placeholders = append(placeholders, generatePlaceholdersTemplate(i*argsQuantity+1, argsQuantity))
		queryArguments = append(queryArguments, buyId, itemId)
	}

	_, err := tx.ExecContext(ctx, fmt.Sprintf(queryTemplate, strings.Join(placeholders, ", ")), queryArguments...)
	return err
}

func (pgdb PgDatabase) addOrderUnwrapped(ctx context.Context, order model.Order) error {
	tx, err := pgdb.database.BeginTx(ctx, nil)
	if err != nil { // TODO decompose repeated code ? conveyor pattern
		tx.Rollback()
		return err
	}

	deliveryId, err := pgdb.addDeliveryTx(ctx, tx, order.Delivery)
	if err != nil {
		tx.Rollback()
		return err
	}

	paymentId, err := pgdb.addPaymentTx(ctx, tx, order.Payment)
	if err != nil {
		tx.Rollback()
		return err
	}

	buyId, err := pgdb.addBuyTx(ctx, tx, order, deliveryId, paymentId)
	if err != nil {
		tx.Rollback()
		return err
	}

	itemIds, err := pgdb.addItemsTx(ctx, tx, order.Items)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = pgdb.addBuyItemsTx(ctx, tx, buyId, itemIds); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func (pgdb PgDatabase) AddOrder(ctx context.Context, order model.Order) error {
	// TODO write middleware
	err := pgdb.addOrderUnwrapped(ctx, order)
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		switch {
		case strings.Contains(err.Error(), "_order_uid_key"):
			return model.ErrOrderConflict
		}
	}
	return err
}
