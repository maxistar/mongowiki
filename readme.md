# Mongo WIKI

GoLang + Minikube Playground

## Access Database Access

1. go to [Mongo Cloud](https://cloud.mongodb.com/) and create database
2. make sure it is accessible from your location, copy connection string

## Set Environment Variables

````shell
export MONGO_DB_NAME=myDB
export MONGO_COLLECTION_NAME=mongowiki
export MONGO_CONNECTION_STRING="mongodb+srv://username:password@cluster0.x5da5.mongodb.net/?retryWrites=true&w=majority"
````

In Windows

```shell
[Environment]::SetEnvironmentVariable('MONGO_DB_NAME', 'myFirstDatabase')
[Environment]::SetEnvironmentVariable('MONGO_COLLECTION_NAME', 'mongowiki')
[Environment]::SetEnvironmentVariable('MONGO_CONNECTION_STRING', 'mongodb+srv://username:password@cluster0.wpell50.mongodb.net/?retryWrites=true&w=majority')

```


## Run Locally

### Run

````shell
go run ./mongowiki.go
````

open browser on http://localhost:8085

### Compile and run

````shell
go build
./mongowiki
````

http://localhost:8085/view/test

## Docker

## Build Docker Image
 
````shell
docker build -t maxistar/app-mongowiki . --target production
````

## Run in docker locally

````shell
docker run -e MONGO_DB_NAME -e MONGO_COLLECTION_NAME -e MONGO_CONNECTION_STRING -p 8085:8085 maxistar/app-mongowiki
````

## Push Docker images

````shell
docker push maxistar/app-mongowiki
````

## Kubernetes

Run Minikube

````shell
alias k="minikube kubectl --"

````
### Apply Pon Configuration

````shell
k apply -f ./kubernetes/manifest.yaml
````

````shell
k expose deployment mongo-wiki-demo --port=8085 --target-port=8085 --type=LoadBalancer
````

### Links:

[Dockerise go application](https://dev.to/karanpratapsingh/dockerize-your-go-app-46pp)
