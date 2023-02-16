package main

import (
	"bufio"
	"os"

	"github.com/theMillenniumFalcon/falconDB/log"
)

// DefaultDepth is the default depth to resolve reference to
const DefaultDepth = 0

func shell(dir string) error {
	log.IsShellMode = true
	log.Info("starting falcondb shell...")
	setup(dir)

	reader := bufio.NewReader(os.Stdin)
	for {
		// input indicator
		log.Prompt("falcondb> ")

		// read keyboad input
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Warn("err reading input: %s", err.Error())
		}

		// Handle the execution of the input.
		if err = execInput(input, dir); err != nil {
			log.Warn("err executing input: %s", err.Error())
		}
	}
}

func execInput(input string, dir string) (err error) {

}
