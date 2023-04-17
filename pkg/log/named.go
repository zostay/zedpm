package log

import (
	"fmt"
	"strings"
)

func Smprintf(f string, args ...any) string {
	names, values := makeNames(args)
	f = namedToPrintf(f, names)
	return fmt.Sprintf(f, values...)
}

func makeNames(args []any) (map[string]int, []any) {
	names := make(map[string]int, len(args)/2)
	values := make([]any, 0, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		if str, isStr := args[i].(string); isStr && i+1 < len(args) {
			names[str] = len(values) + 1
			values = append(values, args[i+1])
		}
	}
	return names, values
}

func namedToPrintf(f string, names map[string]int) string {
	rem := f
	outFmt := &strings.Builder{}
	for {
		// Look for %v substrings
		i := strings.IndexByte(rem, '%')
		if i < 0 {
			outFmt.WriteString(rem)
			return outFmt.String()
		}

		// Consume and ignore anything prior to the %v substring
		outFmt.WriteString(rem[:i])
		rem = rem[i:]

		// Copy %% through...
		if rem[1] == '%' {
			outFmt.WriteString(rem[:2])
			rem = rem[:2]
			continue
		}

		// Consume whatever is next until we reach [
		startBracket := -1
		moreWork := -1
		for i = 1; i < len(rem); i++ {
			if rem[i] == '[' {
				startBracket = i
				break
			}

			if rem[i] == '%' {
				moreWork = i
				break
			}
		}

		// We ran into a % before a [, so we will ignore this interpolation
		if moreWork >= 0 {
			outFmt.WriteString(rem[:moreWork])
			rem = rem[moreWork:]
			continue
		}

		// We didn't find a bracket or a %, so this string is done
		if startBracket < 0 {
			outFmt.WriteString(rem)
			return outFmt.String()
		}

		// We found the [ we're looking for, now find the ]
		stopBracket := -1
		for i = startBracket + 1; i < len(rem); i++ {
			if rem[i] == ']' {
				stopBracket = i
				break
			}

			if rem[i] == '%' {
				moreWork = i
				break
			}
		}

		// We ran into a % before a ], so we will ignore this interpolation
		if moreWork >= 0 {
			outFmt.WriteString(rem[:moreWork])
			rem = rem[moreWork:]
			continue
		}

		// We didn't find a bracket or a %, so this string is done
		if stopBracket < 0 {
			outFmt.WriteString(rem)
			return outFmt.String()
		}

		// We found both brackets, see if we can replace the contents with an index
		if index, hasIndex := names[rem[startBracket+1:stopBracket]]; hasIndex {
			outFmt.WriteString(rem[:startBracket+1])
			_, _ = fmt.Fprintf(outFmt, "%d", index)
			rem = rem[stopBracket:]
			continue
		}

		// If we're still here, we didn't find the name, so consume as-is and move on
		outFmt.WriteString(rem[:stopBracket+1])
		rem = rem[stopBracket+1:]
	}
}
