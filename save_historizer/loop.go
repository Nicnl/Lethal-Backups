package save_historizer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	knownSavesDates = make(map[string]time.Time)

	saveLock = &sync.Mutex{}
)

func Loop() error {
	saveLock.Lock()
	defer saveLock.Unlock()

	entry, err := os.ReadDir(localLowPath)
	if err != nil {
		return err
	}

	for _, file := range entry {
		// Vérifications basiques
		if file.IsDir() {
			continue
		}

		if !strings.HasPrefix(file.Name(), "LCSaveFile") {
			continue
		}

		if strings.Contains(file.Name(), ".") {
			continue
		}

		saveSlot := "slot" + strings.TrimPrefix(file.Name(), "LCSaveFile")

		fileInfo, err := file.Info()
		if err != nil {
			fmt.Println("  Error when obtaining save info :", err)
			continue
		}

		// Vérifier la date de la save en cours et la comparer à la dernière date connue
		saveDate := fileInfo.ModTime().UTC().Truncate(time.Second)
		prevDate, ok := knownSavesDates[saveSlot]
		//fmt.Println("Check save", saveSlot, ":", saveDate, ">", prevDate, "?")
		if !ok || saveDate.After(prevDate) {
			//fmt.Println("saveSlot =", saveSlot)

			err = CheckSave(saveDate, saveSlot, "", filepath.Join(localLowPath, file.Name()), true)
			if err != nil {
				fmt.Println("  Error when checking save :", err)
				continue
			}
		}
	}

	//fmt.Println("loop done")

	return nil
}
