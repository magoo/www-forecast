all:
			revel build github.com/magoo/www-forecast bin/www-forecast prod

debug:
			revel version
			ls $GOPATH/src
			ls $GOPATH/github.com/magoo/www-forecast
			ls $GOPATH/bin
