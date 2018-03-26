# Docker builder for Golang
FROM golang AS builder
LABEL maintainer "Thomas Bouvier <tomatrocho@gmail.com>"

RUN wget https://bootstrap.pypa.io/get-pip.py
RUN python2.7 get-pip.py
RUN pip install colorthief

RUN apt-get update ; apt-get install -y uuid-runtime

WORKDIR /go/src/app
COPY ./src .
RUN set -x && \
    go get -d -v . && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Docker run Golang app
FROM alpine
LABEL maintainer "Thomas Bouvier <tomatrocho@gmail.com>"

WORKDIR /root/
RUN mkdir img

COPY --from=builder /go/src/app .

EXPOSE 9000

CMD ["./app"]
