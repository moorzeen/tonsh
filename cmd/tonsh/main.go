package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/moorzeen/tonsh/internal/handler"
)

var version = "dev"

func main() {
	testnet := hasFlag("--testnet")
	walletFlag := getFlagValue("--wallet")
	command := getCommand()

	switch command {
	case "":
		handler.Interactive(version, testnet)
	case "create":
		handler.Create(testnet)
	case "info":
		handler.Info(walletFlag, testnet)
	case "import":
		handler.Import(testnet)
	case "delete":
		handler.Delete(walletFlag)
	case "version":
		fmt.Printf("TONsh %s\n", version)
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
	fmt.Println("  import          Import a wallet from seed phrase or private key")
	fmt.Println("  info            Show wallet address, balance, version and network")
	fmt.Println("  delete          Remove wallet from keychain")
	fmt.Println("  version         Show version")
	fmt.Println("  help            Show this help message")
	fmt.Println("\nFlags:")
	fmt.Println("  --wallet <addr> Specify wallet address (default: interactive selection)")
	fmt.Println("  --testnet       Use testnet (default: mainnet)")
}

// getCommand returns the first non-flag argument (i.e. the subcommand), or ""
// if none is found (which means interactive mode should be used).
func getCommand() string {
	skipNext := false
	for _, arg := range os.Args[1:] {
		if skipNext {
			skipNext = false
			continue
		}
		if arg == "--wallet" {
			skipNext = true
			continue
		}
		if strings.HasPrefix(arg, "--") {
			continue
		}
		return arg
	}
	return ""
}

func hasFlag(flag string) bool {
	return slices.Contains(os.Args, flag)
}

func getFlagValue(flag string) string {
	for i, arg := range os.Args {
		if arg == flag && i+1 < len(os.Args) {
			return os.Args[i+1]
		}
	}
	return ""
}
