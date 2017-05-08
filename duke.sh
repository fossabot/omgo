#!/bin/bash

# nuke unused docker images

read -p "Are you sure to remove all unused docker images?(yes/no)" -n 1 -r

echo

if [[ $REPLY =~ ^[Yy]$ ]]
then
  docker rm $(docker ps -a -q)
  docker rmi $(docker images -q)
fi
