default: build

build: fix
	go build -v .

fix: *.go
	goimports -l -w .
	gofmt -l -w .

fire:
	./sv4-codec.exe "C:\GOG Games\RollerCoaster Tycoon Deluxe (German)\Saved Games\Dynamite Dunes 2.SV4"

run: build fire
