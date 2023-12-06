insert into url_shortener.url (id, original_url, short_url, user_id, deleted_flag) 
values ($1, $2, $3, $4, $5) 
returning id, original_url, short_url, user_id, deleted_flag;