dt=$(date "+%Y%m%d-%H%m")
echo "Starting get dump into file with timestamp - $dt"
path=/root/dumps/db/dump-$dt.sql.gz
pg_dump -d artchitect -U artchitector -h localhost -p 21431 | gzip > $path
echo "success. created file is $path. finish."
