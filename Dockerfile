FROM golang:1.23.0-alpine3.20 AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY ./src ./src

# Set necessary environment variables needed for our image and build the API server.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o app ./src/main.go

FROM scratch

LABEL maintainer="Iqbal Maulana <iqbal19600@gmail.com>"

COPY --from=build etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /build/app .

ENTRYPOINT ["/app"]