#!/bin/sh

SCRIPTPATH=$(cd "$(dirname "$0")"; pwd)
"$SCRIPTPATH/www-forecast" -importPath  -srcPath "$SCRIPTPATH/src" -runMode dev
