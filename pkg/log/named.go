package log

import (
	"fmt"
	"strings"

	strings2 "github.com/zostay/go-std/strings"
)

func Smprintf(f string, args ...any) string {
	used := findUsedNames(f)
	values := makeValues(used, args)
	f = namedToPrintf(f, used)
	return fmt.Sprintf(f, values...)
}

type usedNames struct {
	name string
	s, e int
}

func findUsedNames(f string) []usedNames {
	rem := f
	consumed := 0
	used := make([]usedNames, 0, 4)
	for {
		// Look for %v substrings
		i := strings.IndexByte(rem, '%')
		if i < 0 {
			return used
		}

		// Consume and ignore anything prior to the %v substring
		consumed += i
		rem = rem[i:]

		// Copy %% through...
		if rem[1] == '%' {
			consumed += 2
			rem = rem[:2]
			continue
		}

		// Consume whatever is next until we reach [
		startBracket := -1
		contextStartBracket := -1
		moreWork := -1
		for i = 1; i < len(rem); i++ {
			if rem[i] == '[' {
				startBracket = i
				contextStartBracket = consumed + i + 1
				break
			}

			if rem[i] == '%' {
				moreWork = i
				break
			}
		}

		// We ran into a % before a [, so we will ignore this interpolation
		if moreWork >= 0 {
			consumed += moreWork
			rem = rem[moreWork:]
			continue
		}

		// We didn't find a bracket or a %, so this string is done
		if startBracket < 0 {
			return used
		}

		// We found the [ we're looking for, now find the ]
		stopBracket := -1
		contextStopBracket := -1
		for i = startBracket + 1; i < len(rem); i++ {
			if rem[i] == ']' {
				stopBracket = i
				contextStopBracket = consumed + i
				break
			}

			if rem[i] == '%' {
				moreWork = i
				break
			}
		}

		// We ran into a % before a ], so we will ignore this interpolation
		if moreWork >= 0 {
			consumed += moreWork
			rem = rem[moreWork:]
			continue
		}

		// We didn't find a bracket or a %, so this string is done
		if stopBracket < 0 {
			return used
		}

		// We found both brackets, see if we can replace the contents with an index
		name := rem[startBracket+1 : stopBracket]
		used = append(used, usedNames{
			name: name,
			s:    contextStartBracket,
			e:    contextStopBracket,
		})

		// If we're still here, we didn't find the name, so consume as-is and move on
		consumed += stopBracket + 1
		rem = rem[stopBracket+1:]
	}
}

func findNamed(name string, args []any) any {
	for i := 0; i < len(args); i += 2 {
		if str, isStr := args[i].(string); isStr && str == name && i+1 < len(args) {
			return args[i+1]
		}
	}
	return nil
}

func makeValues(used []usedNames, args []any) []any {
	values := make([]any, len(used))
	for i, u := range used {
		values[i] = findNamed(u.name, args)
	}
	return values
}

func namedToPrintf(f string, used []usedNames) string {
	bf := strings2.Reverse(f)
	outFmt := &strings.Builder{}
	from := 0
	for i := len(used) - 1; i >= 0; i-- {
		u := used[i]
		e := len(f) - u.e

		outFmt.WriteString(bf[from:e])
		fmt.Fprintf(outFmt, "%d", i+1)

		from = len(f) - u.s
	}
	outFmt.WriteString(bf[from:])
	return strings2.Reverse(outFmt.String())
}
