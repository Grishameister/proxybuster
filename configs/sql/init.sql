create database request_proxy
	with owner postgres
	encoding 'utf8'
    TABLESPACE = pg_default
	;
GRANT ALL PRIVILEGES ON database request_proxy TO postgres;