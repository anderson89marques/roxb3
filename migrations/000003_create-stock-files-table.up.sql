CREATE TABLE stock_files (
	id Serial primary key,
	name varchar(255) not null unique,
	created_at timestamp default current_timestamp
);

