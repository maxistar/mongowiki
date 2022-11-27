# Mongo WIKI

## Access Database Access
https://cloud.mongodb.com/

## Compile
go build
http://localhost:8085/view/test

## Build Dockek Image

docker build -t maxistar/app-mongowiki . --target production
docker run -p 8085:8085 maxistar/app-mongowiki
docker push  maxistar/app-mongowiki


Links:
[Dockerise go application](https://dev.to/karanpratapsingh/dockerize-your-go-app-46pp)