package models

import (
	"regexp"
	"strings"
)

type Matcher interface {
	FindUsage(line string, filePath string, lineNumber int) *Usage
	GetValue() string
}

type RegexMatcher struct {
	value string
	rgx   *regexp.Regexp
}

func NewRegexMatcher(value string) *RegexMatcher {
	rgx := regexp.MustCompile(value)
	return &RegexMatcher{
		value: value,
		rgx:   rgx,
	}
}

func (r RegexMatcher) GetValue() string {
	return r.value
}

func (r *RegexMatcher) FindUsage(line string, filePath string, lineNumber int) *Usage {
	matchInd := r.rgx.FindStringIndex(line)

	if matchInd == nil {
		return nil
	}
	return &Usage{
		FilePath:   filePath,
		Line:       line,
		LineNumber: lineNumber,
		Index:      matchInd[0],
	}
}

type ContainsMatcher struct {
	value string
}

func NewContainsMatcher(value string) *ContainsMatcher {
	return &ContainsMatcher{
		value: value,
	}
}

func (c *ContainsMatcher) GetValue() string {
	return c.value
}

func (c *ContainsMatcher) FindUsage(line string, filePath string, lineNumber int) *Usage {
	if !strings.Contains(line, c.value) {
		return nil
	}
	return &Usage{
		FilePath:   filePath,
		Line:       line,
		LineNumber: lineNumber,
		Index:      strings.Index(line, c.value),
	}
}

func CreateMatchers(findValues []string, isRegex bool) []Matcher {
	var res []Matcher
	for _, value := range findValues {
		if isRegex {
			res = append(res, NewRegexMatcher(value))
		} else {
			res = append(res, NewContainsMatcher(value))
		}
	}
	return res
}
