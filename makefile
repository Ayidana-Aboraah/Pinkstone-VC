all: test win clean

compile:
	go run Main.go

test:
	go test -v ./Test/...

win:
	sudo systemctl start docker
	sudo fyne-cross windows -icon="icon02.png" -app-id="Bronze.Hermes" -name="Pinkstone"
#	env GOOS=windows GOARCH=amd64
#	go build -o BH.exe Main.go CGO=ENABLED 

andr_build:
	fyne package -os android -appID com.redstoneagx.BronzeHermes -icon icon02.png

win_build:
	fyne package -os windows -icon icon02.png
	./BronzeHermes.exe

clean:
	del BronzeHermes.exe
