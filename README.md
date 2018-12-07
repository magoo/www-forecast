# e6e
This is a Revel web application with a DynamoDB backend, requiring Google developer credentials.

## Development workflow
Dependencies before getting started:

1) Setup Go & Dep
2) Clone the github repo
3) Setup Revel
4) Verify the install works
5) Setup AWS credentials
6) Setup DynamoDB
7) Setup Google Identity


### 1. Go & GOPATH
Currently developing with `go1.11.1`. Golang is picky on having a structured `$GOPATH`.

Clone this repo into `$GOPATH/src/www-forecast`.

### 2. Dep
Install the [dep package manager for go](https://github.com/golang/dep), brew install dep and ensure dep.


### 3. Set up Revel
You'll need the revel command line tool.

```bash
go get -u github.com/revel/cmd/revel
```

### 4. Verify the install works.
Enter the `$GOPATH/src/www-forecast` directory and `revel run`. The server should start locally. It won't work quite yet as we haven't setup AWS and Google credentials, but this is a good point to stop and troubleshoot any issues with Go or Revel.

### 5. AWS IAM Account
Create an IAM programmatic user account with permission to modify DynamoDB.

Once you have these, make sure they are loaded in your path:
```
export AWS_SECRET_KEY= (AWS secret key)
export AWS_ACCESS_KEY= (AWS access Key)
```

### 6. DynamoDB
Install terraform. The `tf` folder contains a terraform configuration to create the DynamoDB tables needed to operate.

`terraform apply` within the `tf/` directory to set up tables.

> You can use the `E6E_TABLE_PREFIX` environment variable to point the app at a specific set of tables, but you'll have to modify the terraform script to name these tables with your chosen prefix.

### 7. Google Identity
You'll need a set of Google API/OAuth credentials to work with `http://localhost:9000` or whatever domain you'll be using.

- https://console.developers.google.com/apis/credentials

Once you have these, make sure they are loaded in your path:
```
export E6E_GOOGLE_CLIENT= (Google Client)
export E6E_GOOGLE_SECRET= (Google Secret)
```

## e6e.io Production
This is currently a docker container (`Dockerfile` included) that is pushed to Fargate (An AWS service). The Fargate configuration is manual and not yet documented. Currently, roles and environment are configured in production.

1. `docker build -t scrty .`
2. `docker tag scrty:latest 832911230879.dkr.ecr.us-east-1.amazonaws.com/scrty:latest`
3. `aws ecr get-login --no-include-email --region us-east-1` (change profile if needed)
4. (copy code from #3)
5. `docker push 832911230879.dkr.ecr.us-east-1.amazonaws.com/scrty:latest`
6. `aws ecs update-service --region us-east-1 --force-new-deployment --service e6e-service-prod --cluster e6e-cluster-prod` (change profile if needed)

## Configuration
There are a number of environment variables that can change how the application functions.

### E6E_CONFIG_PATHS
To specify a custom revel config you can add a comma separated list of directories containing app.conf files that revel will check. For more info see https://revel.github.io/manual/appconf.html#external_app.conf

### E6E_QUESTIONS_TABLE_NAME 
To override the default dynamodb table name for use in custom deployments use this variable. The default is `questions-tf`

### E6E_ANSWERS_TABLE_NAME 
To override the default dynamodb table name for use in custom deployments use this variable. The default is `answers-tf`

### E6E_GOOGLE_CLIENT
This is the Oauth client id for logging into the application with Google.

### E6E_GOOGLE_SECRET
This is the Oauth client secret for logging into the application with Google.



