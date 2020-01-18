package keyboard

var comboKeys = []uint32{
	65507, //: "ctrl",
	65513, //: "alt",
	65505, //: "shift",
	65506, //: "rshift",
	65515, //: "cmd",
	65516, //: "rcmd",
	65514, //: "ralt",
}

var keyboardMap = map[uint32]string{
	// more
	65307: "esc",
	65289: "tab",
	32:    "space",
	// 0: "control",
	65507: "ctrl",
	65513: "alt",
	65505: "shift",
	65506: "rshift",
	65515: "cmd",
	65516: "rcmd",
	65514: "ralt",
	// 0: "command",
	65362: "up",
	65364: "down",
	65361: "left",
	65363: "right",
	65293: "enter",
	65288: "delete",
}
