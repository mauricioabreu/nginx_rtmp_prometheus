// MIT License

// Copyright (c) 2020 Mauricio Antunes

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"io"

	"github.com/antchfx/xmlquery"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	namespace = "nginx_rtmp"
)

func newMetric(metricName string, docString string, varLabels []string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "stream", metricName), docString, varLabels, constLabels)
}

// Exporter collects NGINX-RTMP stats from the status page URI
// using the prometheus metrics package
type Exporter struct {
	URI    string
	mutex  sync.RWMutex
	fetch  func() (io.ReadCloser, error)
	logger log.Logger

	bytesIn  *prometheus.Desc
	bytesOut *prometheus.Desc
}

// StreamInfo characteristics of a stream
type StreamInfo struct {
	Name     string
	BytesIn  float64
	BytesOut float64
}

// NewStreamInfo builds a StreamInfo struct from string values
func NewStreamInfo(name string, bytesIn string, bytesOut string) StreamInfo {
	var bytesInInt, bytesOutInt float64
	if n, err := strconv.ParseFloat(bytesIn, 64); err == nil {
		bytesInInt = n
	}
	if n, err := strconv.ParseFloat(bytesOut, 64); err == nil {
		bytesOutInt = n
	}
	return StreamInfo{Name: name, BytesIn: bytesInInt, BytesOut: bytesOutInt}
}

// NewExporter initializes an exporter
func NewExporter(uri string, timeout time.Duration, logger log.Logger) (*Exporter, error) {
	return &Exporter{
		URI:      uri,
		fetch:    fetchStats(uri, timeout),
		logger:   logger,
		bytesIn:  newMetric("bytes_in", "Current total of incoming bytes", []string{"stream"}, nil),
		bytesOut: newMetric("bytes_out", "Current total of outgoing bytes", []string{"stream"}, nil),
	}, nil
}

func fetchStats(uri string, timeout time.Duration) func() (io.ReadCloser, error) {
	client := http.Client{
		Timeout: timeout,
	}

	return func() (io.ReadCloser, error) {
		resp, err := client.Get(uri)
		if err != nil {
			return nil, err
		}
		if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
			resp.Body.Close()
			return nil, fmt.Errorf("HTTP status %d", resp.StatusCode)
		}
		return resp.Body, nil
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect from concurrent collects
	defer e.mutex.Unlock()

	e.scrape(ch)
}

func parseStreams(respBody io.ReadCloser) ([]StreamInfo, error) {
	doc, err := xmlquery.Parse(respBody)
	if err != nil {
		return nil, err
	}

	streams := make([]StreamInfo, 0)
	data := xmlquery.Find(doc, "//stream")

	for _, stream := range data {
		name := stream.SelectElement("name").InnerText()
		bytesIn := stream.SelectElement("bytes_in").InnerText()
		bytesOut := stream.SelectElement("bytes_out").InnerText()
		streams = append(streams, NewStreamInfo(name, bytesIn, bytesOut))
	}
	return streams, nil
}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) {
	data, err := e.fetch()
	if err != nil {
		level.Error(e.logger).Log("msg", "Can't scrape NGINX-RTMP", "err", err)
		return
	}
	defer data.Close()

	streams, err := parseStreams(data)
	if err != nil {
		level.Error(e.logger).Log("msg", "Can't parse XML", "err", err)
		return
	}

	for _, stream := range streams {
		ch <- prometheus.MustNewConstMetric(e.bytesIn, prometheus.CounterValue, stream.BytesIn, stream.Name)
		ch <- prometheus.MustNewConstMetric(e.bytesOut, prometheus.CounterValue, stream.BytesOut, stream.Name)
	}
}

// Describe describes all metrics to be exported to Prometheus
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.bytesIn
	ch <- e.bytesOut
}

func main() {
	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9718").String()
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
		scrapeURI     = kingpin.Flag("nginxrtmp.scrape-uri", "URI on which to scrape HAProxy.").Default("http://localhost:8080/stats").String()
		timeout       = kingpin.Flag("nginxrtmp.timeout", "Timeout for trying to get stats from HAProxy.").Default("5s").Duration()
	)

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)
	level.Info(logger).Log("msg", "Starting nginx_rtmp_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "context", version.BuildContext())

	exporter, err := NewExporter(*scrapeURI, *timeout, logger)
	if err != nil {
		level.Error(logger).Log("msg", "Error creating an exporter", "err", err)
		os.Exit(1)
	}
	prometheus.MustRegister(exporter)
	prometheus.MustRegister(version.NewCollector("nginx_rtmp_exporter"))

	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>NGINX-RTMP exporter</title></head>
			<body>
			<h1>NGINX-RTMP exporter</h1>
			<p><a href='` + *metricsPath + `'>Metrics</a></p>
			</body>
			</html>`,
		))
	})
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
