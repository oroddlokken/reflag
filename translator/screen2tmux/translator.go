package screen2tmux

import (
	"strings"

	"github.com/kluzzebass/reflag/translator"
)

func init() {
	translator.Register(&Translator{})
}

// Translator implements the screen to tmux flag translation
type Translator struct{}

func (t *Translator) Name() string        { return "screen2tmux" }
func (t *Translator) SourceTool() string  { return "screen" }
func (t *Translator) TargetTool() string  { return "tmux" }
func (t *Translator) IncludeInInit() bool { return true }

// Translate converts screen arguments to tmux arguments
func (t *Translator) Translate(args []string, mode string) []string {
	return translateFlags(args)
}

func translateFlags(args []string) []string {
	var (
		operation   = "new" // "new", "attach", "list"
		sessionName string
		configFile  string
		detachFlag  bool // -d or -D
		bigR        bool // -R (reattach or create)
		hasM        bool // -m (ignore $STY, enables detached new-session with -d)
		command     []string
	)

	i := 0
	for i < len(args) {
		arg := args[i]

		// Special compound flags (must be standalone args)
		if arg == "-ls" || arg == "-list" || arg == "-wipe" {
			operation = "list"
			i++
			continue
		}

		// End of options: rest is the command
		if arg == "--" {
			command = append(command, args[i+1:]...)
			break
		}

		// Not a flag: remaining args are the command
		if !strings.HasPrefix(arg, "-") {
			command = append(command, args[i:]...)
			break
		}

		// Process bundled short flags
		j := 1
		for j < len(arg) {
			ch := arg[j]
			switch ch {
			case 'r':
				operation = "attach"
				j++
				// Optional session name: if last in bundle, check next arg
				if j >= len(arg) && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					sessionName = args[i+1]
					i++
				}
			case 'R':
				bigR = true
				j++
				if j >= len(arg) && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					sessionName = args[i+1]
					i++
				}
			case 'x':
				operation = "attach"
				j++
				if j >= len(arg) && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					sessionName = args[i+1]
					i++
				}
			case 'd':
				detachFlag = true
				j++
			case 'D':
				detachFlag = true
				j++
			case 'm':
				hasM = true
				j++
			case 'S':
				j++
				// Session name: rest of this arg or next arg
				if j < len(arg) {
					sessionName = arg[j:]
					j = len(arg)
				} else if i+1 < len(args) {
					sessionName = args[i+1]
					i++
				}
			case 'c':
				j++
				// Config file: rest of this arg or next arg
				if j < len(arg) {
					configFile = arg[j:]
					j = len(arg)
				} else if i+1 < len(args) {
					configFile = args[i+1]
					i++
				}
			case 'f':
				j++
				// -fn or -fa (flow control variants)
				if j < len(arg) && (arg[j] == 'n' || arg[j] == 'a') {
					j++
				}
			case 'l':
				j++
				// -ln (login mode off)
				if j < len(arg) && arg[j] == 'n' {
					j++
				}
			// Flags that consume an argument (all ignored for tmux)
			case 'e', 'h', 'p', 's', 't', 'T':
				j++
				if j < len(arg) {
					j = len(arg) // rest of bundle is the argument value
				} else if i+1 < len(args) {
					i++ // consume next arg
				}
			// Ignored standalone flags
			case 'a', 'A', 'i', 'n', 'O', 'q', 'U':
				j++
			default:
				j++
			}
		}

		i++
	}

	// Resolve -R behavior: with detach → attach, without → new-session -A
	if bigR && operation != "attach" && operation != "list" {
		if detachFlag {
			operation = "attach"
		} else {
			operation = "reattach_or_create"
		}
	}

	// Build tmux args
	var result []string

	// Global flags first (e.g., -f configfile)
	if configFile != "" {
		result = append(result, "-f", configFile)
	}

	switch operation {
	case "list":
		result = append(result, "list-sessions")
	case "attach":
		result = append(result, "attach")
		if detachFlag {
			result = append(result, "-d")
		}
		if sessionName != "" {
			result = append(result, "-t", sessionName)
		}
	case "reattach_or_create":
		result = append(result, "new-session", "-A")
		if sessionName != "" {
			result = append(result, "-s", sessionName)
		}
		if len(command) > 0 {
			result = append(result, "--")
			result = append(result, command...)
		}
	default: // "new"
		result = append(result, "new-session")
		if detachFlag && hasM {
			result = append(result, "-d")
		}
		if sessionName != "" {
			result = append(result, "-s", sessionName)
		}
		if len(command) > 0 {
			result = append(result, "--")
			result = append(result, command...)
		}
	}

	return result
}
