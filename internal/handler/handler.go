package handler

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

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

func confirmOverwrite() bool {
	if !keychain.KeyExists() {
		return true
	}

	fmt.Println("\nWallet already exists in keychain")
	fmt.Println("WARNING: The next action cannot be undone. Make sure you have saved your seed phrase.")

	return confirm("Overwrite the existing wallet? (yes/no): ")
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
	fmt.Println()
}

func Create(testnet bool) {
	if !confirmOverwrite() {
		return
	}

	seed := wallet.NewSeed()
	w, err := wlt.CreateWallet(seed, testnet)
	if err != nil {
		fmt.Printf("Failed to create wallet: %v\n", err)
		return
	}

	seedStr := strings.Join(w.Seed, " ")
	if err := keychain.SaveKey(seedStr); err != nil {
		fmt.Printf("Failed to save key to keychain: %v\n", err)
		return
	}

	fmt.Println("\nWallet successfully created and saved in keychain")
	printInfo(w, strconv.Itoa(0), testnet)
}

func Info(testnet bool) {
	seedStr, err := keychain.LoadKey()
	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println("\nUse 'tonsh create' to create a new wallet.")
		return
	}

	seed := strings.Fields(seedStr)
	if len(seed) != 24 {
		fmt.Printf("invalid seed: expected 24 words, got %d", len(seedStr))
	}

	w, err := wlt.CreateWallet(seed, testnet)
	if err != nil {
		fmt.Printf("Failed to load wallet: %v\n", err)
		return
	}

	fmt.Println("\nGet balance...")

	balance, err := w.GetBalance(testnet)
	if err != nil {
		fmt.Printf("Failed to get balance: %v\n", err)
		return
	}

	fmt.Println("\nWallet info")
	printInfo(w, balance, testnet)
}

func Delete() {
	if !keychain.KeyExists() {
		fmt.Println("Wallet doesn't exist in keychain")
		return
	}

	fmt.Println("\nWARNING: This action cannot be undone! Be sure you have saved your seed phrase.")
	if !confirm("Are you sure you want to delete the wallet from the keychain? (yes/no): ") {
		return
	}

	if err := keychain.DeleteKey(); err != nil {
		fmt.Printf("Failed to delete wallet from keychain: %v\n", err)
		return
	}

	fmt.Println("\nWallet successfully deleted from keychain")
	fmt.Println()
}
