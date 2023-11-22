select id, original_url, short_url, coalesce(correlation_id, '')
from url_shortener.url
where id = $1