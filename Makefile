.PHONY: all clean

all:
	go build ./cmd/udpsender/

clean:
	rm ./udpsender
