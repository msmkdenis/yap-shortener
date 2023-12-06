select id, original_url, short_url, coalesce(correlation_id, ''), user_id, deleted_flag
from url_shortener.url
where user_id = $1;