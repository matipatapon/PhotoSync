package metadata

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
)

type IRawMetadataExtractor interface {
	Extract([]byte) (map[string]any, error)
}

type RawMetadataExtractor struct {
	exiftoolPath string
	logger       *log.Logger
}

func NewRawMetadataExtractor(exiftoolPath string) RawMetadataExtractor {

	return RawMetadataExtractor{
		exiftoolPath: exiftoolPath,
		logger:       log.New(os.Stdout, "[RawMetadataExtractor]: ", log.LstdFlags)}
}

func (me *RawMetadataExtractor) Extract(file []byte) (map[string]any, error) {
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

	metadataBytes, err := exec.Command(
		"perl",
		me.exiftoolPath,
		"-j",                // output in json
		"-G0",               // add group prefix to key e.g. EXIF:GPSLongitude
		"-c",                // format gps
		"%d %d %.2f",        // gps format
		"-d",                // format date
		"%Y.%m.%d %H:%M:%S", // date format
		filename).Output()
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

	if len(metadata) != 1 {
		me.logger.Print("File has no metadata")
		return make(map[string]any), nil
	}

	me.logger.Printf("Successfully harvested metadata from file")

	return metadata[0], nil
}
