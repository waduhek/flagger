#!/usr/bin/env bash

# Name of the MongoDB container.
mongo_container_name=mongo

# Run MongoDB container in detached mode and redirect all output to null.
docker run --quiet --detach --name mongo --publish 27017:27017 $mongo_container_name >/dev/null 2>&1

# Wait until the container is ready.
until [ "$(docker inspect -f \{\{.State.Running\}\} $mongo_container_name)" == "true" ] ; do
    sleep 0.1
done

# Run tests.
go test -tags=integrationtest ./...

# Stop and remove the MongoDB container.
docker rm -f $mongo_container_name >/dev/null 2>&1
