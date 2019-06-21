GOFLAGS=-ldflags="-s -w"

build: less assets
	go build $(GOFLAGS) .

assets:
	go-assets-builder static/shared.css static/favicon.ico views -o assets.go 

run:
	sudo O2=dev go run o2.go

less:
	lessc static/index.less static/shared.css

less-watch:
	ls shared/*.less | entr lessc static/index.less static/shared.css