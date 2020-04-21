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

## Collectors

This exporter collects and exposes the following statistics:

```
# HELP nginx_rtmp_exporter_build_info A metric with a constant '1' value labeled by version, revision, branch, and goversion from which nginx_rtmp_exporter was built.
# TYPE nginx_rtmp_exporter_build_info gauge
nginx_rtmp_exporter_build_info{branch="master",goversion="go1.14",revision="83760dedcd355d0dca52a635751b87982bf189e8",version="0.0.1"} 1
# HELP nginx_rtmp_stream_incoming_bytes_total Current total of incoming bytes
# TYPE nginx_rtmp_stream_incoming_bytes_total counter
nginx_rtmp_stream_incoming_bytes_total{stream="hello"} 2.12292e+06
nginx_rtmp_stream_incoming_bytes_total{stream="hello_240p264kbs"} 275959
nginx_rtmp_stream_incoming_bytes_total{stream="hello_240p528kbs"} 539485
nginx_rtmp_stream_incoming_bytes_total{stream="hello_360p878kbs"} 996627
nginx_rtmp_stream_incoming_bytes_total{stream="hello_480p1128kbs"} 1.324909e+06
nginx_rtmp_stream_incoming_bytes_total{stream="hello_720p2628kbs"} 3.306701e+06
# HELP nginx_rtmp_stream_outgoing_bytes_total Current total of outgoing bytes
# TYPE nginx_rtmp_stream_outgoing_bytes_total counter
nginx_rtmp_stream_outgoing_bytes_total{stream="hello"} 1.938764e+06
nginx_rtmp_stream_outgoing_bytes_total{stream="hello_240p264kbs"} 0
nginx_rtmp_stream_outgoing_bytes_total{stream="hello_240p528kbs"} 0
nginx_rtmp_stream_outgoing_bytes_total{stream="hello_360p878kbs"} 0
nginx_rtmp_stream_outgoing_bytes_total{stream="hello_480p1128kbs"} 0
nginx_rtmp_stream_outgoing_bytes_total{stream="hello_720p2628kbs"} 0
# HELP nginx_rtmp_stream_receive_bytes_per_second Current bandwidth in per second
# TYPE nginx_rtmp_stream_receive_bytes_per_second gauge
nginx_rtmp_stream_receive_bytes_per_second{stream="hello"} 1.5095291137695312
nginx_rtmp_stream_receive_bytes_per_second{stream="hello_240p264kbs"} 0
nginx_rtmp_stream_receive_bytes_per_second{stream="hello_240p528kbs"} 0
nginx_rtmp_stream_receive_bytes_per_second{stream="hello_360p878kbs"} 0
nginx_rtmp_stream_receive_bytes_per_second{stream="hello_480p1128kbs"} 0
nginx_rtmp_stream_receive_bytes_per_second{stream="hello_720p2628kbs"} 0
# HELP nginx_rtmp_stream_transmit_bytes_per_second Current bandwidth out per second
# TYPE nginx_rtmp_stream_transmit_bytes_per_second gauge
nginx_rtmp_stream_transmit_bytes_per_second{stream="hello"} 1.3690261840820312
nginx_rtmp_stream_transmit_bytes_per_second{stream="hello_240p264kbs"} 0
nginx_rtmp_stream_transmit_bytes_per_second{stream="hello_240p528kbs"} 0
nginx_rtmp_stream_transmit_bytes_per_second{stream="hello_360p878kbs"} 0
nginx_rtmp_stream_transmit_bytes_per_second{stream="hello_480p1128kbs"} 0
nginx_rtmp_stream_transmit_bytes_per_second{stream="hello_720p2628kbs"} 0
# HELP nginx_rtmp_stream_uptime Number of seconds since the stream started
# TYPE nginx_rtmp_stream_uptime counter
nginx_rtmp_stream_uptime{stream="hello"} 11.724
nginx_rtmp_stream_uptime{stream="hello_240p264kbs"} 5.44
nginx_rtmp_stream_uptime{stream="hello_240p528kbs"} 5.5
nginx_rtmp_stream_uptime{stream="hello_360p878kbs"} 5.55
nginx_rtmp_stream_uptime{stream="hello_480p1128kbs"} 5.61
nginx_rtmp_stream_uptime{stream="hello_720p2628kbs"} 5.66
```