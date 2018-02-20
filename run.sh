#!/bin/bash

if [ "$1" == "bloomRouter" ]; then
	echo "$1";
	go run "cmd/$1/$1.go";
fi

if [ "$1" == "bloomFilterServer" ]; then
	echo "$2";
	go run "cmd/$1/$1.go" 1000000 10;
fi

