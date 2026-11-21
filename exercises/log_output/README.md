# Log output

Generates a random hash once at startup, keeps it in memory, and prints it with a UTC timestamp every 5 seconds.

```
2026-10-18T11:20:00.123456789Z: 3f9c...e1a2
```

## Run locally

```bash
go run .
```

## Build and push the image (amd64 + arm64)

```bash
docker buildx build --platform linux/amd64,linux/arm64 \
  -t m17an/hashgen:v2 --push .
```

## Deploy to Kubernetes

```bash
kubectl create deployment log-output --image=m17an/hashgen:v2
kubectl logs -f deployment/log-output
```
