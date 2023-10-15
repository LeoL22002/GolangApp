FROM golang:latest
COPY . /go/src/app
WORKDIR /go/src/app
RUN go build -o miapp
CMD ["./miapp"]