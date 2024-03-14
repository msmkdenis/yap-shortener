create schema if not exists url_shortener;

create table if not exists url_shortener.url
(
    id                      text,
    original_url            text,
    short_url               text,
    correlation_id          text,
    user_id                 text,
    deleted_flag            bool default false,
    constraint pk_url primary key (id)
);