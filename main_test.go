package main

import (
	"slices"
	"testing"

	"github.com/kluzzebass/reflag/translator"
	_ "github.com/kluzzebass/reflag/translator/ls2eza"
)

func TestShellQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"with space", "'with space'"},
		{"with\ttab", "'with\ttab'"},
		{"with\nnewline", "'with\nnewline'"},
		{"with'quote", "'with'\"'\"'quote'"},
		{"with\"double", "'with\"double'"},
		{"with$dollar", "'with$dollar'"},
		{"with`backtick", "'with`backtick'"},
		{"with\\backslash", "'with\\backslash'"},
		{"with!exclaim", "'with!exclaim'"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := shellQuote(tt.input)
			if result != tt.expected {
				t.Errorf("shellQuote(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseInitArgs(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedShell  string
		expectedAdd    []string
		expectedRemove []string
	}{
		{
			name:           "no args defaults to bash",
			args:           []string{},
			expectedShell:  "bash",
			expectedAdd:    nil,
			expectedRemove: nil,
		},
		{
			name:           "shell only",
			args:           []string{"fish"},
			expectedShell:  "fish",
			expectedAdd:    nil,
			expectedRemove: nil,
		},
		{
			name:           "add translator",
			args:           []string{"bash", "+dig2doggo"},
			expectedShell:  "bash",
			expectedAdd:    []string{"dig2doggo"},
			expectedRemove: nil,
		},
		{
			name:           "remove translator",
			args:           []string{"zsh", "-ls2eza"},
			expectedShell:  "zsh",
			expectedAdd:    nil,
			expectedRemove: []string{"ls2eza"},
		},
		{
			name:           "add and remove",
			args:           []string{"fish", "+dig2doggo", "-ls2eza"},
			expectedShell:  "fish",
			expectedAdd:    []string{"dig2doggo"},
			expectedRemove: []string{"ls2eza"},
		},
		{
			name:           "multiple adds",
			args:           []string{"+dig2doggo", "+more2moor"},
			expectedShell:  "bash",
			expectedAdd:    []string{"dig2doggo", "more2moor"},
			expectedRemove: nil,
		},
		{
			name:           "multiple removes",
			args:           []string{"zsh", "-ls2eza", "-grep2rg"},
			expectedShell:  "zsh",
			expectedAdd:    nil,
			expectedRemove: []string{"ls2eza", "grep2rg"},
		},
		{
			name:           "shell at end",
			args:           []string{"+dig2doggo", "-ls2eza", "fish"},
			expectedShell:  "fish",
			expectedAdd:    []string{"dig2doggo"},
			expectedRemove: []string{"ls2eza"},
		},
		{
			name:           "multiple shells takes last",
			args:           []string{"bash", "fish", "zsh"},
			expectedShell:  "zsh",
			expectedAdd:    nil,
			expectedRemove: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell, add, remove := parseInitArgs(tt.args)
			if shell != tt.expectedShell {
				t.Errorf("parseInitArgs(%v) shell = %q, want %q", tt.args, shell, tt.expectedShell)
			}
			if !slices.Equal(add, tt.expectedAdd) {
				t.Errorf("parseInitArgs(%v) add = %v, want %v", tt.args, add, tt.expectedAdd)
			}
			if !slices.Equal(remove, tt.expectedRemove) {
				t.Errorf("parseInitArgs(%v) remove = %v, want %v", tt.args, remove, tt.expectedRemove)
			}
		})
	}
}

func TestTranslatorRegistry(t *testing.T) {
	// ls2eza should be registered via init()
	tr := translator.Get("ls", "eza")
	if tr == nil {
		t.Fatal("ls2eza translator not registered")
	}

	if tr.Name() != "ls2eza" {
		t.Errorf("Name() = %q, want %q", tr.Name(), "ls2eza")
	}

	// Test GetByName
	tr2 := translator.GetByName("ls2eza")
	if tr2 == nil {
		t.Fatal("GetByName(ls2eza) returned nil")
	}
	if tr2.Name() != "ls2eza" {
		t.Errorf("GetByName(ls2eza).Name() = %q, want %q", tr2.Name(), "ls2eza")
	}

	// Test List
	names := translator.List()
	if !slices.Contains(names, "ls2eza") {
		t.Error("ls2eza not found in List()")
	}

	// Test Get for non-existent translator
	tr3 := translator.Get("foo", "bar")
	if tr3 != nil {
		t.Error("Get(foo, bar) should return nil")
	}
}
