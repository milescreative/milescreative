package main

import (
	"testing"

	"go-std/internal/config"
)

func BenchmarkConfigLookup(b *testing.B) {
	cfg, _ := config.Config()
	cfg.LoadJSON("../../test.jsonc")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = cfg.GetString("app_name")
	}
}

func BenchmarkConfigLookupCached(b *testing.B) {
	cfg, _ := config.Config()
	cfg.LoadJSON("../../test.jsonc")
	value := cfg.GetString("app_name")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = value
	}
}

func BenchmarkHardcodedVariable(b *testing.B) {
	value := "go-std"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = value
	}
}
