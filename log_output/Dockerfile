FROM --platform=$BUILDPLATFORM golang:1.27 AS build
ARG TARGETOS TARGETARCH
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /todo_app

FROM gcr.io/distroless/static
COPY --from=build /todo_app /todo_app

ENV PORT=8080

CMD ["/todo_app"]
