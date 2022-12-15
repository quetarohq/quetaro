#!/bin/bash
SCRIPT_DIR=$(readlink -f $(dirname $0))

: ${PGHOST:=localhost}
: ${PGPORT:=15432}
: ${PGDATABASE:=qtr}
export PGHOST PGPORT PGDATABASE

export PGUSER=postgres
createdb $PGDATABASE
createuser -s $PGDATABASE

export PGUSER=$PGDATABASE
psql -c "create schema qtr authorization $PGUSER"
psql -c 'alter role qtr set search_path = "$user"'
psql -c 'create extension pgcrypto'
psql -f $SCRIPT_DIR/schema.sql
