package btc

import (
	"github.com/joho/godotenv"

	"os"
	"testing"
)

func TestExportBlock(t *testing.T) {
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		err := godotenv.Load()
		if err != nil {
			t.Fatal(err)
		}
	}

	if err := LoadBlockToFile(".", 265458); err != nil {
		t.Error(err)
	}
}
