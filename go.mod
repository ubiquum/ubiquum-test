module github.com/ubiquum/ubiquum

go 1.13

replace github.com/ubiquum/screenshot => ../screenshot

replace github.com/bradfitz/rfbgo => ../rfbgo

require (
	github.com/BurntSushi/xgb v0.0.0-20160522181843-27f122750802 // indirect
	github.com/bradfitz/rfbgo v0.0.0-00010101000000-000000000000
	github.com/gen2brain/shm v0.0.0-20191025110947-b09d223a76f1 // indirect
	github.com/kbinani/screenshot v0.0.0-20191211154542-3a185f1ce18f
)
