run-nginx-rtmp:
	docker run -it -p 1935:1935 -p 8080:80 --rm alfg/nginx-rtmp

serve-mocked-stats:
	go run tests/mock_stats.go

ingest-stream:
	docker run --net="host" --rm jrottenberg/ffmpeg -hide_banner \
		-re -f lavfi -i "testsrc2=size=1280x720:rate=30" -pix_fmt yuv420p \
		-c:v libx264 -x264opts keyint=30:min-keyint=30:scenecut=-1 \
		-tune zerolatency -profile:v high -preset veryfast -bf 0 -refs 3 \
		-b:v 1400k -bufsize 1400k \
		-utc_timing_url "https://time.akamai.com/?iso" -use_timeline 0 -media_seg_name 'chunk-stream-$RepresentationID$-$Number%05d$.m4s' \
		-init_seg_name 'init-stream1-$RepresentationID$.m4s' \
		-window_size 5  -extra_window_size 10 -remove_at_exit 1 -adaptation_sets "id=0,streams=v id=1,streams=a" -f flv rtmp://localhost:1935/stream/hello