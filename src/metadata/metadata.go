package metadata

type MIMEType int

const (
	UNKNOWN MIMEType = iota
	JPG     MIMEType = iota
)

type Metadata struct {
	CreationDate *Date
	MIMEType     MIMEType
}
