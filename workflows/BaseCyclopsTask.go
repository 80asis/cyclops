package Workflows

import (
	"fmt"
)

type BaseCylopsTask interface {
	Run()
}

func Run() {
	fmt.Println("Test")
}
