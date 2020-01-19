package mouse

import (
	"github.com/go-vgo/robotgo"
	log "github.com/sirupsen/logrus"
	vnc "github.com/unistack-org/go-rfb"
)

var mouseState = mouseS{clickFlag: false, x: 0, y: 0}

func HandleEvent(ev *vnc.PointerEvent) error {

	robotgo.MoveMouse(int(ev.X), int(ev.Y))

	if ev.Mask > 0 {
		if !mouseState.clickFlag {
			mouseState.clickFlag = true
			mouseState.clickType = ev.Mask
			mouseState.x = ev.X
			mouseState.y = ev.Y
		} else {
			if mouseState.x != ev.X || mouseState.y != ev.Y {
				robotgo.MouseToggle("down", "left")
				mouseState.x = ev.X
				mouseState.y = ev.Y
			}
		}
	} else {
		if mouseState.clickFlag {
			mouseState.clickFlag = false
			if mouseState.x != ev.X || mouseState.y != ev.Y {
				robotgo.DragMouse(int(ev.X), int(ev.Y), "left")
				robotgo.MouseToggle("up", "left")
			} else {
				log.Infof("mouse event buttonMask %v\n", ev)
				robotgo.MouseClick(clickMap[mouseState.clickType], false)
			}
		}
	}

	return nil
}
