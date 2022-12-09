FROM golang:1.19-alpine

RUN addgroup -g 777 datasink && adduser -u 777 -S -G datasink datasink

COPY datasink-docker /bin/datasink

EXPOSE 8080 1883

USER datasink:datasink
