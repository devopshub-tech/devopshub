# build app
FROM golang:1.21.0-alpine3.18 as base_build

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN go build -v -o worker-svc cmd/worker/main.go

# create runtime
FROM alpine:3.18.3

WORKDIR /app

RUN apk add ca-certificates

COPY --from=base_build /app/worker-svc .

EXPOSE 8000

ENTRYPOINT ["/app/worker-svc"]