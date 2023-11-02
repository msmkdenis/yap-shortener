create schema if not exists raw_data;

create table if not exists raw_data.sales
(
    id                      integer,
    auto                    text,
    gasoline_consumption    text,
    price                   numeric(9,2),
    date                    date,
    person_name             text,
    phone                   text,
    discount                float,
    brand_origin            text
);

create schema if not exists car_shop;

create table if not exists car_shop.country
(
    country_id              serial,
    name                    varchar(50) not null,
    constraint pk_country primary key (country_id),
    constraint country_unique_name unique (name)
);

create table if not exists car_shop.brand
(
    brand_id                serial,
    name                    varchar(50),
    country_id              integer,
    constraint pk_brand primary key (brand_id),
    constraint fk_country foreign key (country_id) references car_shop.country (country_id),
    constraint brand_unique_name unique (name)
);

create table if not exists car_shop.car
(
    car_id                   serial,
    name                     varchar(50) not null,
    brand_id                 integer not null,
    gasoline_consumption     float,
    constraint pk_car primary key (car_id),
    constraint fk_brand foreign key (brand_id) references car_shop.brand (brand_id),
    constraint car_unique_name unique (name)
);

create table if not exists car_shop.colour
(
    colour_id                serial,
    name                     varchar(50) not null,
    constraint pk_colour primary key (colour_id),
    constraint colour_unique_name unique (name)
);

create table if not exists car_shop.client
(
    client_id                serial,
    name                     varchar(50) not null,
    phone                    varchar(25) not null,
    constraint pk_client primary key (client_id),
    constraint client_unique_phone unique (phone)
);

create table if not exists car_shop.invoice
(
    invoice_id               serial,
    date                     date not null,
    discount                 float not null,
    price                    numeric(9, 2) not null,
    car_id                   integer not null,
    client_id                integer not null,
    colour_id                integer not null,
    constraint pk_invoice primary key (invoice_id),
    constraint fk_car foreign key (car_id) references car_shop.car (car_id),
    constraint fk_client foreign key (client_id) references car_shop.client (client_id),
    constraint fk_colour foreign key (colour_id) references car_shop.colour (colour_id),
    constraint positive_discount check (discount >= 0), -- скидка не может быть отрицательной
    constraint positive_price check (price > 0) -- цена не может быть отрицательной
);