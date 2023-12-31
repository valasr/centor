package cli

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
)

func PrintLogo() {
	myFigure := figure.NewFigure("CENTOR", "", true)
	myFigure.Print()
	fmt.Println()
}
