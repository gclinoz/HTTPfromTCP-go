.PHONY: all clean

all:
	go build ./cmd/tcplistener/

clean:
	rm ./tcplistener
