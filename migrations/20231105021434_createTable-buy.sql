CREATE TABLE buy(
    buy_id             INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    order_uid          VARCHAR(32),
    track_number       VARCHAR(30),
    name_entry         VARCHAR(15),
    delivery_id        INT NOT NULL,
    payment_id         INT NOT NULL,
    locale             VARCHAR(10),
    internal_signature VARCHAR(32),
    customer_slug      VARCHAR(32) NOT NULL,
    delivery_service   VARCHAR(15),
    shardkey           VARCHAR(10),
    sm_id              INT NOT NULL,
    date_created       TIMESTAMPTZ NOT NULL,
    oof_shard          VARCHAR(10),
    FOREIGN KEY (delivery_id) REFERENCES delivery(delivery_id),
    FOREIGN KEY (payment_id)  REFERENCES payment(payment_id)
);
