#Build 
FROM golang:1.21.1-alpine3.18 as builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY .prod.env .

EXPOSE 8080
CMD [ "/app/main" ]
