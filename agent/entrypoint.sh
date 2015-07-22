#!/bin/bash

if [ $# = 0 ]; then
  go run $(ls -1 *.go | grep -v _test.go)
else
  sh -c $*
fi
