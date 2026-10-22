FROM --platform=$BUILDPLATFORM golang:1.27 AS build
ARG TARGETOS TARGETARCH
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /log_output

FROM gcr.io/distroless/static
COPY --from=build /log_output /log_output

# No PORT here: without PORT the binary runs as the writer, with PORT as the reader.
ENTRYPOINT ["/log_output"]
