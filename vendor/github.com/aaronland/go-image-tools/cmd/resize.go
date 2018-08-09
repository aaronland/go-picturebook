package main

import (
	"flag"
	"github.com/aaronland/go-image-tools/resize"
	"log"
)

func main() {

	max := flag.Int("max", 640, "...")

	flag.Parse()

	args := flag.Args()

	for _, path := range args {

		resized_path, err := resize.ResizeMax(path, *max)

		if err != nil {
			log.Fatal(err)
		}

		log.Println(resized_path)
	}

}
