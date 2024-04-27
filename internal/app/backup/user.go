package backup

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
)

type BackupData struct {
	oapi.NewUser
	JWT string `json: "jwt"`
}

// TODO - save JWT, not user.
// Backup current User to tmp file.
func SaveUser(u oapi.NewUser, jwt string) error {
	// Save user to file.
	file, error := os.OpenFile(getBackupPath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if error != nil {
		return error
	}
	defer file.Close()

	// Create backup user.
	user := BackupData{u, jwt}
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	file.Write(data)
	return nil
}

// Load backup User from tmp file.
func LoadUser() (user *BackupData, err error) {
	file, err := os.Open(getBackupPath())
	if err != nil {
		return nil, err
	}
	u := &BackupData{}
	err = json.NewDecoder(file).Decode(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Load backup User from tmp file.
func CleanData() (err error) {
	err = os.Remove(getBackupPath())
	if err != nil {
		return err
	}
	return
}

func getBackupPath() string {
	tmp := os.TempDir()
	return fmt.Sprint(tmp, "/", "gophkeeper.usr")
}
