#!/bin/sh

GOOS=`uname | tr '[:upper:]' '[:lower:]'`
GOARCH=386
if [ "x86_64" = `uname -m` ]; then
  GOARCH=amd64
fi
cmd="$0.$GOOS-$GOARCH"

exec "$cmd" "$@"
