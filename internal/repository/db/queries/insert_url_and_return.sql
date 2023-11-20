insert into url_shortener.url (id, original_url, short_url, user_id) 
values ($1, $2, $3, $4) 
returning id, original_url, short_url, user_id;