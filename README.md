# TONsh

TON wallet manager for the terminal. Seed phrase is stored in the system keychain — never exposed in the terminal.

## Features

- Manage multiple TON wallets (v3r2)
- Create new wallets or import existing ones from seed phrase or private key
- Seed phrases and private keys stored securely via system keychain (macOS, Linux, Windows)
- Sensitive input (seed phrase, private key) is never displayed while typing
- Interactive mode or direct commands
- Mainnet and testnet support
- Balance check via TON network

## Install

**Download pre-built binary** from [Releases](https://github.com/moorzeen/tonsh/releases/latest), or build from source:

```bash
go install github.com/moorzeen/tonsh/cmd/tonsh@latest
```

## Usage

Run without arguments to start interactive mode:

```bash
tonsh
tonsh --testnet
```

Or use commands directly:

```
tonsh create                       Create a new wallet
tonsh import                       Import a wallet from seed phrase or private key
tonsh info                         Show wallet info (interactive selection if multiple)
tonsh info --wallet <address>      Show info for a specific wallet
tonsh delete                       Remove wallet from keychain
tonsh delete --wallet <address>    Remove a specific wallet
tonsh version                      Show version
tonsh help                         Show help
```

Add `--testnet` to any command to use testnet:

```bash
tonsh create --testnet
tonsh info --testnet
```

When multiple wallets exist, wallet is selected interactively:

```
Select wallet:
1. EQD...abc
2. EQD...def
>
```

## Import

`tonsh import` accepts either a **seed phrase** (24 words separated by spaces) or a **private key**:

```
Enter seed phrase or private key:
```

The input is hidden while typing. The app detects the format automatically:

- Multiple words → treated as a BIP39 seed phrase
- Single token → treated as a private key in hex or base64 encoding (32 or 64 bytes)

## Security

- Seed phrases and private keys are never displayed in the terminal — stored directly in the system keychain
- Sensitive input during `import` is hidden while typing (no echo)
- To view stored secrets, open your system keychain manager and search for `tonsh`

## Development

```bash
go test ./...
```

On Linux, install `libsecret-1-dev` first:

```bash
sudo apt install libsecret-1-dev
go test ./...
```

## Requirements

- Go 1.21+
- Linux: `libsecret-1-dev` (`sudo apt install libsecret-1-dev`)

## License

MIT