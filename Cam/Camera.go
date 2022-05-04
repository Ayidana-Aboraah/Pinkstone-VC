package Cam

import (
	"fmt"
	"strconv"
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

	if stream.GetVideoTracks() == nil || len(stream.GetVideoTracks()) == 0 {
		return "E"
	}

	vidTrack := stream.GetVideoTracks()[0]
	videoTrack := vidTrack.(*mediadevices.VideoTrack)
	videoReader := videoTrack.NewReader(false)

	defer vidTrack.Close()

	reader := oned.NewUPCAReader()
	var bmp *gozxing.BinaryBitmap
	var result *gozxing.Result

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

				bmp, _ = gozxing.NewBinaryBitmapFromImage(frame)
				// bmp, err := gozxing.NewBinaryBitmapFromImage(frame)
				// if err != nil{
				// 	fmt.Println(err)
				// }

				result, _ = reader.Decode(bmp, nil)
				// result, err = reader.Decode(bmp, nil)
				// if err != nil{
				// 	fmt.Println(err)
				// }

				if result != nil {
					return result.String()
				}

				release()
				fmt.Println("Iteration: " + strconv.Itoa(i)) //Remove after debugging
			}
		}
		return "X"
	}()
}
