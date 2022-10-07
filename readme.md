

http://localhost:8085/view/test


docker build -t app-prod . --target production

docker run -p 8085:8085 --name app-prod app-prod

docker build -t maxistar/app-mongowiki . --target production
docker push  maxistar/app-mongowiki


Links:
[Dockerise go application](https://dev.to/karanpratapsingh/dockerize-your-go-app-46pp)