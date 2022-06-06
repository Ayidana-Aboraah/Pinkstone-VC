package Cam

import (
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

	return func() string {
		timer := time.NewTimer(10 * time.Second)
		defer timer.Stop()

		time.Sleep(50 * time.Millisecond)

		reader := oned.NewUPCAReader()
		var bmp *gozxing.BinaryBitmap
		var result *gozxing.Result

		// var err error

		for {
			select {
			case <-timer.C:
				return "X"
			case <-done:
				return ""
			default:
				frame, release, _ := videoReader.Read()

				Output.Image = frame //Update Camera UI
				Output.Refresh()

				bmp, _ = gozxing.NewBinaryBitmapFromImage(frame)

				release()

				result, _ = reader.Decode(bmp, map[gozxing.DecodeHintType]interface{}{gozxing.DecodeHintType_TRY_HARDER: true})

				if result != nil {
					return result.String()
				}
			}
		}
	}()
}
