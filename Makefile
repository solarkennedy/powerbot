go: powerbot powerbot-cli
	./powerbot-cli 5330371
	sleep 2s
	./powerbot-cli 5330380
	./powerbot

powerbot: powerbot.go cmd/powerbot/powerbot.go
	go build cmd/powerbot/powerbot.go

powerbot-cli: powerbot.go cmd/powerbot-cli/powerbot-cli.go
	go build cmd/powerbot-cli/powerbot-cli.go

clean:
	rm -f powerbot powerbot-cli

test:
	go test -v

get-deps:
	go get .
