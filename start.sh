#!/bin/sh

nohup ./statussvr > /dev/null 2>&1 &
echo $! > pid
