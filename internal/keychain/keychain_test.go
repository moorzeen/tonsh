package keychain

import (
	"testing"

	"github.com/zalando/go-keyring"
)

func TestSaveAndLoadKey(t *testing.T) {
	keyring.MockInit()

	addr := "EQD_test_addr_1"
	seed := "word1 word2 word3 word4 word5"

	if err := SaveKey(addr, seed); err != nil {
		t.Fatalf("SaveKey: %v", err)
	}

	got, err := LoadKey(addr)
	if err != nil {
		t.Fatalf("LoadKey: %v", err)
	}
	if got != seed {
		t.Errorf("LoadKey = %q, want %q", got, seed)
	}
}

func TestLoadKey_NotFound(t *testing.T) {
	keyring.MockInit()

	_, err := LoadKey("nonexistent_address")
	if err == nil {
		t.Error("expected error for missing key, got nil")
	}
}

func TestKeyExists(t *testing.T) {
	keyring.MockInit()

	addr := "EQD_test_addr_2"
	if KeyExists(addr) {
		t.Error("wallet should not exist before saving")
	}

	_ = SaveKey(addr, "some seed words")

	if !KeyExists(addr) {
		t.Error("wallet should exist after saving")
	}
}

func TestDeleteKey(t *testing.T) {
	keyring.MockInit()

	addr := "EQD_test_addr_3"
	if err := SaveKey(addr, "seed words"); err != nil {
		t.Fatalf("SaveKey: %v", err)
	}
	if err := DeleteKey(addr); err != nil {
		t.Fatalf("DeleteKey: %v", err)
	}
	if KeyExists(addr) {
		t.Error("wallet still exists after deletion")
	}
}

func TestListWallets_Empty(t *testing.T) {
	keyring.MockInit()

	wallets, err := ListWallets()
	if err != nil {
		t.Fatalf("ListWallets on empty store: %v", err)
	}
	if len(wallets) != 0 {
		t.Errorf("expected 0 wallets, got %d", len(wallets))
	}
}

func TestListWallets(t *testing.T) {
	keyring.MockInit()

	addrs := []string{"EQD_addr_a", "EQD_addr_b", "EQD_addr_c"}
	for _, addr := range addrs {
		if err := SaveKey(addr, "seed"); err != nil {
			t.Fatalf("SaveKey(%s): %v", addr, err)
		}
	}

	wallets, err := ListWallets()
	if err != nil {
		t.Fatalf("ListWallets: %v", err)
	}
	if len(wallets) != len(addrs) {
		t.Errorf("expected %d wallets, got %d", len(addrs), len(wallets))
	}
}

func TestSaveKey_NoDuplicatesInIndex(t *testing.T) {
	keyring.MockInit()

	addr := "EQD_dup_addr"
	for i := 0; i < 3; i++ {
		if err := SaveKey(addr, "seed"); err != nil {
			t.Fatalf("SaveKey iteration %d: %v", i, err)
		}
	}

	wallets, _ := ListWallets()
	count := 0
	for _, w := range wallets {
		if w == addr {
			count++
		}
	}
	if count != 1 {
		t.Errorf("address appears %d times in index, want 1", count)
	}
}

func TestDeleteKey_RemovesFromIndex(t *testing.T) {
	keyring.MockInit()

	addrs := []string{"EQD_x", "EQD_y", "EQD_z"}
	for _, addr := range addrs {
		_ = SaveKey(addr, "seed")
	}

	_ = DeleteKey("EQD_y")

	wallets, _ := ListWallets()
	for _, w := range wallets {
		if w == "EQD_y" {
			t.Error("deleted wallet still present in index")
		}
	}
	if len(wallets) != 2 {
		t.Errorf("expected 2 wallets after deletion, got %d", len(wallets))
	}
}

func TestDeleteKey_LastWallet_ClearsIndex(t *testing.T) {
	keyring.MockInit()

	_ = SaveKey("EQD_only", "seed")
	_ = DeleteKey("EQD_only")

	wallets, err := ListWallets()
	if err != nil {
		t.Fatalf("ListWallets after deleting last wallet: %v", err)
	}
	if len(wallets) != 0 {
		t.Errorf("expected empty index after deleting last wallet, got %v", wallets)
	}
}

func TestDeleteKey_Nonexistent(t *testing.T) {
	keyring.MockInit()

	if err := DeleteKey("EQD_never_saved"); err != nil {
		t.Errorf("DeleteKey on nonexistent wallet returned error: %v", err)
	}
}
