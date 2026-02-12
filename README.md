# tonsh

TON wallet manager for the terminal. Seed phrase is stored in the system keychain — never exposed in the terminal.

## Features

- Manage multiple TON wallets (v3r2)
- Seed phrases stored securely via system keychain (macOS, Linux, Windows)
- Interactive wallet selection or `--wallet <address>` flag
- Mainnet and testnet support
- Balance check via TON network

## Install

**Download pre-built binary** from [Releases](https://github.com/moorzeen/tonsh/releases/latest), or build from source:

```bash
go install github.com/moorzeen/tonsh/cmd/tonsh@latest
```

## Usage

```
tonsh create                       Create a new wallet
tonsh info                         Show wallet info (interactive selection if multiple)
tonsh info --wallet <address>      Show info for a specific wallet
tonsh delete                       Remove wallet from keychain
tonsh delete --wallet <address>    Remove a specific wallet
tonsh version                      Show version
```

Add `--testnet` to any command to use testnet:

```bash
tonsh info --testnet
```

Multiple wallets are selected interactively:

```
Select wallet:
  1. EQD...abc
  2. EQD...def
>
```

## Security

- Seed phrase is never printed to the terminal — stored directly in the system keychain
- To view your seed phrase, open your system keychain manager and search for `tonsh`

## Requirements

- Go 1.21+
- Linux: `libsecret-1-dev` (`sudo apt install libsecret-1-dev`)

## License

MIT
