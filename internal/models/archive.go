package models

type Archive struct {
	ArchiveName string
	ArchiveSize int64
	TotalSize float64
	TotalFile float64
	Files []File
}

type File struct {
	FilePath string
	FileSize float64
	MimeType string
}