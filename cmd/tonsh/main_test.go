package main

import (
	"os"
	"testing"
)

func TestGetCommand(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"no args", []string{"tonsh"}, ""},
		{"testnet only", []string{"tonsh", "--testnet"}, ""},
		{"create", []string{"tonsh", "create"}, "create"},
		{"create testnet after", []string{"tonsh", "create", "--testnet"}, "create"},
		{"testnet before create", []string{"tonsh", "--testnet", "create"}, "create"},
		{"info with wallet after", []string{"tonsh", "info", "--wallet", "EQD123"}, "info"},
		{"wallet before info", []string{"tonsh", "--wallet", "EQD123", "info"}, "info"},
		{"wallet only", []string{"tonsh", "--wallet", "EQD123"}, ""},
		{"delete", []string{"tonsh", "delete"}, "delete"},
		{"version", []string{"tonsh", "version"}, "version"},
		{"help", []string{"tonsh", "help"}, "help"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Args
			os.Args = tt.args
			defer func() { os.Args = old }()
			if got := getCommand(); got != tt.want {
				t.Errorf("getCommand() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHasFlag(t *testing.T) {
	tests := []struct {
		args []string
		flag string
		want bool
	}{
		{[]string{"tonsh", "--testnet"}, "--testnet", true},
		{[]string{"tonsh", "create", "--testnet"}, "--testnet", true},
		{[]string{"tonsh", "create"}, "--testnet", false},
		{[]string{"tonsh"}, "--testnet", false},
		{[]string{"tonsh", "--wallet", "EQD"}, "--wallet", true},
		{[]string{"tonsh", "info"}, "--wallet", false},
	}
	for _, tt := range tests {
		old := os.Args
		os.Args = tt.args
		got := hasFlag(tt.flag)
		os.Args = old
		if got != tt.want {
			t.Errorf("hasFlag(%q) with args %v = %v, want %v", tt.flag, tt.args, got, tt.want)
		}
	}
}

func TestGetFlagValue(t *testing.T) {
	tests := []struct {
		args []string
		flag string
		want string
	}{
		{[]string{"tonsh", "--wallet", "EQD123"}, "--wallet", "EQD123"},
		{[]string{"tonsh", "info", "--wallet", "EQD456"}, "--wallet", "EQD456"},
		{[]string{"tonsh", "info"}, "--wallet", ""},
		{[]string{"tonsh", "--wallet"}, "--wallet", ""},
	}
	for _, tt := range tests {
		old := os.Args
		os.Args = tt.args
		got := getFlagValue(tt.flag)
		os.Args = old
		if got != tt.want {
			t.Errorf("getFlagValue(%q) with args %v = %q, want %q", tt.flag, tt.args, got, tt.want)
		}
	}
}
