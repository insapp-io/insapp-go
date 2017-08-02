FROM golang

MAINTAINER ftm VERSION 1.0

RUN wget https://bootstrap.pypa.io/get-pip.py
RUN python2.7 get-pip.py
RUN pip install colorthief

RUN apt-get update ; apt-get install -y uuid-runtime

EXPOSE 9000

RUN mkdir /go/src/app
COPY ./ /go/src/app/
RUN cd /go/src/app/src && go get

ENTRYPOINT cd /go/src/app/src && go run *.go