package models

type Usage struct {
	FilePath   string
	Line       string
	LineNumber int
	Index      int
}

type UsageResult struct {
	FindValue string
	Usage     Usage
}
