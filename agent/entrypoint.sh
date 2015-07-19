#!/bin/sh

printenv

if [ $# = 0 ]; then
  go run *.go
else
  sh -c $*
fi
