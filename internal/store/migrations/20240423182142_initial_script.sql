-- +migrate Up

CREATE TABLE companies (
    id          uuid not null unique primary key,
    name        varchar not null unique
);

CREATE TABLE users (
    id         uuid not null unique primary key,
    company    uuid references companies (id),
    role       varchar,
    name       varchar,
    surname    varchar,
    phone      numeric,
    email      varchar not null unique,
    user_type   varchar,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);

-- +migrate Down

DROP TABLE users, companies CASCADE;