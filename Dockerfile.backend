# builder stage
FROM golang:1.17-alpine AS builder
WORKDIR /app

# copy all the project files
COPY . .
# build the binary
RUN go build -o /app/webapi ./cmd/webapi/main.go


# final stage, the one we actually run
FROM alpine:latest
WORKDIR /app

# copy the important stuff from the builder
COPY --from=builder /app/webapi .
COPY --from=builder /app/service/api/pictureslol ./service/api/pictureslol

EXPOSE 8082
CMD ["./webapi"]