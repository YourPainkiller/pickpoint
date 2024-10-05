/*Requests without index optimisation*/

explain (analyse, verbose, buffers) select order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch from orders where order_id = 10200

/*
Index Scan using orders_pkey on public.orders  (cost=0.42..8.44 rows=1 width=56) (actual time=0.034..0.034 rows=1 loops=1)
  Output: order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch
  Index Cond: (orders.order_id = 10200)
  Buffers: shared hit=7
Planning:
  Buffers: shared hit=55
Planning Time: 1.282 ms
Execution Time: 0.061 ms
*/

explain (analyse, verbose, buffers) update orders set valid_time = '2024-10-10', order_state = 'accepted' where order_id = 10200

/*
Update on public.orders  (cost=0.42..8.44 rows=0 width=0) (actual time=0.135..0.135 rows=0 loops=1)
  Buffers: shared hit=15
  ->  Index Scan using orders_pkey on public.orders  (cost=0.42..8.44 rows=1 width=70) (actual time=0.040..0.041 rows=1 loops=1)
        Output: '2024-10-10'::text, 'accepted'::text, ctid
        Index Cond: (orders.order_id = 10200)
        Buffers: shared hit=7
Planning:
  Buffers: shared hit=34
Planning Time: 0.644 ms
Execution Time: 0.339 ms
*/

explain (analyse, verbose, buffers) select order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch from orders where user_id = 925 and order_state != 'deleted'

/*
Gather  (cost=1000.00..4968.97 rows=15 width=56) (actual time=1.231..28.804 rows=14 loops=1)
  Output: order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch
  Workers Planned: 1
  Workers Launched: 1
  Buffers: shared hit=2258
  ->  Parallel Seq Scan on public.orders  (cost=0.00..3967.47 rows=9 width=56) (actual time=0.477..8.226 rows=7 loops=2)
        Output: order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch
        Filter: ((orders.order_state <> 'deleted'::text) AND (orders.user_id = 925))
        Rows Removed by Filter: 96868
        Buffers: shared hit=2258
        Worker 0:  actual time=0.001..0.001 rows=0 loops=1
Planning:
  Buffers: shared hit=54
Planning Time: 1.000 ms
Execution Time: 28.843 ms
*/

explain (analyse, verbose, buffers) select order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch from orders where order_state = 'returned'

/*
Seq Scan on public.orders  (cost=0.00..4679.75 rows=48926 width=56) (actual time=0.012..15.430 rows=48229 loops=1)
  Output: order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch
  Filter: (orders.order_state = 'returned'::text)
  Rows Removed by Filter: 145521
  Buffers: shared hit=2258
Planning:
  Buffers: shared hit=54
Planning Time: 1.012 ms
Execution Time: 17.114 ms
*/



/*-----------------------------------------------------------------------------------*/
/*Requests with order_id hash optimisation*/

explain (analyse, verbose, buffers) select order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch from orders where order_id = 10200

/*
Index Scan using orders_order_id_hash_idx on public.orders  (cost=0.00..8.02 rows=1 width=56) (actual time=0.041..0.042 rows=1 loops=1)
  Output: order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch
  Index Cond: (orders.order_id = 10200)
  Buffers: shared hit=6
Planning:
  Buffers: shared hit=73
Planning Time: 1.270 ms
Execution Time: 0.055 ms
*/

explain (analyse, verbose, buffers) update orders set valid_time = '2024-10-10', order_state = 'accepted' where order_id = 10200

/*
Update on public.orders  (cost=0.00..8.02 rows=0 width=0) (actual time=0.141..0.141 rows=0 loops=1)
  Buffers: shared hit=14
  ->  Index Scan using orders_order_id_hash_idx on public.orders  (cost=0.00..8.02 rows=1 width=70) (actual time=0.060..0.061 rows=1 loops=1)
        Output: '2024-10-10'::text, 'accepted'::text, ctid
        Index Cond: (orders.order_id = 10200)
        Buffers: shared hit=6
Planning:
  Buffers: shared hit=52
Planning Time: 1.035 ms
Execution Time: 0.475 ms
*/

explain (analyse, verbose, buffers) select order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch from orders where user_id = 925 and order_state != 'deleted'

