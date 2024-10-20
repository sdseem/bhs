CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY,
    username VARCHAR(255) unique,
    password_hash VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS assets(
    id bigserial PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    price VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS assets_library(
    buyer_id bigint REFERENCES users(id),
    asset_id bigint REFERENCES assets(id),
    PRIMARY KEY (buyer_id, asset_id)
);

insert into assets (name, description, price) select 'test1', 'desc1', '100';
insert into assets (name, description, price) select 'test2', 'desc2', '200';
insert into assets (name, description, price) select 'test3', 'desc3', '300';
insert into assets (name, description, price) select 'test4', 'desc4', '400';
insert into assets (name, description, price) select 'test5', 'desc5', '500';
insert into assets (name, description, price) select 'test6', 'desc6', '600';