package metadata

// This interface enables taking in any struct which embeds Metadata
type AnyMetadata interface {
	GetMD() *Metadata
}

// Verify that Metadata implements AnyMetadata
var _ AnyMetadata = &Metadata{}

func (md *Metadata) GetMD() *Metadata {
	return md
}
