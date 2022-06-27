all: test win_build

andr_build:
	fyne package -os android -appID com.redstoneagx.BronzeHermes -icon icon02.png

win_build:
	fyne package -os windows -icon icon02.png
	./BronzeHermes.exe

test:
	go test -v ./Test/...

compile:
	go run Main.go