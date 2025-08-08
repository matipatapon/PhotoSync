package metadata

type MIMEType int16

const (
	UNKNOWN MIMEType = iota
	JPG     MIMEType = iota
)

// TODO IMPROVE IT

func MIMETypeToString(mimeType MIMEType) string {
	if mimeType == JPG {
		return "image/jpeg"
	}
	return "unknown"
}

func StringToMIMEType(str string) MIMEType {
	if str == "image/jpeg" {
		return JPG
	}
	return UNKNOWN
}

type Metadata struct {
	CreationDate *Date
	MIMEType     MIMEType
}
