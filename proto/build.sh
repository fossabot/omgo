#!/usr/bin/env bash
protoc -I=. --go_out=./pb *.proto
