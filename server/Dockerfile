FROM golang:alpine

WORKDIR /app

# Set the working directory and print it
RUN pwd

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

ENV GOOS=linux CGO_ENABLED=0

RUN go build -o /server_binary ./cmd/api

EXPOSE 8080

CMD ["/server_binary"]
