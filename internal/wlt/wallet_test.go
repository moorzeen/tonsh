package wlt

import (
	"testing"

	"github.com/xssnick/tonutils-go/ton/wallet"
)

func TestCreateWallet_Deterministic(t *testing.T) {
	seed := wallet.NewSeed()

	w1, err := CreateWallet(seed, false)
	if err != nil {
		t.Fatalf("first CreateWallet: %v", err)
	}
	w2, err := CreateWallet(seed, false)
	if err != nil {
		t.Fatalf("second CreateWallet: %v", err)
	}

	if w1.Address != w2.Address {
		t.Errorf("same seed produced different addresses: %q vs %q", w1.Address, w2.Address)
	}
}

func TestCreateWallet_MainnetVsTestnet(t *testing.T) {
	seed := wallet.NewSeed()

	mainnet, err := CreateWallet(seed, false)
	if err != nil {
		t.Fatalf("mainnet CreateWallet: %v", err)
	}
	testnet, err := CreateWallet(seed, true)
	if err != nil {
		t.Fatalf("testnet CreateWallet: %v", err)
	}

	if mainnet.Address == testnet.Address {
		t.Error("mainnet and testnet addresses should differ for the same seed")
	}
}

func TestCreateWallet_SeedPreserved(t *testing.T) {
	seed := wallet.NewSeed()

	w, err := CreateWallet(seed, false)
	if err != nil {
		t.Fatalf("CreateWallet: %v", err)
	}

	if len(w.Seed) != len(seed) {
		t.Fatalf("seed length: got %d, want %d", len(w.Seed), len(seed))
	}
	for i, word := range seed {
		if w.Seed[i] != word {
			t.Errorf("seed[%d] = %q, want %q", i, w.Seed[i], word)
		}
	}
}

func TestCreateWallet_VersionIsV3R2(t *testing.T) {
	seed := wallet.NewSeed()

	w, err := CreateWallet(seed, false)
	if err != nil {
		t.Fatalf("CreateWallet: %v", err)
	}

	if w.Version != walletVer {
		t.Errorf("wallet version = %v, want %v", w.Version, walletVer)
	}
}

func TestCreateWallet_KeysNotEmpty(t *testing.T) {
	seed := wallet.NewSeed()

	w, err := CreateWallet(seed, false)
	if err != nil {
		t.Fatalf("CreateWallet: %v", err)
	}

	if len(w.PrivateKey) == 0 {
		t.Error("private key is empty")
	}
	if len(w.PublicKey) == 0 {
		t.Error("public key is empty")
	}
	if w.Address == "" {
		t.Error("address is empty")
	}
}

func TestCreateWallet_InvalidSeed(t *testing.T) {
	_, err := CreateWallet([]string{"invalid", "seed", "words"}, false)
	if err == nil {
		t.Error("expected error for invalid seed, got nil")
	}
}

func TestCreateWallet_DifferentSeeds_DifferentAddresses(t *testing.T) {
	seed1 := wallet.NewSeed()
	seed2 := wallet.NewSeed()

	w1, err := CreateWallet(seed1, false)
	if err != nil {
		t.Fatalf("CreateWallet seed1: %v", err)
	}
	w2, err := CreateWallet(seed2, false)
	if err != nil {
		t.Fatalf("CreateWallet seed2: %v", err)
	}

	if w1.Address == w2.Address {
		t.Error("different seeds produced the same address")
	}
}
