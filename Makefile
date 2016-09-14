go: powerbot
	./powerbot

powerbot: powerbot.go
	go build .

clean:
	rm -f powerbot

test:
	go test -v

get-deps:
	go get .
