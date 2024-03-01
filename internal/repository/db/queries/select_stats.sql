select
    count (distinct short_url) as urls,
    count (distinct user_id) as users
from url_shortener.url