CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.52.0
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
GOOSEBIN := ${BINDIR}/goose
DSN := "host=localhost port=5432 user=postgres password=postgres dbname=url sslmode=disable"
PACKAGE=github.com/Svoevolin/url-shortener/cmd/url-shortener

all: format build test lint

build: bindir
	go build -o ${BINDIR}/bot ${PACKAGE}

test:
	go test ./... -tags '!functional'

functional-test:
	go test ./... -tags 'functional'

run:
	CONFIG_PATH=./config/local.yaml go run ${PACKAGE}

lint: install-lint
	${LINTBIN} run

precommit: format build test lint
	echo "OK"

bindir:
	mkdir -p ${BINDIR}

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

goose-up: install-goose
	${GOOSEBIN} -dir ${CURDIR}/migrations postgres ${DSN} up

goose-status: install-goose
	${GOOSEBIN} -dir ${CURDIR}/migrations postgres ${DSN} status

goose-down: install-goose
	${GOOSEBIN} -dir ${CURDIR}/migrations postgres ${DSN} down

install-goose: bindir
	test -f ${GOOSEBIN} || GOBIN=${BINDIR} go install github.com/pressly/goose/cmd/goose@latest
	sudo chmod +x ${GOOSEBIN}