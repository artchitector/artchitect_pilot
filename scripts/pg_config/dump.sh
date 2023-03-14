dt=$(date "+%Y%m%d-%H%m")
pg_dump -d artchitect -U artchitector -h localhost -p 21431 | gzip > /root/dumps/db/dump-$dt.sql.gz
