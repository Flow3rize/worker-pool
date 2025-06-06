package worker

import (
	"log"
	"os"
)

var (
	InfoLogger  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lmsgprefix)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lmsgprefix)
)
