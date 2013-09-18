package repo

// Structures that represent version data.

type FileSource struct {
    SourceType string
    Url string
}

type FileInfo struct {
    Path string
    Sources []FileSource
    Md5 string
}

type Version struct {
    ApiVersion int
    Id int
    Name string
    Files []FileInfo
}

