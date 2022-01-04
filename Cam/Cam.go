package Cam

import (
	"BronzeHermes/UI"
	"fmt"
	"fyne.io/fyne/v2/canvas"
	"github.com/pion/mediadevices"
	_ "github.com/pion/mediadevices/pkg/driver/camera"
	"github.com/pion/mediadevices/pkg/prop"
	"image"
	"strconv"
	"time"
)

func StartCamera() string{
	stream, errA:= mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Video: func(constraint *mediadevices.MediaTrackConstraints) {
			constraint.Width = prop.Int(600)
			constraint.Height = prop.Int(400)
			constraint.FrameRate = prop.Float(24)
		},
	})
	UI.HandleError(errA)


	vidTrack := stream.GetVideoTracks()[0]
	videoTrack := vidTrack.(*mediadevices.VideoTrack)
	defer videoTrack.Close()

	videoReader := videoTrack.NewReader(false)

	result := ""
	run := func() chan bool{
		ticker := time.NewTicker(10 * time.Millisecond)
		done := make(chan bool, 1)
		defer close(done)
		defer ticker.Stop()

		time.Sleep(50 * time.Millisecond)

		i := 0
		for{
			select{
			case <- done:
				return done
			case <- ticker.C:
				i += 1
				frame, release, _ := videoReader.Read()
				RefreshCam(frame)

				answer := ReadImage(frame)
				if answer != nil {
					result = answer.String()
					CamOutput.FillMode = canvas.ImageFillStretch
					CamOutput.Refresh()
					release()
					return done
				}

				if i >= 250{
					result = "X"
					CamOutput.FillMode = canvas.ImageFillStretch
					CamOutput.Refresh()
					release()
					return done
				}

				fmt.Println("Iteration: " + strconv.Itoa(i))
				release()
			}
		}
	}
	run()

	fmt.Println("Complete")
	return result
}

var CamOutput *canvas.Image

func RefreshCam(newInput image.Image){
	CamOutput.Image = newInput
	CamOutput.Refresh()
}

