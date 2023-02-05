select t.day, concat(sum(t.size) / 1048576.0, ' MB')
from (select i.card_id,
             length(data)                    as size,
             date_trunc('day', c.created_at) as day
      from images i
          join cards c on i.card_id = c.id
      where c.created_at between '2023-01-01' and '2023-02-03'
     ) t group by t.day;