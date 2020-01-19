package screen

import (
	"errors"
	"image"

	"github.com/kbinani/screenshot"
)

func Get(im *image.RGBA) error {

	// Capture each displays.
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		return errors.New("Active display not found")
	}

	// var all image.Rectangle = image.Rect(0, 0, 0, 0)

	// for i := 0; i < n; i++ {
	// 	bounds := screenshot.GetDisplayBounds(i)
	// 	all = bounds.Union(all)
	//
	// 	img, err := screenshot.CaptureRect(bounds)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fileName := fmt.Sprintf("%d_%dx%d.png", i, bounds.Dx(), bounds.Dy())
	// 	save(img, fileName)
	//
	// 	fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fileName)
	// }

	// Capture all desktop region into an image.
	// fmt.Printf("%v\n", im)
	b := im.Bounds()
	simg, err := screenshot.Capture(b.Min.X, b.Min.Y, b.Dx(), b.Dy())
	if err != nil {
		return err
	}

	// draw.Draw(all, b, img, b.Min, draw.Src)
	copy(im.Pix, simg.Pix)
	return nil
}
