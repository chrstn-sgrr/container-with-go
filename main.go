package main

import (
	"os"
)

func main() {
	switch os.Args[1] {
	case "run":
		parent()
	case "child":
		child()
	default:
		panic("what should I do??")
	}
}

func parent() {

}

func child() {

}
