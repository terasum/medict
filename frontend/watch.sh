#!/bin/bash
WPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

if [ -f $WPATH/dist/watch.pid ];then
    WPID=$(cat dist/watch.pid)
    if ps -p $WPID > /dev/null; then
        echo "$WPID is running, killing"
        kill -9 $WPID
    fi
fi

nohup yarn run watch > $WPATH/dist/watch.out 2>&1 &
WPID=$!
echo $WPID > $WPATH/dist/watch.pid