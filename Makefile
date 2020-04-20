# Ensure that 'all' is the default target otherwise it will be the first target from Makefile.common.
all::

# Needs to be defined before including Makefile.common to auto-generate targets
DOCKER_ARCHS ?= amd64 armv7 arm64 ppc64le s390x

include Makefile.common

PROMTOOL_VERSION ?= 2.5.0
PROMTOOL_URL     ?= https://github.com/prometheus/prometheus/releases/download/v$(PROMTOOL_VERSION)/prometheus-$(PROMTOOL_VERSION).$(GO_BUILD_PLATFORM).tar.gz
PROMTOOL         ?= $(FIRST_GOPATH)/bin/promtool

DOCKER_IMAGE_NAME ?= nginx-rtmp-exporter

# Use CGO for non-Linux builds.
ifeq ($(GOOS), linux)
	PROMU_CONF ?= .promu.yml
else
	ifndef GOOS
		ifeq ($(GOHOSTOS), linux)
			PROMU_CONF ?= .promu.yml
		else
			PROMU_CONF ?= .promu-cgo.yml
		endif
	else
		PROMU_CONF ?= .promu-cgo.yml
	endif
endif

PROMU := $(FIRST_GOPATH)/bin/promu --config $(PROMU_CONF)

all:: vet checkrules common-all

.PHONY: checkrules
checkrules: $(PROMTOOL)
	@echo ">> checking rules for correctness"
	find . -name "*rules*.yml" | xargs -I {} $(PROMTOOL) check rules {}

.PHONY: run-nginx-rtmp
run-nginx-rtmp:
	docker run -it -p 1935:1935 -p 8080:80 --rm alfg/nginx-rtmp

.PHONY: serve-mocked-stats
serve-mocked-stats:
	go run tests/mock_stats.go

.PHONY: ingest-stream
ingest-stream:
	docker run --net="host" --rm jrottenberg/ffmpeg -hide_banner \
		-re -f lavfi -i "testsrc2=size=1280x720:rate=30" -pix_fmt yuv420p \
		-c:v libx264 -x264opts keyint=30:min-keyint=30:scenecut=-1 \
		-tune zerolatency -profile:v high -preset veryfast -bf 0 -refs 3 \
		-b:v 1400k -bufsize 1400k \
		-utc_timing_url "https://time.akamai.com/?iso" -use_timeline 0 -media_seg_name 'chunk-stream-$RepresentationID$-$Number%05d$.m4s' \
		-init_seg_name 'init-stream1-$RepresentationID$.m4s' \
		-window_size 5  -extra_window_size 10 -remove_at_exit 1 -adaptation_sets "id=0,streams=v id=1,streams=a" -f flv rtmp://localhost:1935/stream/hello

.PHONY: promtool
promtool: $(PROMTOOL)

$(PROMTOOL):
	$(eval PROMTOOL_TMP := $(shell mktemp -d))
	curl -s -L $(PROMTOOL_URL) | tar -xvzf - -C $(PROMTOOL_TMP)
	mkdir -p $(FIRST_GOPATH)/bin
	cp $(PROMTOOL_TMP)/prometheus-$(PROMTOOL_VERSION).$(GO_BUILD_PLATFORM)/promtool $(FIRST_GOPATH)/bin/promtool
	rm -r $(PROMTOOL_TMP)