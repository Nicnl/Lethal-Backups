package save_historizer

import (
	"encoding/json"
	"fmt"
	"lethal_company_save_manager/save_decoder"
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

			err = CheckSave(saveDate, saveSlot, "", filepath.Join(localLowPath, file.Name()), false)
			if err != nil {
				fmt.Println("  Error when checking save :", err)
				continue
			}
		}
	}

	//fmt.Println("loop done")

	return nil
}

func ListKnownSaves() ([]byte, error) {
	saveLock.Lock()
	defer saveLock.Unlock()

	jsonResp, err := json.Marshal(knownSaves)
	if err != nil {
		return nil, err
	}

	return jsonResp, nil
}

func NbSaves() int {
	saveLock.Lock()
	defer saveLock.Unlock()

	return len(knownSaves)
}

func ObtainSave(hash string) (SaveContainer, bool) {
	saveLock.Lock()
	defer saveLock.Unlock()

	save, ok := knownSaves[hash]
	return save, ok
}

func (sc *SaveContainer) Restore() error {
	saveLock.Lock()
	defer saveLock.Unlock()

	// Read the backup file (gzip of json)
	backupFilePath := filepath.Join(saveHistoryDir, sc.Filename)

	gzipSave, err := os.ReadFile(backupFilePath)
	if err != nil {
		return fmt.Errorf("error reading backup file %s => %s", sc.Filename, err)
	}

	// Decompress the backup file
	jsonSave, err := decompressGzip(gzipSave)
	if err != nil {
		return fmt.Errorf("error decompressing backup file %s => %s", sc.Filename, err)
	}

	// Encrypt the save
	encryptedSave, err := save_decoder.Encrypt(save_decoder.JsonSave{Data: jsonSave})
	if err != nil {
		return fmt.Errorf("error encrypting save file %s => %s", sc.Filename, err)
	}

	// Write the save file
	destinationPath := filepath.Join(localLowPath, "LCSaveFile"+strings.TrimPrefix(sc.Slot, "slot"))
	//fmt.Println("destinationPath =", destinationPath)
	//fmt.Println("encryptedSave =", encryptedSave)

	err = os.WriteFile(destinationPath, encryptedSave.Data, 0644)
	if err != nil {
		return fmt.Errorf("error writing save file %s => %s", sc.Filename, err)
	}

	return nil
}
