.PHONY: all build install uninstall

all: build

build:
	go build .

install: build
	@if [ -f /usr/local/bin/synchronex ]; then \
		sudo rm /usr/local/bin/synchronex; \
	fi
	sudo cp ./synchronex /usr/local/bin/synchronex
	sudo chown $(USER):$(USER) /usr/local/bin/synchronex
	sudo chmod +x /usr/local/bin/synchronex

uninstall:
	@if [ -f /usr/local/bin/synchronex ]; then \
		sudo rm /usr/local/bin/synchronex; \
	fi

clean:
	rm -f synchronex

test:
	go test ./...