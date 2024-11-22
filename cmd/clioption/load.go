package clioption

import (
	"log/slog"

	"github.com/spf13/pflag"
)

type CLIOption struct {
	Port                  uint16
	PersistentDataDirPath string
}

func Load() (CLIOption, error) {
	// Options for key features
	port := pflag.Uint16P("port", "p", 8080, "Port to listen")
	persistentDataDirPath := pflag.StringP("save.dir", "d", "/var/lib/urlshortener/", "Path to directory for persistent data storage")

	// Options for developer
	debugLogEnable := pflag.Bool("debug", false, "Enable debug logs")

	pflag.Parse()

	if *debugLogEnable {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	return CLIOption{
		Port:                  *port,
		PersistentDataDirPath: *persistentDataDirPath,
	}, nil
}
