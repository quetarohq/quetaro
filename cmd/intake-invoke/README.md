# quetaro-intake-invoke

quetaro-intake-invoke is a daemon when it dequeues messages from DB and invokes Lambda.

## Usage

```
Usage: quetaro-intake-invoke [OPTION]
  -aws-endpoint-url string
    	AWS endpoint URL. use $AWS_ENDPOINT_URL env
  -aws-region string
    	AWS region. use $AWS_REGION env (default "ap-northeast-1")
  -dsn string
    	database DSN. use $QTR_DATABASE_DSN env (e.g. 'postgres://username:password@localhost:5432')
  -err-interval duration
    	error wait interval. use $QTR_INTAKE_INVOKE_ERR_INTERVAL env (default 1m0s)
  -interval duration
    	poll interval. use $QTR_INTAKE_INVOKE_INTERVAL env (default 100ms)
  -nagents int
    	number of agents. use $QTR_INTAKE_INVOKE_NAGENTS env (default 1)
  -queue string
    	intake queue name. use $QTR_INTAKE_QUEUE env
  -version
    	print version
```
