-- Create a materialized view
CREATE MATERIALIZED VIEW stock_summary as
select ticker, trade_date, sum(volume) as total_volume, max(trade_price) as max_price  from stocks
group by ticker, trade_date
with data;

create index idx_stock_summary_ticker_trade_date ON stock_summary (ticker, trade_date);

