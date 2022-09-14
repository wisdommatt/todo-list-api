FROM golang:1.19.1-alpine
WORKDIR /app
COPY ./ /app
RUN go mod download
EXPOSE 8080
ENTRYPOINT ["sh", "-c", "go run main.go"]