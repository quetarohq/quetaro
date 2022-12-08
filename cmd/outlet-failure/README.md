# quetaro-outlet-failure

quetaro-outlet-failure is a daemon that receives the results of failed Lambda executions from SQS($QTR_OUTLET_FAILURE_QUEUE) and updates the DB.

![image](https://user-images.githubusercontent.com/117768/206354240-0f0ffdeb-b743-45ac-bbaa-a8e24ecea110.png)

## Usage

```
Usage: quetaro-outlet-failure [OPTION]
  -aws-endpoint-url string
    	AWS endpoint URL. use $AWS_ENDPOINT_URL env
  -aws-region string
    	AWS region. use $AWS_REGION env (default "ap-northeast-1")
  -dsn string
    	database DSN. use $QTR_DATABASE_DSN env (e.g. 'postgres://username:password@localhost:5432')
  -err-interval duration
    	error wait interval. use $QTR_OUTLET_FAILURE_ERR_INTERVAL env (default 1m0s)
  -interval duration
    	poll interval. use $QTR_OUTLET_FAILURE_INTERVAL env (default 1s)
  -max-recv int
    	maximum number of received messages. use $QTR_OUTLET_FAILURE_MAX_RECV env (default 1)
  -nagents int
    	number of agents. use $QTR_OUTLET_FAILURE_NAGENTS env (default 1)
  -queue string
    	outlet-failure queue name. use $QTR_OUTLET_FAILURE_QUEUE env
  -version
    	print version
```
