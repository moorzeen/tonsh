package main

import (
	"fmt"
	"os"

	"github.com/moorzeen/tonsh/internal/handler"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	command := os.Args[1]
	testnet := isTestnet()

	switch command {
	case "create":
		handler.Create(testnet)
	case "info":
		handler.Info(testnet)
	case "delete":
		handler.Delete()
	case "version":
		fmt.Printf("TONsh %s\n", version)
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(0)
	}
}

func printUsage() {
	fmt.Println("TONsh", version)
	fmt.Println("\nUsage:")
	fmt.Println("  tonsh <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  create          Create new TON wallet and save seed phrase in keychain")
	fmt.Println("  info            Show wallet info: address, balance, qr-code, version, network)")
	fmt.Println("  delete          Delete wallet from keychain")
	fmt.Println("  version         Show version")
	fmt.Println("  help            Show this help message")
	fmt.Println("\nFlags:")
	fmt.Println("  --testnet       Use testnet network (default: mainnet)")
}

func isTestnet() bool {
	for _, arg := range os.Args {
		if arg == "--testnet" {
			return true
		}
	}
	return false
}
