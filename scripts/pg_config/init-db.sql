-- Run it from root (postgres) user.
create database artchitect;
create user artchitector with encrypted password 'Zz123456';
ALTER DATABASE artchitect OWNER TO artchitector;
grant all privileges on database artchitect to artchitector;
