# Todo backend

Go HTTP service that stores the project's todo items. The todos are kept in memory, so they are lost when the pod restarts. It listens on the port set by the `PORT` environment variable (`8080` in the image).

Only the todo app talks to it, through the `todo-backend-svc` Service. It is not exposed through the ingress.

Endpoints:

- `GET /todos` returns the todos as a JSON array of strings
- `POST /todos` creates a todo from the `todo-input` form field (1 to 140 characters): `201 Created`, or `400` with the reason
- `GET /api/healthz` returns `200 OK`

## Run locally

```bash
PORT=8080 go run ./cmd/web
curl -X POST localhost:8080/todos --data-urlencode "todo-input=buy milk"
curl localhost:8080/todos
```

## Build and push the image

```bash
docker buildx build -t m17an/todo_backend:<tag> --push .
```

## Deploy to Kubernetes

```bash
kubectl apply -f manifests/
kubectl logs -f deployment/todo-backend
```
