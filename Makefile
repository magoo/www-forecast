all:
			revel build github.com/magoo/www-forecast bin/www-forecast prod

debug:
			revel version
			ls ${GOPATH}/src
			ls ${GOPATH}/src/github.com/magoo/www-forecast
			ls ${GOPATH}/bin

test: revel test github.com/magoo/www-forecast
