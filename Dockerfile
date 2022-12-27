FROM golang:1.19 as builder
LABEL maintainer="Maurício Antunes <mauricio.abreua@gmail.com>"
WORKDIR /usr/src
COPY . .
RUN go build -v -o /usr/src/nginx_rtmp_exporter

CMD [ "/usr/src/nginx_rtmp_exporter" ]
