package btc

import (
	"github.com/joho/godotenv"

	"testing"
	"os"
)

func TestExportBlock(t *testing.T) {
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		err := godotenv.Load()
		if err != nil {
			t.Fatal(err)
		}
	}

	//if err := LoadBlockToFile(10); err != nil {
	if err := LoadBlockToFile(251718); err != nil {
		t.Error(err)
	}
	t.Error()
}
