target = xftp
version=0.0.1
json=info.json
release = $(shell git rev-parse --short HEAD)

.PHONY: rpm

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

clean:
	rm -f $(target)
	rm -f $(json)
	rm -f *.rpm

rpm: clean
	CGO_ENABLED=0 GOOS=linux go build -o $(target)
	python replace.py $(release)
	go-bin-rpm generate -f $(json) --version $(version) -o xftp-$(version)-$(release).x86_64.rpm
