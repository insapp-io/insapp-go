# Docker builder for Golang
FROM golang AS builder
LABEL maintainer "Thomas Bouvier <tomatrocho@gmail.com>"

RUN apt-get update ; apt-get install -y uuid-runtime

WORKDIR /go/src/app
COPY ./src .
RUN set -x && \
    go get -d -v . && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Docker run Golang app
FROM alpine
LABEL maintainer "Thomas Bouvier <tomatrocho@gmail.com>"

RUN apk update && \
    apk add --no-cache bash && \
    apk add ca-certificates && \
    apk add build-base python-dev jpeg-dev zlib-dev

ENV LIBRARY_PATH=/lib:/usr/lib

RUN python -m ensurepip && \
    rm -r /usr/lib/python*/ensurepip && \
    pip install --upgrade pip setuptools && \
    pip install colorthief

RUN rm -rf /var/cache/apk/* && \
    rm -r /root/.cache

WORKDIR /root/
COPY --from=builder /go/src/app .

EXPOSE 9000

CMD ["./app"]
