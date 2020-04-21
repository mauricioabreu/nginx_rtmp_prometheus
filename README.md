# NGINX-RTMP exporter for Prometheus

## Getting started

To run this exporter:

```
./nginx_rtmp_exporter [flags]
```

Help on flags:

```
./nginx_rtmp_exporter -h
```

## Building

To build the exporter:

```
make
```

or

```
make build
```

NGINX-RTMP exposes its metrics in a path specified in the nginx.conf file.

To start collecting metrics you need to update your NGINX-RTMP configuration with the following directives:

```
location /stat {
    rtmp_stat all;

    # Use this stylesheet to view XML as web page
    # in browser
    rtmp_stat_stylesheet stat.xsl;
}

location /stat.xsl {
    # XML stylesheet to view RTMP stats.
    # Copy stat.xsl wherever you want
    # and put the full directory path here
    root /path/to/stat.xsl/;
}
```

Now you can watch NGINX-RTMP statistics. This exports scrapes `localhost:8080/stats` by default you can change it for whatever address you want, juste pass `nginxrtmp.scrape-uri` argument to the exporter:

```
./nginx_rtmp_exporter --nginxrtmp.scrape-uri="localhost:9090/statistics"
```

By default the NGINX-RTMP exporter serves on port `0.0.0.0:9718` at `/metrics`