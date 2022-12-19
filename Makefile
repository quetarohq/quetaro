.PHONY: all
all: vet build

.PHONY: build
build:
	go build ./cmd/intake-invoke
	go build ./cmd/intake-pull
	go build ./cmd/outlet-failure
	go build ./cmd/outlet-success

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test -v -count=1 ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	rm -f intake-invoke intake-pull outlet-failure outlet-success

.PHONY: tf-init
tf-init:
	docker-compose run terraform init -upgrade

.PHONY: tf-apply
tf-apply:
	docker-compose run terraform apply -auto-approve

.PHONY: db
db:
	docker-compose run psql ./setup.sh
	docker-compose run -e PGDATABASE=qtr_test psql ./setup.sh

.PHONY: psql
psql:
	docker-compose run psql 'psql -h db -U qtr'

.PHONY: restart
restart:
	docker-compose restart intake-invoke intake-pull outlet-failure outlet-success lambda-tailf pgweb

.PHONY: message
message:
	docker-compose run \
	awscli sqs send-message --region us-east-1 --endpoint-url http://localstack:4566 \
		--queue-url qtr-intake \
		--message-attributes "FunctionName={StringValue=qtr-job,DataType=String}" \
		--message-body '{"date":"$(shell date)","_fail":"$(JOB_FAIL)"}'

.PHONY: failure
failure: JOB_FAIL:=true
failure: message

.PHONY: message-for-aws
message-for-aws:
	aws sqs send-message \
		--queue-url $(shell aws sqs get-queue-url --queue-name qtr-intake --output text) \
		--message-attributes "FunctionName={StringValue=qtr-job,DataType=String}" \
		--message-body '{"date":"$(shell date)","_fail":"$(JOB_FAIL)"}'

.PHONY: failure-for-aws
failure-for-aws: JOB_FAIL:=true
failure-for-aws:
