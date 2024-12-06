package logger

import (
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	log.SetOutput(os.Stdout)
}
