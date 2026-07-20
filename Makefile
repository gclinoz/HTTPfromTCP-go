.PHONY: all clean

all:
	go build ./cmd/httpserver/

clean:
	rm ./httpserver
