package app

import (
	"log"
	"os"
)

var (
	// InfoLog представялет вид логирования INFO.
	InfoLog *log.Logger

	// WarnLog представялет вид логирования WARN.
	WarnLog *log.Logger

	// ErrLog  представялет вид логирования ERR.
	ErrLog *log.Logger
)

// InitLogger устанавливает настройки для трех видов
// пользовательского логироввания, а именно для INFO, WARN и ERR.
func InitLogger() {

	InfoLog = log.New(os.Stdout, "INFO: ", log.Ltime|log.Ldate|log.Lshortfile)
	WarnLog = log.New(os.Stdout, "WARN: ", log.Ltime|log.Ldate|log.Lshortfile)
	ErrLog = log.New(os.Stdout, "ERR: ", log.Ltime|log.Ldate|log.Lshortfile)
}
