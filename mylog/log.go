package mylog

import (
	"fmt"
	"log"
	"os"
)

var MyLogger = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)

func Der(v error) {
	if v != nil {
		MyLogger.Fatalf("error:%v", v)
	}
}
func Logf(f string, v ...interface{}) {

	MyLogger.Output(2, fmt.Sprintf(f, v...))

}

type logHelper struct {
	prefix string
}

func (l *logHelper) Write(p []byte) (n int, err error) {

	MyLogger.Printf("%s%s\n", l.prefix, p)
	return len(p), nil
	// }
	// return len(p), nil
}

func newLogHelper(prefix string) *logHelper {
	return &logHelper{prefix}
}
