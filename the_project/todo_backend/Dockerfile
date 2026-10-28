FROM --platform=$BUILDPLATFORM golang:1.27 AS build
ARG TARGETOS TARGETARCH
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /todo_backend ./cmd/web

FROM gcr.io/distroless/static
COPY --from=build /todo_backend /todo_backend

ENV PORT=8080
EXPOSE 8080

CMD ["/todo_backend"]
