package table

/*
 mvcrawler logger 日志
*/
import (
	"github.com/yddeng/dutil/log"
)

var logger *log.Logger

func InitLogger() {
	logger = log.NewLogger("./", "event", 1024*1024)
	log.CloseStdOut()
}
