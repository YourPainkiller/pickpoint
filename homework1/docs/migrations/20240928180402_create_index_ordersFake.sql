-- +goose NO TRANSACTION
-- +goose Up
create index concurrently orders_order_id_hash_idx on orders using HASH (order_id);
-- +goose Down
drop index concurrently orders_order_id_hash_idx;
