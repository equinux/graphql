package location

import (
	"regexp"

	"github.com/chris-ramon/graphql-go/language/source"
)

type SourceLocation struct {
	Line   int
	Column int
}

func GetLocation(s *source.Source, position int) SourceLocation {
	line := 1
	column := position + 1
	lineRegexp := regexp.MustCompile("\r\n|[\n\r\u2028\u2029]")
	for _, match := range lineRegexp.FindAllStringIndex(s.Body, -1) {
		matchIndex := match[0]
		if matchIndex < position {
			line += 1
			l := len(s.Body[match[0]:match[1]])
			column = position + 1 - (matchIndex + l)
			continue
		} else {
			break
		}
	}
	return SourceLocation{Line: line, Column: column}
}