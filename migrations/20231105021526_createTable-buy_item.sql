CREATE TABLE buy_item(
    buy_item_id   INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    buy_id        INT NOT NULL,
    item_id       INT NOT NULL,
    FOREIGN KEY (buy_id) REFERENCES buy(buy_id),
    FOREIGN KEY (item_id)  REFERENCES item(item_id)
);
