#!/bin/bash

CONTAINER_IDS=$( docker ps -a | awk '($2 ~ /dev-peer.*/) {print $1}')
if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" == " " ]; then
    echo "---- No containers available for deletion ----"
else
    docker rm -f $CONTAINER_IDS
fi