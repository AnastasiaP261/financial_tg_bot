CURDIR=$(shell pwd)
GOOSE_BIN=$(shell go env GOPATH)/bin/goose
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
MOCKGEN=${BINDIR}/mockgen_${GOVER}
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.49.0
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=gitlab.ozon.dev/apetrichuk/financial-tg-bot/cmd/bot
DEV_CREDS := postgresql://postgres:postgres@127.0.0.1:5432
DEV_DBNAME := finance
DEV_DBCONTAINER := postgres

all: format generate build test lint

build: bindir
	go build -o ${BINDIR}/bot ${PACKAGE}

test:
	go test ./...

run:
	go run ${PACKAGE}

generate: install-mockgen
	${MOCKGEN} -source=internal/model/messages/incoming_msg.go -destination=internal/model/messages/_mocks/mocks.go
	${MOCKGEN} -source=internal/model/purchases/model.go -destination=internal/model/purchases/_mocks/mocks.go
	${MOCKGEN} -source=internal/model/exchange_rates/model.go -destination=internal/model/exchange_rates/_mocks/mocks.go

lint: install-lint
	${LINTBIN} run

precommit: format build test lint
	echo "OK"

bindir:
	mkdir -p ${BINDIR}

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

install-mockgen: bindir
	test -f ${MOCKGEN} || \
		(GOBIN=${BINDIR} go install github.com/golang/mock/mockgen@v1.6.0 && \
		mv ${BINDIR}/mockgen ${MOCKGEN})

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

.PHONY: dev-db-data
__dev-db-data:
	${GOOSE_BIN} -dir ${CURDIR}/migrations postgres "${DEV_CREDS}/${DEV_DBNAME}?sslmode=disable" up
	${GOOSE_BIN} -dir ${CURDIR}/migrations postgres "${DEV_CREDS}/${DEV_DBNAME}?sslmode=disable" reset
	${GOOSE_BIN} -dir ${CURDIR}/migrations postgres "${DEV_CREDS}/${DEV_DBNAME}?sslmode=disable" up

docker-run:
	docker-compose up
