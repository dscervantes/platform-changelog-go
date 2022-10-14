#!/bin/bash

# For setting it up with in a clowder env
#./app-common-bash > /tmp/vars.sh
# source /tmp/vars.sh

# Variables for running the job locally
CLOWDER_DATABASE_HOSTNAME=localhost
CLOWDER_DATABASE_PORT=5432
CLOWDER_DATABASE_NAME=gumbaroo
CLOWDER_DATABASE_USERNAME=crc
CLOWDER_DATABASE_PASSWORD=crc
RETENTION_DAYS=90

RETENTION_DAYS=${RETENTION_DAYS:-7}

PGPASSWORD=$CLOWDER_DATABASE_PASSWORD psql -h $CLOWDER_DATABASE_HOSTNAME --user $CLOWDER_DATABASE_USERNAME --db $CLOWDER_DATABASE_NAME -c "DELETE from timelines WHERE timestamp < (NOW() - interval '$RETENTION_DAYS days');"
PGPASSWORD=$CLOWDER_DATABASE_PASSWORD psql -h $CLOWDER_DATABASE_HOSTNAME --user $CLOWDER_DATABASE_USERNAME --db $CLOWDER_DATABASE_NAME -c "VACUUM ANALYZE timelines;"