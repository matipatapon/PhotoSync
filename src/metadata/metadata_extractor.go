package metadata

import (
	"log"
	"os"
)

type MetadataExtractor struct {
	rme    IRawMetadataExtractor
	logger *log.Logger
}

func NewMetadataExtractor(rme IRawMetadataExtractor) MetadataExtractor {
	return MetadataExtractor{
		rme:    rme,
		logger: log.New(os.Stdout, "[MetadataExtractor]: ", log.LstdFlags),
	}
}

func (me *MetadataExtractor) Extract(file []byte) Metadata {
	meta, err := me.rme.Extract(file)
	if err != nil {
		me.logger.Printf("Failed to extract metadata from a file: '%s'", err.Error())
		return Metadata{}
	}
	return Metadata{
		CreationDate: me.extractCreationDate(meta),
		Location:     me.extractLocation(meta),
		MIMEType:     me.extractMIMeType(meta),
	}
}

func (me *MetadataExtractor) extractCreationDate(meta map[string]any) *Date {
	dateTags := []string{
		"Composite:DateTimeOriginal",
		"EXIF:DateTimeOriginal",
		"XMP:CreateDate",
		"QuickTime:CreateDate",
	}

	for _, dateTag := range dateTags {
		creationDateRaw, ok := meta[dateTag]
		if ok {
			creationDate, err := NewDate(creationDateRaw.(string))
			if err != nil {
				me.logger.Printf("'%s' contains invalid creation date '%s'", dateTag, creationDateRaw)
				continue
			}
			me.logger.Printf("Extracted creation date '%s' from '%s'", creationDateRaw, dateTag)
			return &creationDate
		}
	}
	me.logger.Print("Tag with creation date is missing")
	return nil
}

func (me *MetadataExtractor) extractLocation(meta map[string]any) *GPS {
	locationRaw, ok := meta["Composite:GPSPosition"]
	if ok {
		location, err := NewGPS(locationRaw.(string))
		if err == nil {
			me.logger.Printf("Extracted location '%s' from 'Composite:GPSPosition'", locationRaw)
			return &location
		}
		me.logger.Printf("Tag 'Composite:GPSPosition' contains invalid location '%s'", locationRaw)
	}
	me.logger.Printf("Tag with location is missing")
	return nil
}

func (me *MetadataExtractor) extractMIMeType(meta map[string]any) MIMEType {
	mimeType, ok := meta["File:MIMEType"]
	if ok {
		if mimeType == "image/jpeg" {
			me.logger.Printf("MIMEType is jpg")
			return JPG
		}
		me.logger.Printf("Unknown MIMEType")
	}
	me.logger.Printf("Tag with MIMEType is missing")
	return UNKNOWN
}
