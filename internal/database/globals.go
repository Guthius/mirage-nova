package database

import (
	"log"
	"os"
)

var Log = log.New(os.Stderr, "[Database] ", log.LstdFlags)
