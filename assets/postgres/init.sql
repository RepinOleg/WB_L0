CREATE TABLE orders
(
    id bigserial PRIMARY KEY,
    order_id text,
    order_data jsonb
);