package Cam

import (
	"BronzeHermes/UI"
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2/canvas"
	"github.com/pion/mediadevices"
	_ "github.com/pion/mediadevices/pkg/driver/camera"
	"github.com/pion/mediadevices/pkg/prop"
)

func StartTestCamera(Output *canvas.Image, done chan bool) string {
	stream, errA := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Video: func(constraint *mediadevices.MediaTrackConstraints) {
			constraint.FrameRate = prop.Float(24)
		},
	})
	UI.HandleError(errA)

	vidTrack := stream.GetVideoTracks()[0]
	videoTrack := vidTrack.(*mediadevices.VideoTrack)
	defer videoTrack.Close()

	videoReader := videoTrack.NewReader(false)

	result := ""
	func() chan bool {
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()

		time.Sleep(50 * time.Millisecond)

		i := 0
		for {
			select {
			case <-done:
				time.Sleep(50 * time.Millisecond)
				return done
			case <-ticker.C:
				i += 1
				frame, release, _ := videoReader.Read()

				//Update Camera UI
				Output.Image = frame
				Output.Refresh()

				answer := ReadImage(frame)
				if answer != nil {
					result = answer.String()
					Output.FillMode = canvas.ImageFillStretch
					release()
					return done
				}

				if i >= 250 {
					result = "X"
					Output.FillMode = canvas.ImageFillStretch
					release()
					return done
				}

				fmt.Println("Iteration: " + strconv.Itoa(i))
				release()
			}
		}
	}()

	fmt.Println("Complete")
	return result
}
