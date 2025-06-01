FROM golang:1.24-alpine AS build

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app

FROM alpine:latest AS certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=build /build /
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 9717
ENTRYPOINT [ "/app" ]
