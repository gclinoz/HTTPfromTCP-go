.PHONY: all clean

all:
	go test ./...

clean:
	rm ./udpsender
