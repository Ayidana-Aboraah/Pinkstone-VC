all: test win_build clean

compile:
	go run Main.go

test:
	go test -v ./Test/...

andr_build:
	fyne package -os android -appID com.redstoneagx.BronzeHermes -icon icon02.png

win_build:
	fyne package -os windows -icon icon02.png
	./BronzeHermes.exe

clean:
	del BronzeHermes.exe