/*
    Gather  (cost=1000.00..4969.06 rows=15 width=56) (actual time=1.014..22.168 rows=14 loops=1)
  Output: order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch
  Workers Planned: 1
  Workers Launched: 1
  Buffers: shared hit=2258
  ->  Parallel Seq Scan on public.orders  (cost=0.00..3967.56 rows=9 width=56) (actual time=0.417..7.221 rows=7 loops=2)
        Output: order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch
        Filter: ((orders.order_state <> 'deleted'::text) AND (orders.user_id = 925))
        Rows Removed by Filter: 96868
        Buffers: shared hit=2258
        Worker 0:  actual time=0.001..0.001 rows=0 loops=1
Planning:
  Buffers: shared hit=72
Planning Time: 0.913 ms
Execution Time: 22.196 ms
*/

explain (analyse, verbose, buffers) select order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch from orders where order_state = 'returned'

/*
Seq Scan on public.orders  (cost=0.00..4679.88 rows=48928 width=56) (actual time=0.011..13.482 rows=48229 loops=1)
  Output: order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch
  Filter: (orders.order_state = 'returned'::text)
  Rows Removed by Filter: 145521
  Buffers: shared hit=2258
Planning:
  Buffers: shared hit=72
Planning Time: 1.074 ms
Execution Time: 14.981 ms
*/

/*-----------------------------------------------------------------------------------*/
/*Requests with order_state b-tree optimisation*/

explain (analyse, verbose, buffers) select order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch from orders where order_id = 10200

/*
Index Scan using orders_order_id_hash_idx on public.orders  (cost=0.00..8.02 rows=1 width=56) (actual time=0.052..0.053 rows=1 loops=1)
  Output: order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch
  Index Cond: (orders.order_id = 10200)
  Buffers: shared hit=6
Planning:
  Buffers: shared hit=91
Planning Time: 1.629 ms
Execution Time: 0.071 ms
*/

explain (analyse, verbose, buffers) update orders set valid_time = '2024-10-10', order_state = 'accepted' where order_id = 10200

/*
Update on public.orders  (cost=0.00..8.02 rows=0 width=0) (actual time=0.118..0.118 rows=0 loops=1)
  Buffers: shared hit=14
  ->  Index Scan using orders_order_id_hash_idx on public.orders  (cost=0.00..8.02 rows=1 width=70) (actual time=0.046..0.047 rows=1 loops=1)
        Output: '2024-10-10'::text, 'accepted'::text, ctid
        Index Cond: (orders.order_id = 10200)
        Buffers: shared hit=6
Planning:
  Buffers: shared hit=70
Planning Time: 1.169 ms
Execution Time: 0.381 ms
*/

explain (analyse, verbose, buffers) select order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch from orders where user_id = 925 and order_state != 'deleted'

/*
Gather  (cost=1000.00..4969.06 rows=15 width=56) (actual time=1.043..21.647 rows=14 loops=1)
  Output: order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch
  Workers Planned: 1
  Workers Launched: 1
  Buffers: shared hit=2258
  ->  Parallel Seq Scan on public.orders  (cost=0.00..3967.56 rows=9 width=56) (actual time=0.432..6.906 rows=7 loops=2)
        Output: order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch
        Filter: ((orders.order_state <> 'deleted'::text) AND (orders.user_id = 925))
        Rows Removed by Filter: 96868
        Buffers: shared hit=2258
        Worker 0:  actual time=0.001..0.001 rows=0 loops=1
Planning:
  Buffers: shared hit=92
Planning Time: 0.980 ms
Execution Time: 21.678 ms
*/

explain (analyse, verbose, buffers) select order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch from orders where order_state = 'returned'

/*
Bitmap Heap Scan on public.orders  (cost=551.49..3421.09 rows=48928 width=56) (actual time=1.213..6.589 rows=48229 loops=1)
  Output: order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch
  Recheck Cond: (orders.order_state = 'returned'::text)
  Heap Blocks: exact=2258
  Buffers: shared hit=2301
  ->  Bitmap Index Scan on orders_order_state_btree_idx  (cost=0.00..539.25 rows=48928 width=0) (actual time=0.988..0.988 rows=48229 loops=1)
        Index Cond: (orders.order_state = 'returned'::text)
        Buffers: shared hit=43
Planning:
  Buffers: shared hit=94
Planning Time: 1.423 ms
Execution Time: 7.986 ms
*/