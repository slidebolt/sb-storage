package main

import (
	runtime "github.com/slidebolt/sb-runtime"

	"github.com/slidebolt/sb-storage/app"
)

func main() {
	runtime.Run(app.New(app.DefaultConfig()))
}
