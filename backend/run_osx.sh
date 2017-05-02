#!/usr/bin/env bash

export EXTERNAL_IP=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
exec docker-compose $@
