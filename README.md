# e6e

Dependencies before getting started:

### GOPATH 
Golang is fairly opinionated on having a structured `$GOPATH`, and my code lives within my `$GOPATH/www-forecast`.

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
For a local e6e server, just running `revel run .` from the main `www-forecast` directory.

## Code Layout

The directory structure of a generated Revel application:

    conf/             Configuration directory
        app.conf      Main app configuration file
        routes        Routes definition file

    app/              App sources
        init.go       Interceptor registration
        controllers/  App controllers go here
        views/        Templates directory

    messages/         Message files

    public/           Public static assets
        css/          CSS files
        js/           Javascript files
        images/       Image files

    tests/            Test suites
