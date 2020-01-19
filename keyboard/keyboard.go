package keyboard

import (
	"github.com/go-vgo/robotgo"
	log "github.com/sirupsen/logrus"
	vnc "github.com/unistack-org/go-rfb"
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
	log.Debugf("down %s", skey)
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

func HandleEvent(ev *vnc.KeyEvent) error {

	flag := "up"
	if ev.Down == 1 {
		flag = "down"
	}
	log.Tracef("key:%s %s", ev.Key, flag)

	// handle combo keys [ctrl, alt, shift, cmd]
	key := uint32(ev.Key)
	if isComboKey(key) {
		if ev.Down == 1 {
			addKey(key)
		} else {
			removeKey(key)
		}
		return nil
	}

	// key up
	if ev.Down == 0 {
		return nil
	}

	keyVal := string(rune(ev.Key))

	robotgo.KeyTap(keyVal, keys)
	log.Debugf("key:%s keys:%v", keyVal, keys)

	keys = []string{}

	return nil
}
