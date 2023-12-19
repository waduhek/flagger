# Flagger

Flagger aims to be a simple feature flag management system.

## Goals

* To allow users to sign-up.
* To allow users to authenticate in to the system.
* To allow authenticated users to create projects, environments and flags.
* To allow authenticated users to toggle the status of the flag i.e. on or off.
* To allow applications to fetch the current status of the provided flag.

## Tech Stack

Flagger is built using the Go programming language, MongoDB for database, Redis
as the cache, and GRPC for communications.

## Development

Flagger requires Docker and Docker Compose for running.

* To build the development image, run:

    ```bash
    make build-dev
    ```

* [Optional] To build the debug image, run:

    ```bash
    make build-debug
    ```

* After the image has been built, you'll need to start MongoDB and Redis. It is
  recommended to start these containers in a separate window using the following
  command:

    ```bash
    make run-others
    ```

* To run the development image, run:

    ```bash
    make run-dev
    ```

* [Optional] To run the debug image, run:

    ```bash
    make run-debug
    ```

  The debug container will open a Delve debug port on port 4040. Connect to it
  using:

    ```bash
    dlv connect localhost:4040
    ```

### Tearing down

* To stop all running containers, run:

    ```bash
    make teardown-all
    ```

* To stop only MongoDB and Redis containers, run:

    ```bash
    make teardown-others
    ```

* To stop only the debug container, run:

    ```bash
    make teardown-debug
    ```

* To stop only the development container, run:

    ```bash
    make teardown-dev
    ```
