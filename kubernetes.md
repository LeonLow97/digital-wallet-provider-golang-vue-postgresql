# Kubernetes

- Running k8s locally (single node) with Docker Desktop Kubernetes

## Pushing images to Docker Registry (Docker Hub)

```
cd server
docker login -u <username>
docker build -t lowjiewei/mobile-wallet-be:v1.1 .
docker push lowjiewei/mobile-wallet-be:v1.1
```

## Commands to setup k8s cluster

```
kubectl apply -f ./k8s/
```

## Creating secrets in k8s cluster

- Creating secrets in kubernetes cluster is an **imperative** command. I have created one to store Postgres Password.

```
kubectl create secret generic pgpassword --from-literal PGPASSWORD=verystrongpassword
kubectl create secret generic pguser --from-literal PGUSER=postgres
kubectl create secret generic pghost --from-literal PGHOST=postgres-cluster-ip-service
kubectl create secret generic pgdb --from-literal PGDB=postgres
kubectl create secret generic redishost --from-literal REDISHOST=redis-cluster-ip-service

kubectl get secrets
NAME         TYPE     DATA   AGE
pgpassword   Opaque   1      18s
```
