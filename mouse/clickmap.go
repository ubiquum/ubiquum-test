package mouse

type mouseS struct {
	clickFlag bool
	clickType uint8
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
