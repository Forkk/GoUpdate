package repo

// Structures that represent a GoUpdate index file.

const IndexFileName = "index.json"

type VersionSummary struct {
    Id int
    Name string
}

type ChannelSummary struct {
    Id string
    Name string
}

type Index struct {
    ApiVersion int
    Versions []VersionSummary
    Channels []ChannelSummary
}

// Returns a new, blank index struct.
func NewBlankIndex() Index {
    return Index{ApiVersion: 0, Versions: []VersionSummary {}, Channels: []ChannelSummary {}}
}

