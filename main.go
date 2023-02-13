package main

import (
	_ "github.com/theMillenniumFalcon/falconDB/api"
	"github.com/theMillenniumFalcon/falconDB/index"
)

func main() {
	index.I.Regenerate()
}
