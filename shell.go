package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/theMillenniumFalcon/falconDB/index"
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
	input = strings.TrimSuffix(input, "\n")
	args := strings.Split(input, " ")

	switch args[0] {
	case "index":
		indexWrapper()

	case "exit":
		cleanUp(dir)
		os.Exit(1)

	case "lookup":
		return lookupWrapper(args)

	case "delete":
		return deleteWrapper(args)

	case "regenerate":
		index.I.Regenerate()

	default:
		log.Warn("'%s' is not a valid command.", args[0])
		log.Info("valid commands: index, lookup <key> <depth>, delete <key>, regenerate, exit")
	}
	return err
}

func indexWrapper() {
	files := index.I.List()
	log.Success("found %d files in index:", len(files))

	for _, f := range files {
		log.Info(f)
	}
}

func lookupWrapper(args []string) error {
	return nil
}

func deleteWrapper(args []string) error {
	// assert theres a key
	if len(args) < 2 {
		err := fmt.Errorf("no key provided")
		return err
	}

	key := args[1]

	// lookup key, return err if not found
	file, ok := index.I.Lookup(key)
	if !ok {
		err := fmt.Errorf("key doesn't exist")
		return err
	}

	// attempt delete file
	err := index.I.Delete(file)
	if err != nil {
		return err
	}

	log.Success("deleted key %s", key)
	return nil
}
