package logger

import (
	"log"
	"os"
)

var (
	Info    *log.Logger
	Error   *log.Logger
	Trace   *log.Logger
	Debug   *log.Logger
	Warning *log.Logger
	Fatal   *log.Logger
)

func init() {

	Info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Error = log.New(os.Stdout,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Trace = log.New(os.Stdout,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Debug = log.New(os.Stdout,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Warning = log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Fatal = log.New(os.Stdout,
		"FATAL: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
}

/*В любой пакет нужно испортировать ""github.com/go-park-mail-ru/2019_1_SleeplessNights/log""
log.<log>.Println("commit")*/
