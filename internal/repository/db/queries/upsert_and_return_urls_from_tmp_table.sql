insert into url_shortener.url (id, original_url, short_url, correlation_id, user_id, deleted_flag) 
select id, original_url, short_url, correlation_id, user_id, deleted_flag from pg_temp.%s 
on conflict (id) do update set original_url = excluded.original_url, short_url = excluded.short_url, correlation_id = excluded.correlation_id, user_id = excluded.user_id, deleted_flag = excluded.deleted_flag
returning id, original_url, short_url, correlation_id, user_id, deleted_flag 