FROM golang:1.15.8

EXPOSE 9000

RUN go get -u github.com/revel/cmd/revel

WORKDIR $GOPATH/src/github.com/magoo/www-forecast

COPY . $GOPATH/src/github.com/magoo/www-forecast

RUN revel build github.com/magoo/www-forecast $GOPATH/bin/www-forecast prod

RUN chmod +x $GOPATH/bin/www-forecast

RUN rm -rf $GOPATH/src

CMD $GOPATH/bin/www-forecast/run.sh
