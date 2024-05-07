package cyclopsWorklows

import (
	"fmt"
)

type BaseCylopsTask interface {
	Run()
}

func Run() {
	fmt.Println("Test")
}
