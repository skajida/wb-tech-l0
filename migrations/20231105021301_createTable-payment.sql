CREATE TABLE payment(
    payment_id     INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    transaction_id VARCHAR(32) NOT NULL,
    request_id     VARCHAR(32),
    name_currency  VARCHAR(10),
    name_provider  VARCHAR(15) NOT NULL,
    total_amount   INT,
    payment_dt     BIGINT,
    name_bank      VARCHAR(30),
    delivery_cost  INT,
    total_goods    INT,
    custom_fee     INT
);
