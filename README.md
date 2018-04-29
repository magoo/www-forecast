# e6e

## Development workflow
Dependencies before getting started:

### GOPATH
Golang is fairly opinionated on having a structured `$GOPATH`, and my code lives within my `$GOPATH/src/www-forecast`.

### Dep
Install the [dep package manager for go](https://github.com/golang/dep), `brew install dep` and `ensure dep`.

### DynamoDB
The `tf` folder contains a terraform configuration to create the DynamoDB tables needed to operate.

`terraform apply` within the `tf directory` to set up.

### Google Identity

Need a set of Google API credentials to work with `http://localhost:4000` or whatever domain you're using.

- https://console.developers.google.com/apis/credentials

### Environment Variables

In a local environment, AWS credentials are pulled from the environment and only need DynamoDB access. Currently using Fargate and IAM roles for production (e6e.io).

Currently sourcing a .profile script for local development, like so:

```
export AWS_SECRET_KEY= (AWS secret key)
export AWS_ACCESS_KEY= (AWS access Key)
export E6E_GOOGLE_SECRET= (Google Secret)
export E6E_GOOGLE_CLIENT= (Google Client)
```

### Starting Revel
Install [revel command line tool](https://revel.github.io/tutorial/gettingstarted.html).

For a local e6e server, just running `revel run` from the main `www-forecast` directory.

## Production
This is currently a docker container (`Dockerfile` included) that is pushed to Fargate (An AWS service). Roles and environment are configured in production. 
