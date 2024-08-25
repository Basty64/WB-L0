-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
                        order_uid VARCHAR(255) PRIMARY KEY,
                        id SERIAL UNIQUE,
                        track_number VARCHAR(255) UNIQUE,
                        entry VARCHAR(255),
                        locale VARCHAR(2),
                        internal_signature VARCHAR(255),
                        customer_id VARCHAR(255),
                        delivery_service VARCHAR(255),
                        shardkey VARCHAR(2),
                        sm_id INT,
                        date_created TIMESTAMP WITH TIME ZONE,
                        oofshard VARCHAR(2)
);

CREATE TABLE deliveries (
                            order_uid VARCHAR(255) REFERENCES orders(order_uid),
                            name VARCHAR(255),
                            phone VARCHAR(20),
                            zip VARCHAR(20),
                            city VARCHAR(255),
                            address VARCHAR(255),
                            region VARCHAR(255),
                            email VARCHAR(255),
                            PRIMARY KEY (order_uid)
);

CREATE TABLE payments (
                          order_uid VARCHAR(255) REFERENCES orders(order_uid),
                          transaction VARCHAR(255),
                          request_id VARCHAR(255),
                          currency VARCHAR(3),
                          provider VARCHAR(255),
                          amount NUMERIC(10,2),
                          payment_dt INT,
                          bank VARCHAR(255),
                          delivery_cost NUMERIC(10,2),
                          goods_total NUMERIC(10,2),
                          custom_fee NUMERIC(10,2),
                          PRIMARY KEY (order_uid)
);

CREATE TABLE items (
                       order_uid VARCHAR(255) REFERENCES orders(order_uid),
                       chrt_id INT,
                       track_number VARCHAR(255) UNIQUE,
                       price NUMERIC(10,2),
                       rid VARCHAR(255),
                       name VARCHAR(255),
                       sale INT,
                       size VARCHAR(20),
                       total_price NUMERIC(10,2),
                       nm_id INT,
                       brand VARCHAR(255),
                       status INT,
                       PRIMARY KEY (order_uid, chrt_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS deliveries CASCADE;
DROP TABLE IF EXISTS payments CASCADE;
DROP TABLE IF EXISTS items CASCADE;

-- +goose StatementEnd