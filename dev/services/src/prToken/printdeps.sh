#!/bin/bash

export GOPATH=`cd ../../;pwd`

dep status | awk '{if(NR>1) print $1}' | tr "\n" " "
