CREATE TABLE item(
    item_id      INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    chart_id     BIGINT NOT NULL,
    track_number VARCHAR(30),
    price        INT,
    rid          VARCHAR(32),
    name_item    VARCHAR(50),
    sale         INT,
    size_item    VARCHAR(15),
    total_price  INT,
    nm_id        BIGINT,
    brand        VARCHAR(50),
    status_id    INT NOT NULL
);
