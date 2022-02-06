//go:build tools
// +build tools

// this is a tools.go file for pinning tool versions as recommended by https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module and
// https://github.com/golang/go/issues/25922#issuecomment-413898264

package tools

import (
	_ "github.com/celo-org/celo-blockchain/cmd/abigen"
)
