package table

import (
	"github.com/yddeng/dutil/log"
)

var logger *log.Logger

func InitLogger() {
	logger = log.NewLogger("./", "cmd", 1024*1024)
	//log.CloseStdOut()
	logger.Infoln("init log")
}
