# Docker builder for Golang
FROM golang as builder
LABEL maintainer "Thomas Bouvier <tomatrocho@gmail.com>"

RUN wget https://bootstrap.pypa.io/get-pip.py
RUN python2.7 get-pip.py
RUN pip install colorthief

RUN apt-get update ; apt-get install -y uuid-runtime

WORKDIR /go/src/github.com/tomatrocho/insapp-go
COPY ./src .
RUN set -x && \
    go get -d -v . && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Docker run Golang app
FROM scratch
LABEL maintainer "Thomas Bouvier <tomatrocho@gmail.com>"

WORKDIR /root/
COPY --from=0 /go/src/github.com/tomatrocho/insapp-go .

EXPOSE 9000

CMD ["./app"]
