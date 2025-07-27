package metadata

/*
TODO Error handling !!!!!!!!
*/

type MetadataExtractor struct {
	rme IRawMetadataExtractor
}

func NewMetadataExtractor(rme IRawMetadataExtractor) MetadataExtractor {
	return MetadataExtractor{rme: rme}
}

func (me *MetadataExtractor) Extract(file []byte) Metadata {
	meta, _ := me.rme.Extract(file)
	return Metadata{CreationDate: extractCreationDate(meta), Location: extractLocation(meta), MIMEType: extractMIMeType(meta)}
}

func extractCreationDate(meta map[string]any) *Date {
	dateTags := []string{
		"Composite:DateTimeOriginal",
		"EXIF:DateTimeOriginal",
		"XMP:CreateDate",
		"QuickTime:CreateDate",
	}

	for _, dateTag := range dateTags {
		creationDateRaw, ok := meta[dateTag]
		if ok {
			creationDate, _ := NewDate(creationDateRaw.(string))
			return &creationDate
		}
	}
	return nil
}

func extractLocation(meta map[string]any) *GPS {
	locationRaw, ok := meta["Composite:GPSPosition"]
	if ok {
		location, _ := NewGPS(locationRaw.(string))
		return &location
	}
	return nil
}

func extractMIMeType(meta map[string]any) MIMEType {
	mimeType, ok := meta["File:MIMEType"]
	if ok && mimeType == "image/jpeg" {
		return JPG
	}
	return UNKNOWN
}
