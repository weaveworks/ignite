package storage

type Format uint8

var Formats = map[string]Format{
	".json": FormatJSON,
	".yaml": FormatYAML,
	".yml":  FormatYAML,
}

const (
	FormatJSON Format = iota
	FormatYAML
)
