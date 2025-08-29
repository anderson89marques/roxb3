create table stocks (
	id Serial primary key,
	ticker varchar(20) not null,
	trade_date Date not null,
	trade_price decimal(10,3),
	volume bigint,
	closing_time time not null,
	created_at timestamp default current_timestamp
);

-- Indexes for fast filtering
create index idx_stock_ticker on stocks (ticker);
create index idx_stock_trade_date on stocks (trade_date);
create index idx_stock_date on stocks (ticker, trade_date);
create index idx_stock_ticker_price_volume ON stocks (ticker, trade_price, volume);

