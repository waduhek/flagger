#!/usr/bin/env bash

# Name of the MongoDB container.
mongo_container_name=mongo
redis_container_name=redis

# Runs a docker container in detached mode.
#
# $1 Name to give to the container.
# $2 Port to publish.
# $3 Name of the image to run.
run_docker_container () {
    docker run --quiet --detach --name "$1" --publish "$2" "$3" >/dev/null 2>&1
}

# Waits for a container to change it's state to running.
#
# $1 Name given to the container.
wait_for_container () {
    until [ "$(docker inspect -f \{\{.State.Running\}\} "$1")" == "true" ] ; do
        sleep 0.1
    done
}

# Forcefully stops a docker container.
#
# $1 Name given to the container.
force_stop_container () {
    docker rm -f "$1" >/dev/null 2>&1
}

# Run MongoDB and Redis containers in detached mode and redirect all output to null.
run_docker_container $mongo_container_name "27017:27017" mongo
run_docker_container $redis_container_name "6379:6379" redis

# Wait until the container is ready.
wait_for_container $mongo_container_name
wait_for_container $redis_container_name

# Run tests.
go list ./... | grep -Ev "cmd|proto" | tr "\n" " " | xargs go test -coverprofile cover.out

# Stop and remove the MongoDB container.
force_stop_container $mongo_container_name
force_stop_container $redis_container_name
