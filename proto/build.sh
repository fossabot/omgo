#!/bin/bash
protoc -I=. --go_out=./pb *.proto
