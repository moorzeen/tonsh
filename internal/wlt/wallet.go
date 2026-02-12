package wlt

import (
	"context"
	"crypto/ed25519"
	"fmt"

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

const walletVer = wallet.V3R2

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

func (w *Wallet) GetBalance(testnet bool) (string, error) {
	ctx := context.Background()

	var connection *liteclient.ConnectionPool
	var err error

	if testnet {
		connection = liteclient.NewConnectionPool()
		err = connection.AddConnectionsFromConfigUrl(ctx, "https://ton.org/testnet-global.config.json")
	} else {
		connection = liteclient.NewConnectionPool()
		err = connection.AddConnectionsFromConfigUrl(ctx, "https://ton.org/global-config.json")
	}
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
