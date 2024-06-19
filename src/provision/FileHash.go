package provision

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func GetFileChecksum(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	return hashString, nil
}

func FilesEqual(file1, file2 string) (bool, error) {
	hash1, err := GetFileChecksum(file1)
	if err != nil {
		return false, err
	}

	hash2, err := GetFileChecksum(file2)
	if err != nil {
		return false, err
	}

	return hash1 == hash2, nil
}
