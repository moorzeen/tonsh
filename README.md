# TONsh

TON wallet manager for the terminal. Seed phrases and private keys are stored in the system keychain — never exposed in the terminal or on disk.

## Features

- Manage multiple TON wallets (V3R2)
- Create new wallets or import existing ones from seed phrase or private key
- Seed phrases and private keys stored securely via system keychain (macOS Keychain, Linux Secret Service, Windows Credential Manager)
- Sensitive input (seed phrase, private key) is hidden while typing — no echo
- Interactive mode or direct commands
- Mainnet and testnet support
- Balance check via TON liteserver network

## Install

**Download pre-built binary** from [Releases](https://github.com/moorzeen/tonsh/releases/latest), or install from source:

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
tonsh info                         Show wallet info (auto-selected if only one)
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

When multiple wallets exist, selection is interactive:

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

Input is hidden while typing. The format is detected automatically:

- Multiple words → BIP39 mnemonic seed phrase
- Single token → private key in hex or base64 encoding (32-byte seed or 64-byte full Ed25519 key)

## Security

- Seed phrases and private keys are never written to disk — stored only in the system keychain
- Sensitive input during `import` is never echoed to the terminal
- Each wallet is stored as a separate keychain entry under the service name `tonsh`
- To inspect or back up stored secrets, open your system keychain manager and search for `tonsh`

## Development

```bash
go test ./...
```

On Linux, install `libsecret-1-dev` first:

```bash
sudo apt install libsecret-1-dev
go test ./...
```

Build from source:

```bash
go build -o bin/tonsh ./cmd/tonsh
```

## Requirements

- Go 1.21+
- Linux: `libsecret-1-dev` (`sudo apt install libsecret-1-dev`)

## License

MIT
