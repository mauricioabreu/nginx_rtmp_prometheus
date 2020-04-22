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
nginx_rtmp_exporter_build_info{branch="master",goversion="go1.14",revision="fe3d8ac350cec520648b07cf9ceb613f12362e2b",version="0.0.1"} 1
# HELP nginx_rtmp_server_incoming_bytes_total Current total of incoming bytes
# TYPE nginx_rtmp_server_incoming_bytes_total counter
nginx_rtmp_server_incoming_bytes_total 2.40895361e+08
# HELP nginx_rtmp_server_outgoing_bytes_total Current total of outgoing bytes
# TYPE nginx_rtmp_server_outgoing_bytes_total counter
nginx_rtmp_server_outgoing_bytes_total 5.3739661e+07
# HELP nginx_rtmp_server_receive_bytes Current bandwidth in per second
# TYPE nginx_rtmp_server_receive_bytes gauge
nginx_rtmp_server_receive_bytes 6.592750549316406
# HELP nginx_rtmp_server_transmit_bytes Current bandwidth out per second
# TYPE nginx_rtmp_server_transmit_bytes gauge
nginx_rtmp_server_transmit_bytes 1.4704513549804688
# HELP nginx_rtmp_server_uptime_seconds_total Number of seconds NGINX-RTMP started
# TYPE nginx_rtmp_server_uptime_seconds_total counter
nginx_rtmp_server_uptime_seconds_total 310
# HELP nginx_rtmp_stream_incoming_bytes_total Current total of incoming bytes
# TYPE nginx_rtmp_stream_incoming_bytes_total counter
nginx_rtmp_stream_incoming_bytes_total{stream="hello"} 5.3836313e+07
nginx_rtmp_stream_incoming_bytes_total{stream="hello_240p264kbs"} 7.772386e+06
nginx_rtmp_stream_incoming_bytes_total{stream="hello_240p528kbs"} 1.5462727e+07
nginx_rtmp_stream_incoming_bytes_total{stream="hello_360p878kbs"} 2.8888155e+07
nginx_rtmp_stream_incoming_bytes_total{stream="hello_480p1128kbs"} 3.849418e+07
nginx_rtmp_stream_incoming_bytes_total{stream="hello_720p2628kbs"} 9.5998269e+07
# HELP nginx_rtmp_stream_outgoing_bytes_total Current total of outgoing bytes
# TYPE nginx_rtmp_stream_outgoing_bytes_total counter
nginx_rtmp_stream_outgoing_bytes_total{stream="hello"} 5.365292e+07
nginx_rtmp_stream_outgoing_bytes_total{stream="hello_240p264kbs"} 0
nginx_rtmp_stream_outgoing_bytes_total{stream="hello_240p528kbs"} 0
nginx_rtmp_stream_outgoing_bytes_total{stream="hello_360p878kbs"} 0
nginx_rtmp_stream_outgoing_bytes_total{stream="hello_480p1128kbs"} 0
nginx_rtmp_stream_outgoing_bytes_total{stream="hello_720p2628kbs"} 0
# HELP nginx_rtmp_stream_receive_bytes Current bandwidth in per second
# TYPE nginx_rtmp_stream_receive_bytes gauge
nginx_rtmp_stream_receive_bytes{stream="hello"} 1.468170166015625
nginx_rtmp_stream_receive_bytes{stream="hello_240p264kbs"} 0.20711517333984375
nginx_rtmp_stream_receive_bytes{stream="hello_240p528kbs"} 0.4146881103515625
nginx_rtmp_stream_receive_bytes{stream="hello_360p878kbs"} 0.77984619140625
nginx_rtmp_stream_receive_bytes{stream="hello_480p1128kbs"} 1.0345001220703125
nginx_rtmp_stream_receive_bytes{stream="hello_720p2628kbs"} 2.619903564453125
# HELP nginx_rtmp_stream_transmit_bytes Current bandwidth out per second
# TYPE nginx_rtmp_stream_transmit_bytes gauge
nginx_rtmp_stream_transmit_bytes{stream="hello"} 1.468170166015625
nginx_rtmp_stream_transmit_bytes{stream="hello_240p264kbs"} 0
nginx_rtmp_stream_transmit_bytes{stream="hello_240p528kbs"} 0
nginx_rtmp_stream_transmit_bytes{stream="hello_360p878kbs"} 0
nginx_rtmp_stream_transmit_bytes{stream="hello_480p1128kbs"} 0
nginx_rtmp_stream_transmit_bytes{stream="hello_720p2628kbs"} 0
# HELP nginx_rtmp_stream_uptime_seconds_total Number of seconds since the stream started
# TYPE nginx_rtmp_stream_uptime_seconds_total counter
nginx_rtmp_stream_uptime_seconds_total{stream="hello"} 306.805
nginx_rtmp_stream_uptime_seconds_total{stream="hello_240p264kbs"} 300.55
nginx_rtmp_stream_uptime_seconds_total{stream="hello_240p528kbs"} 300.6
nginx_rtmp_stream_uptime_seconds_total{stream="hello_360p878kbs"} 300.65
nginx_rtmp_stream_uptime_seconds_total{stream="hello_480p1128kbs"} 300.7
nginx_rtmp_stream_uptime_seconds_total{stream="hello_720p2628kbs"} 300.75
```