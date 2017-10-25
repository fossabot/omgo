#!/usr/bin/env bash

mongodump --db somedb --collection somecollection --out - | gzip > collectiondump.gz

mongodump --db somedb --collection somecollection --out - | gzip > dump_`date "+%Y-%m-%d"`.gz

gunzip

mongorestore
