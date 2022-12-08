# quetaro-outlet-success

quetaro-outlet-success is a daemon that receives the results of a successful Lambda execution from SQS($QTR_OUTLET_SUCCESS_QUEUE) and updates the DB.

## Usage

```
Usage: quetaro-outlet-success [OPTION]
  -aws-endpoint-url string
    	AWS endpoint URL. use $AWS_ENDPOINT_URL env
  -aws-region string
    	AWS region. use $AWS_REGION env (default "ap-northeast-1")
  -dsn string
    	database DSN. use $QTR_DATABASE_DSN env (e.g. 'postgres://username:password@localhost:5432')
  -err-interval duration
    	error wait interval. use $QTR_OUTLET_SUCCESS_ERR_INTERVAL env (default 1m0s)
  -interval duration
    	poll interval. use $QTR_OUTLET_SUCCESS_INTERVAL env (default 1s)
  -max-recv int
    	maximum number of received messages. use $QTR_OUTLET_SUCCESS_MAX_RECV env (default 1)
  -nagents int
    	number of agents. use $QTR_OUTLET_SUCCESS_NAGENTS env (default 1)
  -queue string
    	outlet-success queue name. use $QTR_OUTLET_SUCCESS_QUEUE env
  -version
    	print version
```
