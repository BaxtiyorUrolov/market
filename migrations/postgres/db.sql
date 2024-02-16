CREATE TYPE payment_type_enum AS ENUM ('card', 'cash');
CREATE TYPE status_enum AS ENUM ('in_process', 'success', 'cancel');
CREATE TYPE transaction_type_enum AS ENUM ('withdraw', 'topup');
CREATE TYPE source_type_enum AS ENUM ('bonus', 'sales');
CREATE TYPE tarif_type_enum AS ENUM ('percent', 'fixed');
CREATE TYPE staff_type_enum AS ENUM ('shop_assistant', 'cashier');
create type repostitory_transaction_type_enum as enum ('minus', 'plus');

create table categories(
    id VARCHAR(40) primary key not null ,
    name VARCHAR(30),
    parent_id VARCHAR(40) references categories(id) default null,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at INTEGER DEFAULT 0
);

create table products(
    id uuid PRIMARY KEY NOT NULL ,
    name VARCHAR(30) UNIQUE,
    price INT ,
    barcode INT UNIQUE,
    category_id VARCHAR(40) REFERENCES categories(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at INTEGER DEFAULT 0
);

CREATE TABLE branches (
    id uuid PRIMARY KEY NOT NULL,
    name VARCHAR(30),
    address VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at INTEGER DEFAULT 0
);

create table repositories(
    id uuid primary key not null ,
    product_id uuid references products(id),
    branch_id uuid references branches(id),
    count int,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at INTEGER DEFAULT 0
);

CREATE TABLE sales (
    id UUID PRIMARY KEY NOT NULL,
    branch_id UUID REFERENCES branches (id),
    shop_assistant_id VARCHAR(80),
    cashier_id VARCHAR(80),
    payment_type payment_type_enum,
    price numeric DEFAULT 0,
    status status_enum DEFAULT 'in_process',
    client_name VARCHAR(30),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at INTEGER DEFAULT 0
);

create table baskets(
    id uuid primary key ,
    sale_id uuid references sales(id),
    product_id uuid references products(id),
    quantity int,
    price int,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at INTEGER DEFAULT 0
);

CREATE TABLE staff_tarifs (
    id UUID PRIMARY KEY,
    name VARCHAR(30) UNIQUE NOT NULL,
    tarif_type tarif_type_enum NOT NULL,
    amount_for_cash INT,
    amount_for_card INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at INTEGER DEFAULT 0
);

CREATE TABLE staffs (
    id UUID PRIMARY KEY,
    branch_id UUID REFERENCES branches(id),
    tariff_id UUID REFERENCES staff_tarifs(id),
    staff_type staff_type_enum NOT NULL,
    name VARCHAR(30),
    balance INT DEFAULT 0,
    age INT,
    birth_date DATE,
    login VARCHAR(15) UNIQUE,
    password VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at INTEGER DEFAULT 0
);

CREATE TABLE transactions (
    id uuid PRIMARY KEY NOT NULL,
    sale_id uuid REFERENCES sales (id),
    staff_id uuid REFERENCES staffs (id),
    transaction_type transaction_type_enum,
    source_type source_type_enum,
    amount numeric,
    description text,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at INTEGER DEFAULT 0
);

create table repository_transactions(
    id uuid primary key ,
    branch_id uuid references branches(id),
    staff_id uuid references staffs(id),
    product_id uuid references products(id),
    repository_transaction_type repostitory_transaction_type_enum,
    price int,
    quantity int,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at INTEGER DEFAULT 0
);