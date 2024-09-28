-- +goose Up
create table if not exists orders (
    order_id bigint primary key, 
    user_id bigint,
    valid_time text not null CHECK (length(valid_time) <= 20),
    order_state text CHECK (length(order_state) <= 20),
    price bigint,
    weight bigint, 
    package text CHECK (length(package) <= 20),
    additional_stretch boolean
);

-- +goose Down
drop table if exists orders;
