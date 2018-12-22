package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"syscall/js"

	"github.com/madlambda/qi/qi"
	"github.com/madlambda/qi/yin/dev"
)

// Version of yin
const Version = "v0.1"

func main() {
	fmt.Printf("yin version %s\n", Version)

	err := bootstrap()
	if err != nil {
		fmt.Printf("err: %s\n", err)
		os.Exit(1)
	}
}

var (
	doc           js.Value
	ctx           js.Value
	width, height float64

	// in-memory filesystem
	fs qi.Filesystem

	// root file (aka /)
	root qi.File

	// registered devices
	devices map[string]qi.File
)

func init() {
	devices = make(map[string]qi.File)
}

func bootstrap() error {
	// Init Canvas
	doc = js.Global().Get("document")
	canvasEl := doc.Call("getElementById", "screen")
	width = doc.Get("body").Get("clientWidth").Float()
	height = doc.Get("body").Get("clientHeight").Float()
	canvasEl.Set("width", width)
	canvasEl.Set("height", height)
	canvasEl.Set("style", "cursor: none")
	ctx = canvasEl.Call("getContext", "2d")

	fs := qi.NewRoot()

	mouseFile := dev.MouseInit()
	moveEvent := js.NewCallback(func(args []js.Value) {
		e := args[0]
		mouseFile.UpdateCoords(int(e.Get("clientX").Float()), int(e.Get("clientY").Float()))
	})

	defer moveEvent.Release()

	doc.Call("addEventListener", "mousemove", moveEvent)
	devices["#m"] = mouseFile

	err := fs.Append("/dev", mouseFile)
	abortonerr(err)

	fmt.Printf("Filesystem:\n")
	fs.Walk("/", func(f qi.File) {
		fmt.Printf("\t/%s\n", f.Name())
	})

	fmt.Printf("Devices loaded: \n")
	for devname := range devices {
		fmt.Printf("\t%s\n", devname)
	}

	fmt.Printf("Mouse:\n")

	for {
		mouseFile.Seek(io.SeekStart, 0)
		mouseData, err := ioutil.ReadAll(mouseFile)
		abortonerr(err)

		fmt.Printf("\t>> %s\n", string(mouseData))
	}

	return nil
}

func abortonerr(err error) {
	if err != nil {
		panic(err)
	}
}
