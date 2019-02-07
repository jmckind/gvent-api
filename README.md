# gvent-api

Repo for the gvent-api component.

## Overview

gvent is a "simple" event management system written in Go. It is intended to
dempnstrate best practices for deployment of a cloud native application
based on a microservice architecture.

The gvent-api component is responsible for providing a REST API to manage the
system data model.

## Deployment

These instructions assume a working Kubernetes cluster with admin rights.

Optionally, create a namespace for the application.

```bash
kubectl create ns gvent
```

The first step is to ensure there is a rethinkdb database available. Use the included manifest to deploy a sample database. A deployment with 1 replica, along with a service should be created.

```bash
kubectl create -n gvent -f deploy/gvent-api.yaml
```

Once the database service is available, use the included manifest to deploy the application. A deployment with 1 replica, along with a service should be created.

```bash
kubectl create -n gvent -f deploy/gvent-api.yaml
```

View the list of pods in the namespace to verify everything is running.

```bash
kubectl get pods -n gvent
```

The api and database pods should be in a `Running` state.

```bash
NAME                         READY     STATUS      RESTARTS   AGE
gvent-api-59b6f55fd8-5zthf   1/1       Running     0          15m
gvent-db-fc4c4bb6d-9mx85     1/1       Running     0          16m
```

## License

gvent-api is released under the Apache 2.0 license. See the
[LICENSE][license_file] file for details.

[license_file]:./LICENSE
