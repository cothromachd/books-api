FROM golang:latest

RUN mkdir /app
ADD . /app/
WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o main .

CMD ["./main", "-dsn", "-ra", ]