select id, user_id, deleted_flag
from url_shortener.url
where user_id = $1 and id = $2 for update