all: test win_build run

andr_build:
	fyne package -os android -appID com.redstoneagx.BronzeHermes -icon icon02.png

win_build:
	fyne package -os windows -icon icon02.png

test:
	go test -v ./Test/...

run:
	./BronzeHermes.exe

fun:
	go run Main.go