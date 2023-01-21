-- Run it from root
sudo -u postgres psql -- create DB
create database artchitect;
create user artchitector with encrypted password 'Zz123456';
ALTER DATABASE artchitect OWNER TO artchitector;
sudo -u postgres psql -d artchitect
grant all privileges on database artchitect to artchitector;

-- # dumps
-- backup
pg_dump -d artchitect -U artchitector -h localhost -p 21431 | gzip > dump20230120.sql.gz
-- restore
gunzip dump20230120.sql.gz
psql -U artchitector -d artchitect -h localhost -p 21431 -f dump20230120.sql -v ON_ERROR_STOP=1