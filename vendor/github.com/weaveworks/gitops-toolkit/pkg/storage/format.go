package storage

// Format is an enum describing the format
// of data in a file. This is used by Storages
// to determine the encoding format for a file.
type Format uint8

const (
	FormatJSON Format = iota
	FormatYAML
)

// Formats describes the connection between
// file extensions and a encoding formats.
var Formats = map[string]Format{
	".json": FormatJSON,
	".yaml": FormatYAML,
	".yml":  FormatYAML,
}
