create database request_proxy
	with owner docker
	encoding 'utf8'
    TABLESPACE = pg_default
	;
GRANT ALL PRIVILEGES ON database request_proxy TO docker;