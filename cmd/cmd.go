package cmd

import (
	_ "github.com/arlettebrook/serv00-ct8/configs"
	"github.com/common-nighthawk/go-figure"
)

func Start() {
	figure.NewColorFigure("serv00&ct8", "doom",
		"cyan", true).Print()

	Execute()

}
