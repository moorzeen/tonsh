package keychain

import (
	"errors"
	"strings"

	"github.com/zalando/go-keyring"
)

const (
	serviceName  = "tonsh"
	indexAccount = "index"
)

func SaveKey(address, seed string) error {
	if err := keyring.Set(serviceName, address, seed); err != nil {
		return err
	}
	return addToIndex(address)
}

func LoadKey(address string) (string, error) {
	seed, err := keyring.Get(serviceName, address)
	if errors.Is(err, keyring.ErrNotFound) {
		return "", errors.New("wallet not found in keychain")
	}
	return seed, err
}

func DeleteKey(address string) error {
	err := keyring.Delete(serviceName, address)
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return err
	}
	return removeFromIndex(address)
}

func KeyExists(address string) bool {
	_, err := keyring.Get(serviceName, address)
	return err == nil
}

func ListWallets() ([]string, error) {
	data, err := keyring.Get(serviceName, indexAccount)
	if errors.Is(err, keyring.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if data == "" {
		return nil, nil
	}
	return strings.Split(data, "\n"), nil
}

func addToIndex(address string) error {
	wallets, err := ListWallets()
	if err != nil {
		return err
	}
	for _, w := range wallets {
		if w == address {
			return nil
		}
	}
	wallets = append(wallets, address)
	return keyring.Set(serviceName, indexAccount, strings.Join(wallets, "\n"))
}

func removeFromIndex(address string) error {
	wallets, err := ListWallets()
	if err != nil {
		return err
	}
	filtered := make([]string, 0, len(wallets))
	for _, w := range wallets {
		if w != address {
			filtered = append(filtered, w)
		}
	}
	if len(filtered) == 0 {
		err := keyring.Delete(serviceName, indexAccount)
		if errors.Is(err, keyring.ErrNotFound) {
			return nil
		}
		return err
	}
	return keyring.Set(serviceName, indexAccount, strings.Join(filtered, "\n"))
}
