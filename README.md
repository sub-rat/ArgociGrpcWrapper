# Argo Ci wrapper
REST api on top of Argo Ci workflow using GRPC connection

## Steps to run project
## setup argo ci
```
kubectl create namespace argo
kubectl apply -n argo -f https://raw.githubusercontent.com/argoproj/argo/stable/manifests/install.yaml
```

#### Grant default service account
```
kubectl create rolebinding default-admin --clusterrole=admin --serviceaccount=default:default
```

For your own service account and detail configuration of Argo ci visit to their official documentation
[https://argoproj.github.io/](https://argoproj.github.io/)

For creating Workflow visit to example of argo workflow [https://argoproj.github.io/docs/argo/examples/readme.html](https://argoproj.github.io/docs/argo/examples/readme.html)

## Install Postgres 
Postgres is used to maintain user database for authentication

follow official setup guide [https://www.postgresql.org/download/](https://www.postgresql.org/download/)

## Running Api 
first clone this repository and run following command
```
go mod tidy
```

Change the environment variables as per your environment

to directly run use following
```
go run main.go
```

to build binary and run go with this method
```
go build .
./ArogciGrpcWrapper
```
