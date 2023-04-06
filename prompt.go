package aime

import (
	_ "embed"
	"strings"
)

//go:embed system/prompt/bazel.txt
var BazelPrompt string

//go:embed system/prompt/golang.txt
var GolangPrompt string

func LookupSystemPrompt(key string) string {
	switch strings.ToLower(key) {
	case "bazel", "bzl":
		return BazelPrompt
	case "go", "golang":
		return GolangPrompt
	default:
		return ""
	}
}