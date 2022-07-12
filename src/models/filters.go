package models

import (
	"github.com/thoas/go-funk"
	"regexp"
	"strings"
)

type FileNameFilter interface {
	Match(fileName string) bool
}

type FileNameFilterRegex struct {
	inp   string
	regex *regexp.Regexp
}

func (f FileNameFilterRegex) Match(fileName string) bool {
	//TODO implement me
	panic("implement me")
}

type FileNameFilterExtension struct {
	inp string
}

func (f FileNameFilterExtension) Match(fileName string) bool {
	return strings.HasSuffix(fileName, f.inp)
}

func newFileNameFilter(input string) FileNameFilter {
	trimmed := strings.TrimSpace(input)
	if strings.HasPrefix(trimmed, "/") && strings.HasSuffix(trimmed, "/") {
		rgx := regexp.MustCompile(trimmed)
		return FileNameFilterRegex{
			inp:   trimmed,
			regex: rgx,
		}
	} else {
		if !strings.HasPrefix(trimmed, ".") {
			trimmed = "." + trimmed
		}
		return FileNameFilterExtension{
			inp: trimmed,
		}
	}
}

func CreateFileNameFilters(filtersStrArr []string) []FileNameFilter {
	return funk.Map(filtersStrArr, func(filter string) FileNameFilter {
		return newFileNameFilter(filter)
	}).([]FileNameFilter)
}

type FileNameFilters struct {
	Include []FileNameFilter
	Exclude []FileNameFilter
}

func (fnf *FileNameFilters) IsRelevant(fileName string) bool {
	res := false
	for _, filter := range fnf.Include {
		if filter.Match(fileName) {
			res = true
			break
		}
	}

	if !res {
		return false
	}

	for _, filter := range fnf.Exclude {
		if filter.Match(fileName) {
			res = false
			break
		}
	}

	return res
}
