CREATE TABLE public.balances
(
    currency       varchar(10) NULL,
    balance        numeric NULL,
    locked_balance numeric NULL,
    updated_date   int8 NULL,
    user_id        varchar(36) null,

    UNIQUE (currency, user_id)
);