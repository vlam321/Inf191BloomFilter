#!/bin/bash

if [ "$1" == "bloomRouter" ]; then
	go run "cmd/$1/$1.go";
fi

if [ "$1" == "bloomFilterServer" ]; then
	go run "cmd/$1/$1.go";
fi

