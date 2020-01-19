module github.com/ubiquum/ubiquum

go 1.13

replace github.com/ubiquum/screenshot => ../screenshot

require (
	github.com/go-vgo/robotgo v0.0.0-20200105130635-9c21c09ef9f3
	github.com/kbinani/screenshot v0.0.0-20191211154542-3a185f1ce18f
	github.com/sirupsen/logrus v1.4.2
	github.com/ubiquum/screenshot v0.0.0-00010101000000-000000000000 // indirect
	github.com/unistack-org/go-rfb v0.0.0-20181217213206-b1cfa9ee459a
	google.golang.org/appengine v1.6.5
)
