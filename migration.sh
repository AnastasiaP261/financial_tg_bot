#!/bin/bash

CONFIG_PATH=./build/config/config.yaml
MIGRATION_DIR=./build/database/migrations
GOOSE_BIN=go/bin/goose

arr=($(cat ${CONFIG_PATH} | grep "dbUri: "))
POSTGRESQL_DSN=${arr[1]}
#echo "postgresql dsn: ${POSTGRESQL_DSN}"

${GOOSE_BIN} -dir ${MIGRATION_DIR} postgres "${POSTGRESQL_DSN}" up
${GOOSE_BIN} -dir ${MIGRATION_DIR} postgres "${POSTGRESQL_DSN}" reset
${GOOSE_BIN} -dir ${MIGRATION_DIR} postgres "${POSTGRESQL_DSN}" up
