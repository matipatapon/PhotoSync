package metadata

type MIMEType int

const (
	JPG     MIMEType = iota
	UNKNOWN MIMEType = iota
)

type Metadata struct {
	CreationDate *Date
	Location     GPS
	MIMEType     MIMEType
}
