FROM golang:1.21.0-alpine3.18

WORKDIR /home

COPY ./pkg /home

RUN cd /home && go build -o library

CMD ["/home/library"]


