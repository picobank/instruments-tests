package etl

import (
	"fmt"
	"testing"
)

func TestExtract(t *testing.T) {
	extractor := Extract()

	sigend := Load(extractor)

	<-sigend

	fmt.Println("Extract and Load finished !!!!!!!")
}
