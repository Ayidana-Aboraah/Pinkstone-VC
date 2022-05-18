package Cam

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/canvas"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
	"github.com/pion/mediadevices"
	_ "github.com/pion/mediadevices/pkg/driver/camera"
	"github.com/pion/mediadevices/pkg/prop"
)

func StartCamera(Output *canvas.Image, done chan bool) string {
	stream, _ := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Video: func(constraint *mediadevices.MediaTrackConstraints) {
			constraint.FrameRate = prop.Float(24)
		}})

	vidTrack := stream.GetVideoTracks()[0]
	videoTrack := vidTrack.(*mediadevices.VideoTrack)
	videoReader := videoTrack.NewReader(false)

	defer videoTrack.Close()

	reader := oned.NewUPCAReader()

	return func() string {
		timer := time.NewTimer(10 * time.Second)
		defer timer.Stop()

		time.Sleep(50 * time.Millisecond)

		var bmp *gozxing.BinaryBitmap
		var result *gozxing.Result
		var err error

		for {
			select {
			case <-timer.C:
				return "X"
			case <-done:
				return ""
			default:
				frame, release, _ := videoReader.Read()

				//Update Camera UI
				Output.Image = frame
				Output.Refresh()

				// bmp, _ = gozxing.NewBinaryBitmapFromImage(frame)
				bmp, err = gozxing.NewBinaryBitmapFromImage(frame)
				if err != nil {
					fmt.Println(err)
				}

				result, err = reader.Decode(bmp, nil)
				if err != nil {
					fmt.Println(err)
				}

				if result != nil {
					return result.String()
				}

				release()
			}
		}
	}()
}
