go: powerbot
	./powerbot

powerbot: powerbot.go
	go build .

clean:
	rm -f powerbot
