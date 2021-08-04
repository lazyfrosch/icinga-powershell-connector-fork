package main

import (
	"regexp"
	"strings"
)

var (
	reUnquoteString = regexp.MustCompile(`(^["'\s]+|["'\s]+$)`)
	reArrayWrap     = regexp.MustCompile(`(^@\('?"?|'?"?\)$)`)
	reArraySplit    = regexp.MustCompile(`["'],[",]`)
)

// GetPowershellArgs returns remaining args and parse them as if we were powershell.exe
func GetPowershellArgs(args []string) (command string, arguments map[string]interface{}) {
	arguments = map[string]interface{}{}
	l := len(args)

	for i := 0; i < l; i++ {
		arg := args[i]

		// ignore any flag with double dash (our flags)
		if strings.HasPrefix(arg, "--") {
			if i+1 < l && args[i+1][0] != '-' {
				i++ // skip the next one also
			}

			continue
		}

		// retrieve command
		if strings.EqualFold(arg, "-Command") || strings.EqualFold(arg, "-C") {
			if i+1 < l && args[i+1][0] != '-' {
				command = args[i+1]
				i++
			}

			continue
		}

		// all other flags
		if i+1 >= l || args[i+1][0] == '-' {
			// next argument is also a flag, so this is a switch
			arguments[arg] = true
		} else {
			arguments[arg] = BuildPowershellType(args[i+1])
			i++
		}
	}

	if command != "" {
		command = ParsePowershellTryCatch(command)
	}

	return
}

func BuildPowershellType(value string) interface{} {
	if strings.EqualFold(value, `$null`) {
		return nil
	} else if strings.EqualFold(value, `$true`) {
		return true
	} else if strings.EqualFold(value, `$false`) {
		return false
	} else if IsPowershellArray(value) {
		return ConvertPowershellArray(value)
	} else {
		return value
	}
}

// ConvertPowershellArray to a golang type.
//
// Examples:
//  @() -> []string{}
//  @('abc') -> []string{"abc"}
//  @('abc','def') -> []string{"abc","def}
//
func ConvertPowershellArray(value string) []string {
	value = reArrayWrap.ReplaceAllString(value, "")

	if value == "" {
		return []string{}
	} else {
		return reArraySplit.Split(value, -1)
	}
}

// ParsePowershellTryCatch parses the actual command from a try/catch PowerShell code snippet.
//
// Examples:
//
//  try { Use-Icinga -Minimal; } catch { <# something #> exit 3; };
// 	  Exit-IcingaExecutePlugin -Command 'Invoke-IcingaCheckUsedPartitionSpace'
//  try { Use-Icinga -Minimal; } catch { <# something #> exit 3; }; Invoke-IcingaCheckUsedPartitionSpace
//  Invoke-IcingaCheckUsedPartitionSpace
//
func ParsePowershellTryCatch(command string) string {
	command = strings.TrimSpace(command)

	// for now just parse the last word, dequote it and use it as command
	parts := strings.Split(command, " ")
	command = parts[len(parts)-1]
	command = reUnquoteString.ReplaceAllString(command, "")

	return command
}

func IsPowershellArray(s string) bool {
	l := len(s)
	if l < 3 {
		return false
	}

	return s[0] == '@' && s[1] == '(' && s[l-1] == ')'
}
