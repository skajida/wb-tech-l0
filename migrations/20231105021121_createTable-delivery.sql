CREATE TABLE delivery(
    delivery_id   INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name_receiver VARCHAR(70),
    phone_number  VARCHAR(15),
    zip_code      VARCHAR(10),
    name_city     VARCHAR(35),
    direction     VARCHAR(100),
    region        VARCHAR(20),
    email         VARCHAR(255)
);
