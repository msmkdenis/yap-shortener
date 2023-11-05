create schema if not exists url_shortener;

create table if not exists url_shortener.url
(
    id                      varchar(16),
    iriginal_url            text,
    short_url               text,
    constraint pk_url primary key (id)
);