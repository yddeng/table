package table

import (
	"encoding/json"
	"fmt"
	"github.com/sniperHW/kendynet"
	listener "github.com/sniperHW/kendynet/socket/listener/websocket"
	"github.com/yddeng/table/conf"
	"net/http"
)

func Start(path string) {
	conf.LoadConfig(path)
	_conf := conf.GetConfig()

	InitLogger()
	go Loop()

	server, err := listener.New("tcp4", _conf.WSAddr, "/table")
	if server != nil {
		fmt.Printf("webasocket start on %s\n", _conf.WSAddr)
		go func() {
			err = server.Serve(func(session kendynet.StreamSession) {
				fmt.Println("new Session", session.RemoteAddr())
				session.SetCloseCallBack(func(sess kendynet.StreamSession, reason string) {
					PostTask(func() {
						OnClose(sess, reason)
					})
				})
				session.Start(func(event *kendynet.Event) {
					if event.EventType == kendynet.EventTypeError {
						event.Session.Close(event.Data.(error).Error(), 0)
					} else {
						msg := map[string]interface{}{}
						err := json.Unmarshal(event.Data.(kendynet.Message).Bytes(), &msg)
						if err == nil {
							PostTask(func() {
								Dispatcher(msg, session)
							})
						}
					}
				})
			})
			if nil != err {
				panic(fmt.Sprintf("TcpServer start failed%s\n", err))
			}
		}()
	} else {
		panic(fmt.Sprintf("NewTcpServer failed %s\n", err))
	}

	fmt.Printf("http start on %s, LoadDir on %s\n", _conf.HttpAddr, _conf.LoadDir)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(_conf.LoadDir))))

	// http handler
	http.HandleFunc("/createTable", HandleCreateTable)
	http.HandleFunc("/deleteTable", HandleDeleteTable)
	http.HandleFunc("/getAllTable", HandleGetAllTable)
	http.HandleFunc("/downloadTable", HandleDownloadTable)
	err = http.ListenAndServe(_conf.HttpAddr, nil)
	if err != nil {
		fmt.Println(err)
	}
}
