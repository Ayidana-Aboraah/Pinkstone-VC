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

func StartCamera(Output *canvas.Image, done chan bool) string {
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

	return func() string {
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()

		time.Sleep(50 * time.Millisecond)

		for i := 0; i < 100; i++ {
			select {
			case <-done:
				return "X"
			case <-ticker.C:
				frame, release, _ := videoReader.Read()

				//Update Camera UI
				Output.Image = frame
				Output.Refresh()

				if answer := ReadImage(frame); answer != nil {
					return answer.String()
				}

				release()
				fmt.Println("Iteration: " + strconv.Itoa(i))
			}
		}
		return "X"
	}()
}
