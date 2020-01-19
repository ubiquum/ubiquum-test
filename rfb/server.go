package rfb

import (
	"context"
	"fmt"
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
	tick := time.NewTicker(time.Second / 2)
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
		case <-tick.C:
			err := screen.Get(im)
			if err != nil {
				log.Errorf("Failed to grab screen: %s", err)
			}
			fmt.Printf("tick\n")
		case msg := <-chClient:
			switch msg.Type() {
			default:
				log.Printf("11 Received message type:%v msg:%v\n", msg.Type(), msg)
			}
		case msg := <-chServer:
			switch msg.Type() {
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
				log.Printf("22 Received message type:%v msg:%v\n", msg.Type(), msg)
			}
		}
	}
}
