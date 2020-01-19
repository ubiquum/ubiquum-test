package rfb

import (
	"context"
	"image"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ubiquum/ubiquum/keyboard"
	"github.com/ubiquum/ubiquum/mouse"
	"github.com/ubiquum/ubiquum/screen"
	vnc "github.com/unistack-org/go-rfb"
)

const (
	width  = 1920
	height = 1080
)

func Serve(host string) error {

	ln, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("Error listen. %v", err)
		return err
	}

	chServer := make(chan vnc.ClientMessage)
	chClient := make(chan vnc.ServerMessage)

	im := image.NewRGBA(image.Rect(0, 0, width, height))
	tick := time.NewTicker(time.Second * 2)
	defer tick.Stop()

	cfg := &vnc.ServerConfig{
		Width:            width,
		Height:           height,
		Handlers:         vnc.DefaultServerHandlers,
		SecurityHandlers: []vnc.SecurityHandler{&vnc.ClientAuthNone{}},
		Encodings:        []vnc.Encoding{&vnc.RawEncoding{}},
		PixelFormat:      vnc.PixelFormat32bit,
		ClientMessageCh:  chServer,
		ServerMessageCh:  chClient,
		Messages:         vnc.DefaultClientMessages,
	}

	go vnc.Serve(context.Background(), ln, cfg)
	log.Infof("Listening on %s", host)

	// Process messages coming in on the ClientMessage channel.
	for {
		select {
		case err := <-cfg.ErrorCh:
			log.Errorf("VNC error: %s", err)
		case <-tick.C:

			err := screen.Get(im)
			if err != nil {
				log.Errorf("Failed to grab screen: %s", err)
			}

			// fmt.Printf("tick\n")
		case msg := <-chClient:
			switch msg.Type() {
			default:
				log.Printf("client: Received message type:%v msg:%v\n", msg.Type(), msg)
			}
		case msg := <-chServer:
			switch msg.Type() {
			case vnc.FramebufferUpdateRequestMsgType:

				ev, ok := msg.(*vnc.FramebufferUpdateRequest)
				if !ok {
					continue
				}

				log.Infof("FramebufferUpdateRequest %dx%d [%d,%d]", ev.Width, ev.Height, ev.X, ev.Y)

				// r := image.Rect(int(ev.X), int(ev.Y), int(ev.X+ev.Width), int(ev.Y+ev.Height))
				// im1 := im.SubImage(r)

				msg := new(vnc.FramebufferUpdate)
				rawEnc := new(vnc.RawEncoding)
				rawEnc.Colors = []vnc.Color{}

				r1 := vnc.NewRectangle()
				r1.EncType = vnc.EncRaw
				r1.Enc = rawEnc

				r1.Width = ev.Width
				r1.Height = ev.Height
				r1.X = ev.X
				r1.Y = ev.Y

				msg.NumRect = 1
				msg.Rects = []*vnc.Rectangle{
					r1,
				}

				chClient <- msg
				log.Info("Sent response")

			case vnc.PointerEventMsgType:
				ev, ok := msg.(*vnc.PointerEvent)
				if !ok {
					continue
				}
				mouse.HandleEvent(ev)
			case vnc.KeyEventMsgType:
				ev, ok := msg.(*vnc.KeyEvent)
				if !ok {
					continue
				}
				keyboard.HandleEvent(ev)
			default:
				log.Printf("server: Received message type:%v msg:%v\n", msg.Type(), msg)
			}
		}
	}
}
