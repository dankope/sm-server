package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"
	"time"

	datafile "github.com/d-ank/otdata"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 10 * time.Second
)

var (
	port      = flag.String("port", ":8080", "overlay port")
	hook      datafile.Hook
	homeTempl = template.Must(template.New("").Parse(homeData))
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func writer(ws *websocket.Conn) {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case data := <-hook.Reader:
			// tell the client all the data, you could also do this a number of different ways
			// maybe you would prefer to run this code without the middleman with WASM
			if data != nil {
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
					return
				}
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	go writer(ws)
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}
func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTempl.Execute(w, nil)
}

func main() {
	flag.Parse()
	var err error
	hook, err = datafile.Add("streamermode.data")
	if err != nil {
		log.Fatal(err)
	}
	defer hook.Close()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	log.Fatal(http.ListenAndServe(*port, nil))
}

var homeData = `<!DOCTYPE html>
<html lang="en">

<head>
    <title>stream overlay</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/pixi.js/5.1.3/pixi.min.js"></script>
    <script>
        window.onload = function () {
            var conn = new WebSocket("ws://" + location.host + "/ws");
            const app = new PIXI.Application({
                antialias: true,
                transparent: true,
                resizeTo: window,
                autoDensity: true,
                resolution: devicePixelRatio
            });
            document.body.appendChild(app.view);
            let text = new PIXI.Text('Onetap Anti Leak', {
                fontFamily: 'Tahoma',
                fontSize: 26,
                fill: 0xFFFFFF,
                align: 'left'
            });
            text.alpha = 0;
            const graphics = new PIXI.Graphics();
            conn.onclose = function (evt) {
                graphics.clear();
            }
            conn.onmessage = function (evt) {
                data = JSON.parse(evt.data);
                graphics.clear();
                if (data["MENU"]["OPEN"]) {
                    graphics.beginFill(0x232328);
                    graphics.drawRect(data["MENU"]["INFO"][0], data["MENU"]["INFO"][1], data["MENU"]["INFO"][2], data["MENU"]["INFO"][3]);
                    graphics.endFill();
                }
                if (data["CONSOLE"]) {
                    graphics.beginFill(0x232328);
                    graphics.drawRect(0, 0, app.screen.width, app.screen.height);
                    graphics.endFill();
                }
            }
            app.stage.addChild(graphics);
        }
    </script>
</head>

<body>
    <div id="main">

    </div>
</body>

</html>`
