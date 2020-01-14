package main

import (
	"flag"
	"image"
	"math"
	"net"
	"os"
	"runtime/pprof"
	"time"

	"github.com/bradfitz/rfbgo/rfb"
	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
	log "github.com/sirupsen/logrus"
)

type mouseS struct {
	clickFlag bool
	x         uint16
	y         uint16
}

var clickMap = map[uint8]string{
	1:  "left",
	2:  "center",
	4:  "right",
	8:  "scrollup",
	16: "scrolldown",
}

var keyboardMap = map[uint32]string{
	65288: "delete",
}

func main() {
	serve()
}

func getScreen(all *image.RGBA) {
	// Capture each displays.
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		panic("Active display not found")
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
	// fmt.Printf("%v\n", all)
	b := all.Bounds()
	img, err := screenshot.Capture(b.Min.X, b.Min.Y, b.Dx(), b.Dy())
	if err != nil {
		panic(err)
	}

	// draw.Draw(all, b, img, b.Min, draw.Src)
	copy(all.Pix, img.Pix)
}

var (
	listen  = flag.String("listen", ":5900", "listen on [ip]:port")
	profile = flag.Bool("profile", false, "write a cpu.prof file when client disconnects")
)

const (
	width  = 1920
	height = 1080
)

func serve() {
	flag.Parse()

	ln, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatal(err)
	}

	s := rfb.NewServer(width, height)
	go func() {
		err = s.Serve(ln)
		log.Fatalf("rfb server ended with: %v", err)
	}()
	log.Infof("Waiting for connection on port 5900\n")

	for c := range s.Conns {
		handleConn(c)
	}
}

func handleConn(c *rfb.Conn) {

	if *profile {
		f, err := os.Create("cpu.prof")
		if err != nil {
			log.Fatal(err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("profiling CPU")
		defer pprof.StopCPUProfile()
		defer log.Printf("stopping profiling CPU")
	}

	mouse := mouseS{clickFlag: false, x: 0, y: 0}

	im := image.NewRGBA(image.Rect(0, 0, width, height))
	li := &rfb.LockableImage{Img: im}

	closec := make(chan bool)
	go func() {
		slide := 0
		tick := time.NewTicker(time.Second / 30)
		defer tick.Stop()
		haveNewFrame := false
		for {
			feed := c.Feed
			if !haveNewFrame {
				feed = nil
			}
			_ = feed
			select {
			case feed <- li:
				haveNewFrame = false
			case <-closec:
				log.Printf("Close conn")
				return
			case <-tick.C:
				slide++
				li.Lock()
				getScreen(im)
				li.Unlock()
				haveNewFrame = true
			}
		}
	}()

	for e := range c.Event {
		//log.Infof("got event: %#v", e)
		if ev, ok := e.(rfb.KeyEvent); ok {
			log.Infof("keyboard  key:%d down:%d", ev.Key, ev.DownFlag)
			if k, ok := keyboardMap[ev.Key]; ok {
				log.Infof("keyboard mapped %d", k)
				if ev.DownFlag == 1 {
					robotgo.KeyTap(k)
				}
			} else {
				if ev.DownFlag == 1 {
					robotgo.UnicodeType(ev.Key)
				}
			}
		}

		if ev, ok := e.(rfb.PointerEvent); ok {
			robotgo.MoveMouse(int(ev.X), int(ev.Y))
			log.Infof("mouse event buttonMask %d\n", ev.ButtonMask)

			if ev.ButtonMask > 0 {
				if !mouse.clickFlag {
					mouse.clickFlag = true
					mouse.x = ev.X
					mouse.y = ev.Y
				} else {
					if mouse.x != ev.X || mouse.y != ev.Y {
						robotgo.MouseToggle("down", "left")
						mouse.x = ev.X
						mouse.y = ev.Y
					}
				}
			} else {
				if mouse.clickFlag {
					mouse.clickFlag = false
					if mouse.x != ev.X || mouse.y != ev.Y {
						robotgo.DragMouse(int(ev.X), int(ev.Y), "left")
						robotgo.MouseToggle("up", "left")
					} else {
						robotgo.MouseClick(clickMap[ev.ButtonMask], false)
					}
				}
			}
		}
	}
	defer close(closec)

	log.Printf("Client disconnected")
}

func drawImage(im *image.RGBA, anim int) {
	pos := 0
	const border = 50
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var r, g, b uint8
			switch {
			case x < border*2.5 && x < int((1.1+math.Sin(float64(y+anim*2)/40))*border):
				r = 255
			case x > width-border*2.5 && x > width-int((1.1+math.Sin(math.Pi+float64(y+anim*2)/40))*border):
				g = 255
			case y < border*2.5 && y < int((1.1+math.Sin(float64(x+anim*2)/40))*border):
				r, g = 255, 255
			case y > height-border*2.5 && y > height-int((1.1+math.Sin(math.Pi+float64(x+anim*2)/40))*border):
				b = 255
			default:
				r, g, b = uint8(x+anim), uint8(y+anim), uint8(x+y+anim*3)
			}
			im.Pix[pos] = r
			im.Pix[pos+1] = g
			im.Pix[pos+2] = b
			pos += 4 // skipping alpha
		}
	}
}
