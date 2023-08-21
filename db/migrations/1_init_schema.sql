BEGIN;

CREATE TABLE IF NOT EXISTS purchases
(
    purchase_id varchar(50) NOT NULL,
    purchase_date timestamp NOT NULL,
    purchase_type varchar(50) NULL,
    amount numeric NOT NULL,
    customer_id varchar(50) NULL,
    mall_address varchar(50) NULL,
    mall_id varchar(50) NULL,
    mall_name varchar(50) NULL,
    order_source varchar(50) NULL,
    

    PRIMARY KEY (purchase_id)
);