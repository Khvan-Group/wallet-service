create table if not exists t_wallets
(
    username varchar(255) primary key not null,
    total    int                      not null default 0
);