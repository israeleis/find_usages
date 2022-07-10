package models

type FlowChannels struct {
	Errors      chan error
	Directories chan string
	FilePath    chan string
	FileContent chan FileContent
	UsageResult chan UsageResult
}

func (c FlowChannels) Close() {
	close(c.Errors)
	close(c.Directories)
	close(c.FilePath)
	close(c.FileContent)
	close(c.UsageResult)
}
