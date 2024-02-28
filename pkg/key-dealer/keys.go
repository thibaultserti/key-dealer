package keydealer

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/iam/v1"
)

const keysDir = "./keys"

func MakeKeyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		saEmail := r.PathValue("sa_email")

		// Check key is present
		keyPath := fmt.Sprintf(keysDir+"/%s.json", saEmail)
		if _, err := os.Stat(keyPath); err == nil {
			keyBytes, err := os.ReadFile(keyPath)
			if err != nil {
				logrus.Errorf("Error when reading key: %v", err)
				http.Error(w, "Error when reading key", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(keyBytes)
			return
		}

		// Generate key
		gcpSaKey, err := createGCPKey(saEmail)
		if err != nil {
			logrus.Errorf("Error when generating GCP SA key: %v", err)
		}
		newKey, _ := base64.StdEncoding.DecodeString(gcpSaKey.PrivateKeyData)

		err = os.WriteFile(keyPath, []byte(string(newKey)), 0644)
		if err != nil {
			logrus.Errorf("Error when writing key to file: %v", err)
			http.Error(w, "Error when writing key to file", http.StatusInternalServerError)
			return
		}

		// Renvoyer la clé générée
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(string(newKey)))
	}
}

// createKey creates a service account key.
func createGCPKey(serviceAccountEmail string) (*iam.ServiceAccountKey, error) {
	ctx := context.Background()
	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("iam.NewService: %w", err)
	}

	resource := "projects/-/serviceAccounts/" + serviceAccountEmail
	request := &iam.CreateServiceAccountKeyRequest{}
	key, err := service.Projects.ServiceAccounts.Keys.Create(resource, request).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.ServiceAccounts.Keys.Create: %w", err)
	}

	return key, nil
}

func DeleteKeys() error {
	files := deleteKeysOnDisk()
	for _, file := range files {
		serviceAccountEmail := strings.TrimSuffix(file.Name(), ".json")
		err := deleteAllGCPKeys(serviceAccountEmail)
		if err != nil {
			return fmt.Errorf("error when deleting key for sa %s", serviceAccountEmail)
		}
	}
	return nil
}

func deleteAllGCPKeys(serviceAccountEmail string) error {
	ctx := context.Background()
	service, err := iam.NewService(ctx)
	if err != nil {
		return fmt.Errorf("iam.NewService: %w", err)
	}

	resource := "projects/-/serviceAccounts/" + serviceAccountEmail

	keysList, err := service.Projects.ServiceAccounts.Keys.List(resource).Do()
	if err != nil {
		return fmt.Errorf("Projects.ServiceAccounts.Keys.List: %w", err)
	}
	for _, key := range keysList.Keys {
		if key.KeyType == "USER_MANAGED" {
			_, err = service.Projects.ServiceAccounts.Keys.Delete(key.Name).Do()
			if err != nil {
				return fmt.Errorf("Projects.ServiceAccounts.Keys.Delete: %w", err)
			}
		}
	}
	return nil
}

func deleteKeysOnDisk() []fs.DirEntry {
	// Delete all keys in keysDir
	files, err := os.ReadDir(keysDir)
	if err != nil {
		logrus.Errorf("Failed to read keys directory: %v", err)
	}

	for _, file := range files {
		err := os.Remove(filepath.Join(keysDir, file.Name()))
		if err != nil {
			logrus.Errorf("Failed to delete key %s: %v", file.Name(), err)
		} else {
			logrus.Infof("Deleted key: %s", file.Name())
		}
	}
	return files

}
