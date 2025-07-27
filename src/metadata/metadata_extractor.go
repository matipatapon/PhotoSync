package metadata

type MetadataExtractor struct {
	rme IRawMetadataExtractor
}

func NewMetadataExtractor(rme IRawMetadataExtractor) MetadataExtractor {
	return MetadataExtractor{rme: rme}
}

func (me *MetadataExtractor) Extract(file []byte) *Metadata {
	meta, _ := me.rme.Extract(file)

	if len(meta) == 0 {
		return nil
	}

	location, _ := NewGPS(meta["Composite:GPSPosition"].(string))

	return &Metadata{CreationDate: extractCreationDate(meta), Location: location, MIMEType: JPG}
}

func extractCreationDate(meta map[string]any) *Date {
	dateTags := []string{
		"EXIF:DateTimeOriginal",
		"XMP:CreateDate",
		"QuickTime:CreateDate",
		"Composite:DateTimeOriginal",
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
