package screen2tmux

import (
	"reflect"
	"testing"
)

func TestTranslateFlags(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		// Default: new session
		{
			name:     "bare screen",
			input:    []string{},
			expected: []string{"new-session"},
		},

		// New session with name (-S)
		{
			name:     "new session with name",
			input:    []string{"-S", "mysession"},
			expected: []string{"new-session", "-s", "mysession"},
		},
		{
			name:     "new session with bundled S",
			input:    []string{"-Smysession"},
			expected: []string{"new-session", "-s", "mysession"},
		},

		// New session with command
		{
			name:     "new session with command",
			input:    []string{"vim"},
			expected: []string{"new-session", "--", "vim"},
		},
		{
			name:     "new session with command and args",
			input:    []string{"vim", "file.txt"},
			expected: []string{"new-session", "--", "vim", "file.txt"},
		},
		{
			name:     "new session with name and command",
			input:    []string{"-S", "dev", "vim", "file.txt"},
			expected: []string{"new-session", "-s", "dev", "--", "vim", "file.txt"},
		},
		{
			name:     "command after -- separator",
			input:    []string{"--", "vim"},
			expected: []string{"new-session", "--", "vim"},
		},
		{
			name:     "name and command after separator",
			input:    []string{"-S", "dev", "--", "vim", "file.txt"},
			expected: []string{"new-session", "-s", "dev", "--", "vim", "file.txt"},
		},

		// Detached new session (-dm)
		{
			name:     "detached new session",
			input:    []string{"-dm", "-S", "bg"},
			expected: []string{"new-session", "-d", "-s", "bg"},
		},
		{
			name:     "detached new session separate flags",
			input:    []string{"-d", "-m", "-S", "bg"},
			expected: []string{"new-session", "-d", "-s", "bg"},
		},
		{
			name:     "detached new session bundled dmS",
			input:    []string{"-dmS", "bg"},
			expected: []string{"new-session", "-d", "-s", "bg"},
		},
		{
			name:     "detached new session with command",
			input:    []string{"-dmS", "bg", "python", "server.py"},
			expected: []string{"new-session", "-d", "-s", "bg", "--", "python", "server.py"},
		},
		{
			name:     "d without m does not start detached",
			input:    []string{"-d", "-S", "mysession"},
			expected: []string{"new-session", "-s", "mysession"},
		},

		// Reattach (-r)
		{
			name:     "reattach no name",
			input:    []string{"-r"},
			expected: []string{"attach"},
		},
		{
			name:     "reattach with name",
			input:    []string{"-r", "mysession"},
			expected: []string{"attach", "-t", "mysession"},
		},

		// Multi-attach (-x)
		{
			name:     "multi attach no name",
			input:    []string{"-x"},
			expected: []string{"attach"},
		},
		{
			name:     "multi attach with name",
			input:    []string{"-x", "mysession"},
			expected: []string{"attach", "-t", "mysession"},
		},

		// Detach and reattach (-d -r)
		{
			name:     "detach and reattach",
			input:    []string{"-d", "-r", "mysession"},
			expected: []string{"attach", "-d", "-t", "mysession"},
		},
		{
			name:     "detach and reattach bundled",
			input:    []string{"-dr", "mysession"},
			expected: []string{"attach", "-d", "-t", "mysession"},
		},

		// Power detach and reattach (-D -R)
		{
			name:     "power detach and reattach",
			input:    []string{"-D", "-R", "mysession"},
			expected: []string{"attach", "-d", "-t", "mysession"},
		},
		{
			name:     "power detach and reattach bundled",
			input:    []string{"-DR", "mysession"},
			expected: []string{"attach", "-d", "-t", "mysession"},
		},

		// Reattach or create (-R alone)
		{
			name:     "reattach or create no name",
			input:    []string{"-R"},
			expected: []string{"new-session", "-A"},
		},
		{
			name:     "reattach or create with name",
			input:    []string{"-R", "mysession"},
			expected: []string{"new-session", "-A", "-s", "mysession"},
		},

		// List sessions
		{
			name:     "list sessions -ls",
			input:    []string{"-ls"},
			expected: []string{"list-sessions"},
		},
		{
			name:     "list sessions -list",
			input:    []string{"-list"},
			expected: []string{"list-sessions"},
		},
		{
			name:     "wipe sessions",
			input:    []string{"-wipe"},
			expected: []string{"list-sessions"},
		},

		// Config file (-c)
		{
			name:     "config file",
			input:    []string{"-c", "myconfig"},
			expected: []string{"-f", "myconfig", "new-session"},
		},
		{
			name:     "config file with session",
			input:    []string{"-c", "myconfig", "-S", "mysession"},
			expected: []string{"-f", "myconfig", "new-session", "-s", "mysession"},
		},
		{
			name:     "config file with attach",
			input:    []string{"-c", "myconfig", "-r", "mysession"},
			expected: []string{"-f", "myconfig", "attach", "-t", "mysession"},
		},

		// Ignored flags
		{
			name:     "ignored flags a A",
			input:    []string{"-aA"},
			expected: []string{"new-session"},
		},
		{
			name:     "ignored flag U",
			input:    []string{"-U"},
			expected: []string{"new-session"},
		},
		{
			name:     "ignored flag q",
			input:    []string{"-q"},
			expected: []string{"new-session"},
		},
		{
			name:     "ignored flag i",
			input:    []string{"-i"},
			expected: []string{"new-session"},
		},
		{
			name:     "ignored flag O",
			input:    []string{"-O"},
			expected: []string{"new-session"},
		},
		{
			name:     "ignored flag with arg h",
			input:    []string{"-h", "1000"},
			expected: []string{"new-session"},
		},
		{
			name:     "ignored flag with arg t",
			input:    []string{"-t", "title"},
			expected: []string{"new-session"},
		},
		{
			name:     "ignored flag with arg T",
			input:    []string{"-T", "xterm"},
			expected: []string{"new-session"},
		},
		{
			name:     "ignored flag with arg e",
			input:    []string{"-e", "^Aa"},
			expected: []string{"new-session"},
		},
		{
			name:     "ignored flag with arg s",
			input:    []string{"-s", "/bin/zsh"},
			expected: []string{"new-session"},
		},
		{
			name:     "ignored flag with arg p",
			input:    []string{"-p", "0"},
			expected: []string{"new-session"},
		},
		{
			name:     "ignored flow control fn",
			input:    []string{"-fn"},
			expected: []string{"new-session"},
		},
		{
			name:     "ignored flow control fa",
			input:    []string{"-fa"},
			expected: []string{"new-session"},
		},
		{
			name:     "ignored flow control f alone",
			input:    []string{"-f"},
			expected: []string{"new-session"},
		},
		{
			name:     "ignored login mode l",
			input:    []string{"-l"},
			expected: []string{"new-session"},
		},
		{
			name:     "ignored login mode ln",
			input:    []string{"-ln"},
			expected: []string{"new-session"},
		},

		// Complex bundled flags
		{
			name:     "ignored flags mixed with session name",
			input:    []string{"-aUS", "mysession"},
			expected: []string{"new-session", "-s", "mysession"},
		},
		{
			name:     "multiple ignored flags then attach",
			input:    []string{"-q", "-r", "mysession"},
			expected: []string{"attach", "-t", "mysession"},
		},

		// Edge cases
		{
			name:     "r does not consume flag-like arg as name",
			input:    []string{"-r", "-S", "name"},
			expected: []string{"attach", "-t", "name"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := translateFlags(tt.input)
			if len(result) == 0 && len(tt.expected) == 0 {
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("translateFlags(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTranslator(t *testing.T) {
	tr := &Translator{}

	if tr.Name() != "screen2tmux" {
		t.Errorf("Name() = %s, want screen2tmux", tr.Name())
	}
	if tr.SourceTool() != "screen" {
		t.Errorf("SourceTool() = %s, want screen", tr.SourceTool())
	}
	if tr.TargetTool() != "tmux" {
		t.Errorf("TargetTool() = %s, want tmux", tr.TargetTool())
	}
	if !tr.IncludeInInit() {
		t.Error("IncludeInInit() = false, want true")
	}

	// Test Translate method delegates to translateFlags
	input := []string{"-S", "test", "vim"}
	expected := []string{"new-session", "-s", "test", "--", "vim"}
	result := tr.Translate(input, "")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Translate(%v) = %v, want %v", input, result, expected)
	}
}
