-- init.sql
DROP DATABASE IF EXISTS gcs_db WITH (FORCE);
CREATE DATABASE gcs_db
    WITH OWNER = admin
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TEMPLATE = template0;
