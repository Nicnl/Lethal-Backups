package save_historizer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var knownSavesDates = make(map[string]time.Time)

func Loop() error {
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

		fmt.Println()

		saveSlot := strings.TrimPrefix(file.Name(), "LCSaveFile")
		fmt.Println("Check save : slot", saveSlot)

		fileInfo, err := file.Info()
		if err != nil {
			fmt.Println("  Error when obtaining save info :", err)
			continue
		}

		// Vérifier la date de la save en cours et la comparer à la dernière date connue
		saveDate := fileInfo.ModTime()
		prevDate, ok := knownSavesDates[saveSlot]
		if !ok || saveDate.After(prevDate) {
			knownSavesDates[saveSlot] = saveDate

			err = CheckSave(saveDate, saveSlot, "", filepath.Join(localLowPath, file.Name()), true)
			if err != nil {
				fmt.Println("  Error when checking save :", err)
				continue
			}
		}
		fmt.Println(" => ok")
	}

	return nil
}
