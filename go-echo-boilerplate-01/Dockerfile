# BUILD STAGE
FROM golang:alpine as builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY ./go.mod ./go.sum ./
COPY . .
RUN go mod tidy
RUN go build -o /app/my-service /app/cmd/main.go

# Run Stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/my-service .
COPY --from=builder /app/configs/env.yaml /app/configs/env.yaml

EXPOSE 8100
CMD ["/app/my-service"]