package wlt

import (
	"encoding/base64"
	"encoding/hex"
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

func TestImportFromPrivateKey_HexFull(t *testing.T) {
	seed := wallet.NewSeed()
	w, err := CreateWallet(seed, false)
	if err != nil {
		t.Fatalf("CreateWallet: %v", err)
	}

	hexKey := hex.EncodeToString(w.PrivateKey)
	imported, err := ImportFromPrivateKey(hexKey, false)
	if err != nil {
		t.Fatalf("ImportFromPrivateKey (hex full): %v", err)
	}
	if imported.Address != w.Address {
		t.Errorf("address mismatch: got %q, want %q", imported.Address, w.Address)
	}
}

func TestImportFromPrivateKey_HexSeed(t *testing.T) {
	seed := wallet.NewSeed()
	w, err := CreateWallet(seed, false)
	if err != nil {
		t.Fatalf("CreateWallet: %v", err)
	}

	// 32-byte private key seed
	hexSeed := hex.EncodeToString(w.PrivateKey.Seed())
	imported, err := ImportFromPrivateKey(hexSeed, false)
	if err != nil {
		t.Fatalf("ImportFromPrivateKey (hex seed): %v", err)
	}
	if imported.Address != w.Address {
		t.Errorf("address mismatch: got %q, want %q", imported.Address, w.Address)
	}
}

func TestImportFromPrivateKey_Base64(t *testing.T) {
	seed := wallet.NewSeed()
	w, err := CreateWallet(seed, false)
	if err != nil {
		t.Fatalf("CreateWallet: %v", err)
	}

	b64Key := base64.StdEncoding.EncodeToString(w.PrivateKey)
	imported, err := ImportFromPrivateKey(b64Key, false)
	if err != nil {
		t.Fatalf("ImportFromPrivateKey (base64): %v", err)
	}
	if imported.Address != w.Address {
		t.Errorf("address mismatch: got %q, want %q", imported.Address, w.Address)
	}
}

func TestImportFromPrivateKey_InvalidKey(t *testing.T) {
	_, err := ImportFromPrivateKey("notavalidkey!!!", false)
	if err == nil {
		t.Error("expected error for invalid key, got nil")
	}
}

func TestImportFromPrivateKey_WrongLength(t *testing.T) {
	// Valid hex but wrong length (16 bytes)
	_, err := ImportFromPrivateKey(hex.EncodeToString(make([]byte, 16)), false)
	if err == nil {
		t.Error("expected error for wrong key length, got nil")
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
