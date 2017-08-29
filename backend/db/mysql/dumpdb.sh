#!/bin/bash

if [ "$1" == "" ] ; then
    echo "filename missing"
    echo "dumpdb.sh [filename]"
    exit
fi

# sh dumpdb.sh docker-entrypoint-initdb.d/backup.sql

mysqldump -h127.0.0.1 -P 3306 -uroot -psupersecret --routines --flush-privileges --databases mysql master > "$1"
