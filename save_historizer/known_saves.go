package save_historizer

import (
	"crypto/md5"
	"fmt"
	"lethal_company_save_manager/save_decoder"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SaveContainer struct {
	Filename string                      `json:"filename"`
	Time     time.Time                   `json:"time"`
	Slot     string                      `json:"slot"`
	Hash     string                      `json:"hash"`
	Infos    save_decoder.LethalSaveInfo `json:"infos"`
}

var (
	saveHistoryDir string
	knownSaves     = make(map[string]SaveContainer)
)

func init() {
	// Défine the save history directory location
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	saveHistoryDir = filepath.Join(pwd, "backups")

	// Create it if missing
	if _, err := os.Stat(saveHistoryDir); os.IsNotExist(err) {
		err := os.Mkdir(saveHistoryDir, 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// List all files in the save history directory
	files, err := os.ReadDir(saveHistoryDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Listing existing backups:")
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Expected format: "{unix_timestamp}__{slot}__{md5_hash}.lethal"
		fileName := file.Name()
		if filepath.Ext(fileName) != ".lethal" {
			continue
		}

		split := strings.Split(fileName, "__")
		if len(split) != 3 {
			continue
		}

		rawTimestamp, slot, hash := split[0], split[1], split[2]
		hash = strings.TrimSuffix(hash, ".lethal")

		// Parse the timestamp
		timestamp, err := time.Parse("2006-01-02_15-04-05", rawTimestamp)
		if err != nil {
			fmt.Println("Error parsing timestamp for", fileName, ":", err)
			continue
		}

		err = CheckSave(timestamp, slot, hash, filepath.Join(saveHistoryDir, fileName), true)
		if err != nil {
			fmt.Println("Error checking save for", fileName, ":", err)
			continue
		} else {
			fmt.Println("  -", fileName)
		}
	}
	fmt.Println("Done listing backups.")
	fmt.Println()
}

func CheckSave(timestamp time.Time, slot string, hash string, filePath string, isBackup bool) error {
	timestamp = timestamp.UTC().Truncate(time.Second)

	if hash != "" {
		_, conflict := knownSaves[hash]
		if conflict {
			fmt.Println("Hash conflict A for hash=", hash)
			return nil
		}
	}

	fileName := filepath.Base(filePath)

	// Load the save file
	rawFile, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading raw file %s => %s", fileName, err)
	}

	// Decrypt the save
	var jsonSave save_decoder.JsonSave

	if isBackup {
		decompressed, err := decompressGzip(rawFile)
		if err != nil {
			return fmt.Errorf("error decompressing save file %s => %s", fileName, err)
		}

		jsonSave = save_decoder.JsonSave{Data: decompressed}
	} else {
		jsonSave, err = save_decoder.Decrypt(save_decoder.EncryptedSave{Data: rawFile})
		if err != nil {
			return fmt.Errorf("error decrypting save file %s => %s", fileName, err)
		}
	}

	if hash == "" {
		hash = fmt.Sprintf("%032x", md5.Sum(jsonSave.Data))
	}

	_, conflict := knownSaves[hash]
	if conflict {
		fmt.Println("Hash conflict B for hash=", hash)
		return nil
	}

	// Read the save
	saveInfos, err := save_decoder.Read(jsonSave)
	if err != nil {
		return fmt.Errorf("error reading save file %s => %s", fileName, err)
	}

	// Add the save to memory
	saveContainer := SaveContainer{
		Filename: fileName,
		Time:     timestamp,
		Slot:     slot,
		Hash:     hash,
		Infos:    saveInfos,
	}

	if !isBackup {
		// Backup the save
		saveContainer.Filename = fmt.Sprintf("%s__%s__%s.lethal", timestamp.Format("2006-01-02_15-04-05"), slot, hash)

		compressed, err := compressGzip(jsonSave.Data)
		if err != nil {
			return fmt.Errorf("error compressing save file %s => %s", fileName, err)
		}

		fmt.Println("New save:")
		fmt.Println("  -", saveContainer.Filename)
		err = os.WriteFile(filepath.Join(saveHistoryDir, saveContainer.Filename), compressed, 0644)
		if err != nil {
			return fmt.Errorf("error backing up save file %s => %s", fileName, err)
		}

		knownSavesDates[slot] = timestamp
	} else {
		curDate, ok := knownSavesDates[slot]
		if !ok || timestamp.After(curDate) {
			knownSavesDates[slot] = timestamp
		}
	}

	//knownSaves[hash] = saveContainer
	//for hash := range knownSaves {
	//	fmt.Println("  Known save =>", hash, "(", knownSaves[hash].Time, ")")
	//}

	//fmt.Printf("  Found save => %s (%s)\n", slot, timestamp)
	//fmt.Printf("    Data = %+v\n", saveInfos)
	//fmt.Println()
	return nil
}
