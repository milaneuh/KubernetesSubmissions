# Todo app

Go HTTP server for the course project. It listens on the port set by the `PORT` environment variable (`8080` in the image) and prints `Server started in port <PORT>` at startup.

Endpoints:

- `GET /healthz` returns `200 OK`

## Run locally

```bash
PORT=8080 go run .
curl localhost:8080/healthz
```

## Build and push the image (amd64 + arm64)

```bash
docker buildx build --platform linux/amd64,linux/arm64 \
  -t <dockerhub-user>/todo-app:<tag> --push .
```

## Deploy to Kubernetes

```bash
kubectl create deployment todo-app --image=<dockerhub-user>/todo-app:<tag>
kubectl logs -f deployment/todo-app
```
