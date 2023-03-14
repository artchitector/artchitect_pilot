rm ./dump.sql.gz
pg_dump -d artchitect -U artchitector -h *** -p *** | \
gzip > ./dump.sql.gz