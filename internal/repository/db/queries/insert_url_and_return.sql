insert into url_shortener.url (id, original_url, short_url) 
values ($1, $2, $3) 
returning id, original_url, short_url