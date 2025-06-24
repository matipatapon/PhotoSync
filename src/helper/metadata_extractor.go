package helper

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
)

type IMetadataExtractor interface {
	Extract([]byte) (map[string]any, error)
}

type MetadataExtractor struct {
	logger *log.Logger
}

func NewMetadataExtractor() MetadataExtractor {
	return MetadataExtractor{logger: log.New(os.Stdout, "[MetadataExtractor]: ", log.LstdFlags)}
}

func (me *MetadataExtractor) Extract(file []byte) (map[string]any, error) {
	tmp, err := os.CreateTemp("", "")
	if err != nil {
		me.logger.Printf("Failed to create temporary file: '%s'", err.Error())
		return nil, err
	}
	filename := tmp.Name()
	tmp.Close()
	defer os.Remove(filename)

	err = os.WriteFile(filename, file, 0655)
	if err != nil {
		me.logger.Printf("Failed to write to temporary file: '%s'", err.Error())
		return nil, err
	}

	metadataBytes, err := exec.Command("perl", "exiftool/exiftool", "-j", filename).Output()
	if err != nil {
		me.logger.Printf("Failed to read metadata: '%s'", err.Error())
		return nil, err
	}

	var metadata []map[string]any
	err = json.Unmarshal(metadataBytes, &metadata)
	if err != nil {
		me.logger.Printf("Failed to convert metadata: '%s'", err.Error())
		return nil, err
	}
	return metadata[0], nil
}
