package main

import (
	"flag"
	"fmt"
	"os"
)

// App holds the parsed flag values.
type App struct {
	address    string
	configFile string
	version    bool
}

func main() {
	var app App
	flag.StringVar(&app.address, "a", "", "ip:port to listen on (runs as lambda if empty)")
	flag.StringVar(&app.configFile, "c", "./config.yaml", "name of the config file")
	flag.BoolVar(&app.version, "v", false, "print version information")
	flag.Parse()

	if app.version {
		fmt.Println(versionInfo())
		return
	}

	c, err := NewConfig(app.configFile)
	exitOnErr(err)

	s := NewServer(app.address, c)
	s.run()
}

func exitOnErr(errs ...error) {
	errNotNil := false
	for _, err := range errs {
		if err == nil {
			continue
		}
		errNotNil = true
		fmt.Fprintf(os.Stderr, "ERROR: %s", err.Error())
	}
	if errNotNil {
		fmt.Print("\n")
		os.Exit(-1)
	}
}
