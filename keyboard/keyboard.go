package keyboard

import (
	"github.com/bradfitz/rfbgo/rfb"
	"github.com/go-vgo/robotgo"
	log "github.com/sirupsen/logrus"
)

var keys = []string{}

func getKey(key uint32) string {
	if k, ok := keyboardMap[key]; ok {
		return k
	}

	log.Warnf("key not mapped: %d", key)
	return ""
}

func addKey(k uint32) {
	skey := getKey(k)
	for _, key := range keys {
		if key == skey {
			break
		}
	}
	keys = append(keys, skey)
	log.Infof("down %s", skey)
}

func removeKey(k uint32) {
	for i, key := range keys {
		skey := getKey(k)
		if key == skey {
			keys[i] = keys[len(keys)-1]
			keys[len(keys)-1] = ""
			keys = keys[:len(keys)-1]
			log.Debugf("up %s", skey)
			break
		}
	}
}

func isComboKey(k uint32) bool {
	for _, key := range comboKeys {
		if key == k {
			return true
		}
	}
	return false
}

func HandleKey(ev rfb.KeyEvent) error {

	flag := "up"
	if ev.DownFlag == 1 {
		flag = "down"
	}
	log.Tracef("key:%d %s", ev.Key, flag)

	// handle combo keys [ctrl, alt, shift, cmd]
	if isComboKey(ev.Key) {
		if ev.DownFlag == 1 {
			addKey(ev.Key)
		} else {
			removeKey(ev.Key)
		}
		return nil
	}

	// key up
	if ev.DownFlag == 0 {
		return nil
	}

	keyVal := string(rune(ev.Key))

	robotgo.KeyTap(keyVal, keys)
	log.Infof("key:%s keys:%v", keyVal, keys)

	keys = []string{}

	return nil
}
