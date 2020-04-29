ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:glibc
LABEL maintainer="Maurício Antunes <mauricio.abreua@gmail.com>"

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/nginx_rtmp_exporter /bin/nginx_rtmp_exporter

EXPOSE 9728
USER nobody
ENTRYPOINT [ "/bin/nginx_rtmp_exporter" ]