package keychain

import (
	"errors"

	"github.com/zalando/go-keyring"
)

const (
	serviceName = "tonsh"
	accountName = "wallet-seed"
)

func SaveKey(seed string) error {
	return keyring.Set(serviceName, accountName, seed)
}

func LoadKey() (string, error) {
	seed, err := keyring.Get(serviceName, accountName)
	if errors.Is(err, keyring.ErrNotFound) {
		return "", errors.New("wallet key not found in keychain")
	}
	return seed, err
}

func DeleteKey() error {
	err := keyring.Delete(serviceName, accountName)
	if errors.Is(err, keyring.ErrNotFound) {
		return nil
	}
	return err
}

func KeyExists() bool {
	_, err := keyring.Get(serviceName, accountName)
	return err == nil
}
