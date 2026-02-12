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
	testnet := hasFlag("--testnet")
	walletFlag := getFlagValue("--wallet")

	switch command {
	case "create":
		handler.Create(testnet)
	case "info":
		handler.Info(walletFlag, testnet)
	case "delete":
		handler.Delete(walletFlag)
	case "version":
		fmt.Printf("tonsh %s\n", version)
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("tonsh", version)
	fmt.Println("\nUsage:")
	fmt.Println("  tonsh <command> [flags]")
	fmt.Println("\nCommands:")
	fmt.Println("  create          Create a new TON wallet")
	fmt.Println("  info            Show wallet address, balance, version and network")
	fmt.Println("  delete          Remove wallet from keychain")
	fmt.Println("  version         Show version")
	fmt.Println("  help            Show this help message")
	fmt.Println("\nFlags:")
	fmt.Println("  --wallet <addr> Specify wallet address (default: interactive selection)")
	fmt.Println("  --testnet       Use testnet (default: mainnet)")
}

func hasFlag(flag string) bool {
	for _, arg := range os.Args {
		if arg == flag {
			return true
		}
	}
	return false
}

func getFlagValue(flag string) string {
	for i, arg := range os.Args {
		if arg == flag && i+1 < len(os.Args) {
			return os.Args[i+1]
		}
	}
	return ""
}
