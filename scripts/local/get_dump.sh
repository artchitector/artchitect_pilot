mkdir -p ~/dumps/
dt=$(date "+%Y%m%d")
path=~/dumps/dump$dt.sql.gz
scp memory:/root/dump.sql.gz ~/dumps/dump-$dt.sql.gz
