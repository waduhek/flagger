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

Flagger requires Docker and Kubernetes for running.

0. Copy the sample env files `sample.env` and `secret_sample.env` files to
   `k8s/base` using:

     ```bash
     cp sample.env k8s/base/.env
     cp secret_sample.env k8s/base/.env.secret
     ```

   Once copied enter the values for the variables in `k8s/base/.env.secret`.


1. Start a Kubernetes cluster on your machine. I've used
   [`minikube`](https://minikube.sigs.k8s.io/docs/) and so the instructions here
   are going to be with reference to `minikube`. Start a `minikube` cluster
   using:

     ```bash
     minikube start
     ```

2. Set up the Docker environment for `minikube` using:

     ```bash
     eval $(minikube docker-env)
     ```

   This will allow us to build the images directly into `minikube`.

3. To build the development image, run:

     ```bash
     make build-dev
     ```

4. [Optional] To build the debug image, run:

     ```bash
     make build-debug
     ```

5. To start the application in development mode, run:

     ```bash
     kubectl apply -k k8s/overlays/dev
     ```

6. [Optional] To start the application in debug mode, run:

     ```bash
     kubectl apply -k k8s/overlays/debug
     ```

7. When the application first starts the API server will not be able to run as
   MongoDB replicaset has not been setup up yet. To do so first, we to spin up a
   new MongoDB pod using:

     ```bash
     kubectl run mongo --image mongo --rm -it -- bash
     ```

8. After a bash session has started, connect to the primary MongoDB instance
   using:

     ```bash
     mongosh mongodb://mongo-0.mongo-hlsvc
     ```

9. Initiate the replicaset using the following `mongosh` command:

     ```js
     rs.initiate({
         _id: "rs0",
         members: [
             { _id: 0, host: "mongo-0.mongo-hlsvc" },
             { _id: 1, host: "mongo-1.mongo-hlsvc" },
             { _id: 2, host: "mongo-2.mongo-hlsvc" },
         ],
     });
     ```

10. The replicaset will take a little while to initiate. The status of the
    replicaset can be checked using the following `mongosh` command:

      ```js
      rs.status();
      ```

    Once you see one primary server and 2 secondary servers, you are ready to
    go.

11. Now, just wait for the API server to restart automatically and every thing
    should work.

12. To make requests to the API, you'll first need to get the IP of `minikube`
    using:

      ```bash
      minikube ip
      ```

13. Next you'll need the port on which the API server is running. To get that
    run:

      ```bash
      kubectl get service/api-server
      ```

    The port that you need will be mapped under the "PORT(S)" column.

14. Now you can make requests to the service on the IP `$(minikube ip):<PORT>`.

**Note**: Once you delete your local `minikube` cluster, you'll have to repeat
steps 7 through 11 to setup the replicaset.

### Tearing down

To stop the all the deployed units in the cluster, run:

  ```bash
  kubectl delete -k k8s/overlays/<overlay>
  ```

Here, replace `<overlay>` with the overlay used to deploy units.

To delete all resources started by `minikube`, run:

  ```bash
  minikube delete
  ```
