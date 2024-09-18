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
	Filename string
	Time     time.Time
	Slot     string
	Hash     string
	Infos    save_decoder.LethalSaveInfo
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

	fmt.Println("Listing known saves:")
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

		if strings.HasPrefix(slot, "slot") {
			slot = strings.TrimPrefix(slot, "slot")
		}

		// Parse the timestamp
		timestamp, err := time.Parse("2006-01-02_15-04-05", rawTimestamp)
		if err != nil {
			fmt.Println("Error parsing timestamp for", fileName, ":", err)
			continue
		}

		err = CheckSave(timestamp, slot, hash, filepath.Join(saveHistoryDir, fileName), false)
	}
}

func CheckSave(timestamp time.Time, slot string, hash string, filePath string, shouldBackup bool) error {
	if hash != "" {
		_, conflict := knownSaves[slot]
		if conflict {
			return nil
		}
	}

	fileName := filepath.Base(filePath)

	// Load the save file
	rawSave, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading save file %s => %s", fileName, err)
	}

	if hash == "" {
		hash = fmt.Sprintf("%032x", md5.Sum(rawSave))
	}

	_, conflict := knownSaves[slot]
	if conflict {
		return nil
	}

	// Decrypt the save
	decrypted, err := save_decoder.Decrypt(rawSave)
	if err != nil {
		return fmt.Errorf("error decrypting save file %s => %s", fileName, err)
	}

	// Read the save
	infos, err := save_decoder.Read(decrypted)
	if err != nil {
		return fmt.Errorf("error reading save file %s => %s", fileName, err)
	}

	// Add the save to memory
	knownSaves[slot] = SaveContainer{
		Filename: fileName,
		Time:     timestamp,
		Slot:     slot,
		Hash:     hash,
		Infos:    infos,
	}

	if shouldBackup {
		// Backup the save
		backupName := fmt.Sprintf("%s__%s__%s.lethal", timestamp.Format("2006-01-02_15-04-05"), slot, hash)
		err = os.WriteFile(filepath.Join(saveHistoryDir, backupName), rawSave, 0644)
		if err != nil {
			return fmt.Errorf("error backing up save file %s => %s", fileName, err)
		}
	}

	fmt.Printf("  Found save => %s (%s)\n", slot, timestamp)
	fmt.Printf("    Data = %+v\n", infos)
	fmt.Println()
	return nil
}
