#!/usr/bin/env bash

scp -P5555 ./libs/log/*.go work@120.132.50.66:/home/work/go/src/GoSpider/libs/log
scp -P5555 ./spider/*.go work@120.132.50.66:/home/work/go/src/GoSpider/spider