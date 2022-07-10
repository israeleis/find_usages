package models

import (
	"fmt"
	"github.com/thoas/go-funk"
	"strings"
)

type Usage struct {
	FilePath   string
	Line       string
	LineNumber int
	Index      int
}

func (u Usage) String() string {
	return fmt.Sprintf("%s %d:%d, %s", u.FilePath, u.LineNumber, u.Index, u.Line)
}

type UsageResult struct {
	FindValue string
	Usage     Usage
}

type UsagesResults struct {
	FindValue string
	Usages    []Usage
}

func (r *UsagesResults) AddResult(usage Usage) {
	r.Usages = append(r.Usages, usage)
}

func (r UsagesResults) String() string {
	res := fmt.Sprintf("%s %d results found\n", r.FindValue, len(r.Usages))
	usagesStrArr := funk.Map(r.Usages, func(u Usage) string {
		return u.String()
	}).([]string)
	res += strings.Join(usagesStrArr, "\n")
	return res
}
