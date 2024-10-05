-- +goose NO TRANSACTION
-- +goose Up
create index concurrently orders_order_state_btree_idx on orders(order_state);
-- +goose Down
drop index concurrently orders_order_state_btree_idx;