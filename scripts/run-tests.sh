#!/usr/bin/env bash

# Name of the MongoDB container.
mongo_container_name=mongo
mongo_container_port=27017
redis_container_name=redis
redis_container_port=6379

# Runs a docker container in detached mode.
#
# $1 Name to give to the container.
# $2 Port to publish.
# $3 Name of the image to run.
run_docker_container () {
    if ! docker run --quiet --detach --name "$1" --publish "$2" "$3" >/dev/null 2>&1
    then
        echo "Error while starting container with name $1"
        exit 1
    fi
}

# Gets the IP address of a container
#
# $1 Name given to the container.
get_container_ip () {
    docker inspect -f \{\{.NetworkSettings.IPAddress\}\} "$1"
}

# Waits for a container to change it's state to running.
#
# $1 Name given to the container.
# $2 Port the container is running on.
wait_for_container () {
    if ! command -v nc >/dev/null
    then
        echo "nc does not exist"
        exit 1
    fi

    container_ip=$(get_container_ip "$1")

    until nc -z "$container_ip" "$2" ; do
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
run_docker_container $mongo_container_name "$mongo_container_port:$mongo_container_port" mongo
run_docker_container $redis_container_name "$redis_container_port:$redis_container_port" redis

# Wait until the container is ready.
wait_for_container $mongo_container_name $mongo_container_port
wait_for_container $redis_container_name $redis_container_port

# Run tests.
go list ./... | grep -Ev "cmd|proto" | tr "\n" " " | xargs go test -coverprofile cover.out

# Stop and remove the MongoDB container.
force_stop_container $mongo_container_name
force_stop_container $redis_container_name
