# syntax=docker/dockerfile:1

# Build
FROM golang:alpine3.17 AS build
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
# Build development version
ENV BUILD_PLATFORMS -osarch=linux/amd64

EXPOSE 8000

RUN go build -o /mct

# Deploy
FROM alpine:latest
WORKDIR /app

EXPOSE 8000

COPY local.json .
COPY templates .
COPY --from=build /mct .

ENTRYPOINT ["/app/mct"]
