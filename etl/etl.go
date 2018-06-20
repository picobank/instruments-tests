package etl

import (
	"fmt"
)

func run() {
	extractor := Extract()

	for {
		data, open := <-extractor
		if !open {
			break
		}
		fmt.Println("Extracted data: ", data)
	}
}
