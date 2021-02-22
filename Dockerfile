FROM golang:1.15-alpine

WORKDIR /build

COPY . .
RUN go build -o main .

EXPOSE 8080

CMD ["/build/main"]