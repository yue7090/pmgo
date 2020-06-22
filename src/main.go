package main

import (
	"flag"
	"fmt"
	"pm2go/lib"
	"os"
)

func main() {
	configPath := flag.String("c", "", "config file path")
	flag.Parse()

	if *configPath != "" {
		if _, err := os.Stat(*configPath); err != null {
			fmt.Printf("Can't find config file `%s`\n", *configPath)
			os.Exit(1)
		}else {
			os.Setenv("RUNNER_CONFIG_PATH", *configPath)
		}
	}

	lib.Start()
}
