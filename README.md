![log](https://user-images.githubusercontent.com/117768/205684311-87faae92-5f36-4545-a504-01f3a742477d.png)

Quetaro is a job queue system using SQS, Lambda and PostgreSQL.

[![test](https://github.com/quetarohq/quetaro/actions/workflows/test.yml/badge.svg)](https://github.com/quetarohq/quetaro/actions/workflows/test.yml)

## System Architecture

cf. https://github.com/quetarohq/quetaro/tree/main/_etc/terraform

![fig](https://user-images.githubusercontent.com/117768/206375134-d25af7f9-2d6a-494a-80b7-5ab7350e9c31.png)

## Getting Started for local

### Start docker-compose

```sh
$ docker compose up
```

### Create AWS resources to LocalStack

```sh
$ make tf-init
$ make tf-apply
```

NOTE: When docker-compose is stopped, LocalStack resources will disappear. After restarting docker-compose, run `make tf-apply restart` again.

### Setup DB

```sh
$ make db
```

### Restart daemons

```sh
$ make restart
```

### Send a message

```sh
$ make message
# or `make failure`
```

### Browse DB

Open http://localhost:8081/

## Getting Started for AWS

### Create AWS resources

```sh
$ cd _etc/terraform/
$ terraform workspace new aws
$ terraform init
$ terraform apply
```

### Start docker-compose for AWS

```sh
$ docker compose -f compose.aws.yml up
```

### Setup DB

```sh
$ make db
```

### Send a message

```sh
$ make message-for-aws
# or `make failure-for-aws`
```

## Run tests

```sh
$ docker compose up localstack db
```

```sh
$ make tf-init
$ make tf-apply
$ make db
$ make test
```
