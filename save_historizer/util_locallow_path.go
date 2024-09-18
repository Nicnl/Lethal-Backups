package save_historizer

import (
	"os"
	"path/filepath"
)

var localLowPath string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	localLowPath = filepath.Join(homeDir, "AppData", "LocalLow", "ZeekerssRBLX", "Lethal Company")
	//fmt.Println("localLowPath =", localLowPath)
}
