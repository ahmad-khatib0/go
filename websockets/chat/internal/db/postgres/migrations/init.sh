#!/bin/sh

set -e

psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "postgres" <<-EOSQL
	  CREATE DATABASE chat ENCODING = utf8;
	    
	  CREATE USER chat WITH ENCRYPTED PASSWORD 'chat_db_password';
	    
	  ALTER DATABASE chat OWNER TO chat;
EOSQL
