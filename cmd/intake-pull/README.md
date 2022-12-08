# quetaro-intake-pull

quetaro-intake-pull is a daemon that receives messages from SQS($QTR_INTAKE_QUEUE) and enqueues them to DB.

![image](https://user-images.githubusercontent.com/117768/206353932-09fc27fb-2f08-4f29-af53-17aeb3747f4a.png)

## Usage

```
Usage: quetaro-intake-pull [OPTION]
  -aws-endpoint-url string
    	AWS endpoint URL. use $AWS_ENDPOINT_URL env
  -aws-region string
    	AWS region. use $AWS_REGION env (default "ap-northeast-1")
  -dsn string
    	database DSN. use $QTR_DATABASE_DSN env (e.g. 'postgres://username:password@localhost:5432')
  -err-interval duration
    	error wait interval. use $QTR_INTAKE_PULL_ERR_INTERVAL env (default 1m0s)
  -interval duration
    	poll interval. use $QTR_INTAKE_PULL_INTERVAL env (default 1s)
  -max-recv int
    	maximum number of received messages. use $QTR_INTAKE_PULL_MAX_RECV env (default 1)
  -nagents int
    	number of agents. use $QTR_INTAKE_PULL_NAGENTS env (default 1)
  -queue string
    	intake queue name. use $QTR_INTAKE_QUEUE env
  -version
    	print version
```
