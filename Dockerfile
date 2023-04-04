# syntax=docker/dockerfile:1

## Build
FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

# Build development version
ENV BUILD_PLATFORMS -osarch=linux/amd64


RUN go build -o /gps-tracking

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /gps-tracking /gps-tracking

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/gps-tracking"]