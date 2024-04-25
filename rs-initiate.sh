#!/bin/bash

initialise() {
    while true
    do
        if mongosh --host mongodb --eval "quit()" &> /dev/null; then
            break
        fi
    done
    mongosh --host mongodb --eval "rs.initiate()"
}

initialise &

mongod --replSet dbrs --bind_ip_all