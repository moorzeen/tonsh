package handler

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"

	"github.com/moorzeen/tonsh/internal/keychain"
	"github.com/moorzeen/tonsh/internal/wlt"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

func confirm(prompt string) bool {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(response)) != "yes" {
		fmt.Println("Canceled")
		return false
	}
	return true
}

func selectWallet(wallets []string) (string, error) {
	fmt.Println()
	fmt.Println("Select wallet:")
	for i, addr := range wallets {
		fmt.Printf("%d. %s\n", i+1, addr)
	}
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	n, err := strconv.Atoi(input)
	if err != nil || n < 1 || n > len(wallets) {
		return "", fmt.Errorf("invalid selection")
	}
	return wallets[n-1], nil
}

func resolveWallet(walletFlag string) (string, error) {
	if walletFlag != "" {
		if !keychain.KeyExists(walletFlag) {
			return "", fmt.Errorf("wallet not found: %s", walletFlag)
		}
		return walletFlag, nil
	}

	wallets, err := keychain.ListWallets()
	if err != nil {
		return "", err
	}
	if len(wallets) == 0 {
		return "", fmt.Errorf("\nno wallets found\n\nUse 'tonsh create' to create a new wallet.")
	}
	if len(wallets) == 1 {
		return wallets[0], nil
	}
	return selectWallet(wallets)
}

func printTonscanLink(address string, testnet bool) {
	if testnet {
		fmt.Printf("https://testnet.tonscan.org/address/%s\n", address)
	} else {
		fmt.Printf("https://tonscan.org/address/%s\n", address)
	}
}

func printNetwork(testnet bool) {
	network := "Mainnet"
	if testnet {
		network = "Testnet"
	}
	fmt.Printf("Network: %s\n", network)
}

func printInfo(w *wlt.Wallet, balance string, testnet bool) {
	fmt.Printf("Address: %s\n", w.Address)
	fmt.Printf("Balance: %s TON\n", balance)
	fmt.Printf("Version: %s\n", w.Version)
	printNetwork(testnet)
	printTonscanLink(w.Address, testnet)
}

func Interactive(ver string, testnet bool) {
	network := "mainnet"
	if testnet {
		network = "testnet"
	}
	fmt.Printf("TONsh %s (%s)\n", ver, network)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println()
		fmt.Println("Select action:")
		fmt.Println("1. Info")
		fmt.Println("2. Create")
		fmt.Println("3. Import")
		fmt.Println("4. Delete")
		fmt.Println("5. Exit")
		fmt.Print("> ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			Info("", testnet)
		case "2":
			Create(testnet)
		case "3":
			Import(testnet)
		case "4":
			Delete("")
		case "5":
			fmt.Println()
			fmt.Println("Bye!")
			return
		default:
			fmt.Println("Invalid selection")
		}
	}
}

func Create(testnet bool) {
	seed := wallet.NewSeed()
	w, err := wlt.CreateWallet(seed, testnet)
	if err != nil {
		fmt.Printf("Failed to create wallet: %v\n", err)
		return
	}

	seedStr := strings.Join(w.Seed, " ")
	if err := keychain.SaveKey(w.Address, seedStr); err != nil {
		fmt.Printf("Failed to save key to keychain: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("Wallet successfully created and saved in keychain")
	fmt.Println()
	printInfo(w, strconv.Itoa(0), testnet)
	fmt.Println()
	fmt.Println("To view your seed phrase, open your system keychain manager and search for \"tonsh\"")
}

// loadWalletFromKeychain loads a wallet from the keychain by address.
// The stored value is either a space-separated seed phrase or "privkey:<hex>".
func loadWalletFromKeychain(address string, testnet bool) (*wlt.Wallet, error) {
	stored, err := keychain.LoadKey(address)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(stored, "privkey:") {
		return wlt.ImportFromPrivateKey(strings.TrimPrefix(stored, "privkey:"), testnet)
	}
	return wlt.CreateWallet(strings.Fields(stored), testnet)
}

func Info(walletFlag string, testnet bool) {
	address, err := resolveWallet(walletFlag)
	if err != nil {
		fmt.Println(err)
		return
	}

	w, err := loadWalletFromKeychain(address, testnet)
	if err != nil {
		fmt.Printf("Failed to load wallet: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("Get balance...")
	balance, err := w.GetBalance(testnet)
	if err != nil {
		fmt.Printf("Failed to get balance: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("Wallet info")
	fmt.Println()
	printInfo(w, balance, testnet)
}

func Import(testnet bool) {
	fmt.Println()
	fmt.Print("Enter seed phrase or private key: ")

	// Read without echo so the secret is not displayed in the terminal.
	raw, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		// Fallback to plain input when stdin is not a terminal (e.g. tests/pipes).
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		raw = []byte(strings.TrimSpace(line))
	}

	input := strings.TrimSpace(string(raw))
	if input == "" {
		fmt.Println("Nothing entered")
		return
	}

	var w *wlt.Wallet
	var storedValue string

	words := strings.Fields(input)
	if len(words) > 1 {
		// Multiple words → seed phrase
		w, err = wlt.CreateWallet(words, testnet)
		if err != nil {
			fmt.Printf("Failed to import from seed phrase: %v\n", err)
			return
		}
		storedValue = strings.Join(words, " ")
	} else {
		// Single token → private key (hex or base64)
		w, err = wlt.ImportFromPrivateKey(input, testnet)
		if err != nil {
			fmt.Printf("Failed to import from private key: %v\n", err)
			return
		}
		storedValue = "privkey:" + hex.EncodeToString(w.PrivateKey)
	}

	if keychain.KeyExists(w.Address) {
		fmt.Println()
		fmt.Printf("Wallet %s is already in keychain\n", w.Address)
		return
	}

	if err := keychain.SaveKey(w.Address, storedValue); err != nil {
		fmt.Printf("Failed to save to keychain: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("Wallet successfully imported and saved in keychain")
	fmt.Println()
	fmt.Printf("Address: %s\n", w.Address)
	fmt.Printf("Version: %s\n", w.Version)
	printNetwork(testnet)
	printTonscanLink(w.Address, testnet)
	fmt.Println()
	fmt.Println("Use 'tonsh info' to check the balance")
}

func Delete(walletFlag string) {
	address, err := resolveWallet(walletFlag)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println()
	fmt.Println("WARNING: This action cannot be undone! Be sure you have saved your seed phrase.")
	if !confirm(fmt.Sprintf("Delete wallet %s? (yes/no): ", address)) {
		return
	}

	if err := keychain.DeleteKey(address); err != nil {
		fmt.Printf("Failed to delete wallet: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("Wallet successfully deleted from keychain")
}
