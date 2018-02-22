FROM golang:1.8

RUN go get github.com/revel/cmd/revel

WORKDIR $GOPATH

ADD . $GOPATH/src/www-forecast

RUN revel build www-forecast $GOPATH/bin/www-forecast prod

RUN chmod +x $GOPATH/bin/www-forecast

RUN rm -rf $GOPATH/src

CMD $GOPATH/bin/www-forecast/run.sh

EXPOSE 9000

### docker run -i -t -p 8080:8080
