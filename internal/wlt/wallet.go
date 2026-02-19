package wlt

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

const (
	walletVer     = wallet.V3R2
	mainConfigUrl = "https://ton.org/global-config.json"
	testConfigUrl = "https://ton.org/testnet-global.config.json"
)

type Wallet struct {
	Seed       []string
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
	Version    wallet.Version
	Address    string
}

func CreateWallet(seed []string, testnet bool) (*Wallet, error) {
	w, err := wallet.FromSeedWithOptions(nil, seed, walletVer)
	if err != nil {
		return nil, fmt.Errorf("failed to create wlt from seed: %v", err)
	}

	return &Wallet{
		Seed:       seed,
		PrivateKey: w.PrivateKey(),
		PublicKey:  w.PrivateKey().Public().(ed25519.PublicKey),
		Version:    walletVer,
		Address:    w.WalletAddress().Testnet(testnet).String(),
	}, nil
}

// ImportFromPrivateKey creates a Wallet from a hex- or base64-encoded private key.
// Accepts 32-byte (raw seed) or 64-byte (full Ed25519) keys in either encoding.
func ImportFromPrivateKey(keyStr string, testnet bool) (*Wallet, error) {
	keyBytes, err := decodeKey(keyStr)
	if err != nil {
		return nil, err
	}

	var privKey ed25519.PrivateKey
	switch len(keyBytes) {
	case ed25519.SeedSize: // 32 bytes
		privKey = ed25519.NewKeyFromSeed(keyBytes)
	case ed25519.PrivateKeySize: // 64 bytes
		privKey = ed25519.PrivateKey(keyBytes)
	default:
		return nil, fmt.Errorf("invalid key length %d: expected 32 or 64 bytes", len(keyBytes))
	}

	w, err := wallet.FromPrivateKey(nil, privKey, walletVer)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet from private key: %v", err)
	}

	return &Wallet{
		PrivateKey: privKey,
		PublicKey:  privKey.Public().(ed25519.PublicKey),
		Version:    walletVer,
		Address:    w.WalletAddress().Testnet(testnet).String(),
	}, nil
}

// decodeKey tries to decode keyStr as hex, then as standard/URL/raw base64.
func decodeKey(keyStr string) ([]byte, error) {
	if b, err := hex.DecodeString(keyStr); err == nil {
		return b, nil
	}
	for _, enc := range []*base64.Encoding{
		base64.StdEncoding,
		base64.URLEncoding,
		base64.RawStdEncoding,
		base64.RawURLEncoding,
	} {
		if b, err := enc.DecodeString(keyStr); err == nil {
			return b, nil
		}
	}
	return nil, fmt.Errorf("invalid private key: must be hex or base64 encoded (32 or 64 bytes)")
}

func (w *Wallet) GetBalance(testnet bool) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	configUrl := mainConfigUrl
	if testnet {
		configUrl = testConfigUrl
	}

	connection := liteclient.NewConnectionPool()
	err := connection.AddConnectionsFromConfigUrl(ctx, configUrl)
	if err != nil {
		return "", fmt.Errorf("failed to connect to TON network: %v", err)
	}
	defer connection.Stop()

	api := ton.NewAPIClient(connection, ton.ProofCheckPolicyFast).WithRetry()

	tonWallet, err := wallet.FromPrivateKey(api, w.PrivateKey, w.Version)
	if err != nil {
		return "", fmt.Errorf("failed to create wallet: %v", err)
	}

	block, err := api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get current block: %v", err)
	}

	balance, err := tonWallet.GetBalance(ctx, block)
	if err != nil {
		return "", fmt.Errorf("failed to get balance: %v", err)
	}

	return balance.String(), nil
}
