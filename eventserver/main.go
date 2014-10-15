package main

import (
	"log"
	"net/http"
	"runtime"

	"code.google.com/p/go.net/websocket"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
)

type command int

const (
	COMMAND_PLAY command = iota
	COMMAND_PAUSE
	COMMAND_VOLUME_UP
	COMMAND_VOLUME_DOWN
    COMMAND_NEXT
    COMMAND_PREVIOUS
)

var description_map = map[command]string{
    COMMAND_PLAY : "Play",
    COMMAND_PAUSE : "Pause",
    COMMAND_VOLUME_UP: "Volume Up",
    COMMAND_VOLUME_DOWN: "Volume Down",
    COMMAND_NEXT: "Next Item",
    COMMAND_PREVIOUS: "Previous Item",
}

// This seems pretty weird, aw well
var socketChannels []chan command

func YtfdServer(ws *websocket.Conn) {
	myChan := make(chan command)
	socketChannels = append(socketChannels, myChan)
	for {
		cmd := <-myChan
		cmdData := map[string]command{
			"command": cmd,
		}
		websocket.JSON.Send(ws, &cmdData)
	}
	myChan = nil
}

func send_command(command_type command) {
    log.Println("Sending "+description_map[command_type])
	for i := range socketChannels {
		if socketChannels[i] != nil {
			socketChannels[i] <- command_type
		}
	}
}

func main() {
	runtime.GOMAXPROCS(2)

	socketChannels = make([]chan command, 0)
	go func() {
		var pause_state int
		pause_state = 1

		X, _ := xgbutil.NewConn()
		keybind.Initialize(X)

		keybind.KeyPressFun(func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			if pause_state == 1 {
                send_command(COMMAND_PAUSE)
				pause_state = 0
			} else {
				send_command(COMMAND_PLAY)
				pause_state = 1
			}
		}).Connect(X, X.RootWin(), "XF86AudioPlay", true)

		keybind.KeyPressFun(func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			send_command(COMMAND_VOLUME_UP)
		}).Connect(X, X.RootWin(), "XF86AudioRaiseVolume", true)

		keybind.KeyPressFun(func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			send_command(COMMAND_VOLUME_DOWN)
		}).Connect(X, X.RootWin(), "XF86AudioLowerVolume", true)

        keybind.KeyPressFun(func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
            send_command(COMMAND_VOLUME_DOWN)
        }).Connect(X, X.RootWin(), "XF86AudioLowerVolume", true)

        keybind.KeyPressFun(func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
            send_command(COMMAND_PREVIOUS)
        }).Connect(X, X.RootWin(), "XF86Back", true)

        keybind.KeyPressFun(func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
            send_command(COMMAND_NEXT)
        }).Connect(X, X.RootWin(), "XF86Mail", true)

		xevent.Main(X)

		/*
			X, _ := xgb.NewConn()
			screen := xproto.Setup(X).DefaultScreen(X)
			wid := screen.Root
			xproto.GrabKey(X, false, wid, xproto.ModMaskControl|xproto.ModMask1|xproto.ModMask2, 118, xproto.GrabModeAsync, xproto.GrabModeAsync)
		*/
	}()

	http.Handle("/ws", websocket.Handler(YtfdServer))
	err := http.ListenAndServe(":42050", nil)
	if err != nil {
		log.Println(err)
	}
}
