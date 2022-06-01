package main

import (
	"flag"
	"fmt"
	"os"
)

const configFileEnv = "CONFIG_FILE"

// App holds the parsed flag values.
type App struct {
	address    string
	configFile string
	mode       string
	version    bool
}

func main() {
	var app App
	flag.StringVar(&app.address, "a", "", "ip:port to listen on when run locally")
	flag.StringVar(&app.configFile, "c", "./config.yaml", "name of the config file")
	flag.StringVar(&app.mode, "m", "local", "mode, can me either 'local', 'azurefunc' or 'awslambda'")
	flag.BoolVar(&app.version, "v", false, "print version information")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s \n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "(alternatively the -c parameter can be specified via environment variable '%s')\n\n", configFileEnv)
		flag.PrintDefaults()
	}
	flag.Parse()

	envConfig := os.Getenv(configFileEnv)
	if envConfig != "" {
		app.configFile = envConfig
	}

	if app.version {
		fmt.Println(versionInfo())
		return
	}

	c, err := NewConfig(app.configFile)
	exitOnErr(err)

	s := NewServer(app.address, c)
	s.run(app.mode)
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
