# e6e

## Development workflow
Dependencies before getting started:

1) Clone the github repo
2) Setup Go & Dep
3) Setup Revel
4) Verify the install works
5) Setup Google Identity
6) Setup AWS credentials
7) Setup DynamoDB

### Go & GOPATH
Standard install of go works. Golang is fairly opinionated on having a structured `$GOPATH`, and my code lives within my `$GOPATH/src/www-forecast`.

### Dep
Install the [dep package manager for go](https://github.com/golang/dep), brew install dep and ensure dep.

### AWS IAM Account
Create a DynamoDB table and an IAM programmatic user account that can read/write to it.

Once you have these, make sure they are loaded in your path:
```
export AWS_SECRET_KEY= (AWS secret key)
export AWS_ACCESS_KEY= (AWS access Key)
```

### DynamoDB
Install terraform. The `tf` folder contains a terraform configuration to create the DynamoDB tables needed to operate.

`terraform apply` within the `tf directory` to set up.

### Google Identity
Need a set of Google API credentials to work with `http://localhost:4000` or whatever domain you're using.

- https://console.developers.google.com/apis/credentials

Once you have these, make sure they are loaded in your path:
```
export E6E_GOOGLE_SECRET= (Google Secret)
export E6E_GOOGLE_CLIENT= (Google Client)
```

### Installing and starting Revel
Install [revel command line tool](https://revel.github.io/tutorial/gettingstarted.html).

For a local e6e server, just running `revel run` from the main `www-forecast` directory.

## Production
This is currently a docker container (`Dockerfile` included) that is pushed to Fargate (An AWS service). Roles and environment are configured in production.
