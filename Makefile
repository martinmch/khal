.PHONY: all clean

all: khal khal.1

khal: khal.go
	go build -o khal khal.go

khal.1: khal.1.md
	pandoc -s -t man khal.1.md -o khal.1

clean:
	rm khal.1 khal
