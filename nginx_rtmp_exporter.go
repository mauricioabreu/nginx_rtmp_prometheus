package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"io"

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

type metrics map[int]*prometheus.Desc

func (m metrics) String() string {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	s := make([]string, len(keys))
	for i, k := range keys {
		s[i] = strconv.Itoa(k)
	}
	return strings.Join(s, ",")
}

func newMetric(metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "stream", metricName), docString, []string{}, constLabels)
}

var (
	streamMetrics = metrics{
		1: newMetric("bytes_in", "Current total of incoming bytes", nil),
		2: newMetric("bytes_out", "Current total of outgoing bytes", nil),
	}
)

// Application struct to hold all NGINX-RTMP applications
type Application struct {
	XMLName     xml.Name `xml:"application"`
	Name        string   `xml:"name"`
	LiveStreams []Live   `xml:"live"`
}

// Live struct to hold all NGINX-RTMP live streams
type Live struct {
	XMLName xml.Name `xml:"live"`
	Streams []Stream `xml:"stream"`
}

// Stream struct to hold all NGINX-RTMP streams
type Stream struct {
	XMLName  xml.Name `xml:"stream"`
	Name     string   `xml:"name"`
	BytesIn  string   `xml:"bytes_in"`
	BytesOut string   `xml:"bytes_out"`
}

// Exporter collects NGINX-RTMP stats from the status page URI
// using the prometheus metrics package
type Exporter struct {
	URI    string
	mutex  sync.RWMutex
	fetch  func() (io.ReadCloser, error)
	logger log.Logger
}

// NewExporter initializes an exporter
func NewExporter(uri string, timeout time.Duration, logger log.Logger) (*Exporter, error) {
	return &Exporter{
		URI:    uri,
		fetch:  fetchStats(uri, timeout),
		logger: logger,
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

func parseXML(respBody io.ReadCloser) (*Application, error) {
	var application Application
	data, err := ioutil.ReadAll(respBody)
	if err != nil {
		return nil, err
	}
	xml.Unmarshal(data, &application)
	return &application, nil
}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) {
	body, err := e.fetch()
	if err != nil {
		level.Error(e.logger).Log("msg", "Can't scrape NGINX-RTMP", "err", err)
		return
	}
	defer body.Close()
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range streamMetrics {
		ch <- m
	}
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
