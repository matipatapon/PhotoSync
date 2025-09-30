package helper

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"log"
	"os"
	"photosync/src/metadata"

	"github.com/nfnt/resize"
)

const maxSizeWithoutAThumbnail = 600
const thumbnailSize = 400

type IThumbnailCreator interface {
	Create(file []byte, mimeType metadata.MIMEType) ([]byte, error)
}

type ThumbnailCreator struct {
	logger *log.Logger
}

func NewThumbnailCreator() ThumbnailCreator {
	return ThumbnailCreator{logger: log.New(os.Stdout, "[ThumbnailCreator]: ", log.LstdFlags)}
}

func (tc *ThumbnailCreator) Create(file []byte, mimeType metadata.MIMEType) ([]byte, error) {
	if mimeType == metadata.JPG {
		image, _, err := image.Decode(bytes.NewReader((file)))
		if err != nil {
			log.Printf("Failed to decode image: '%s'", err.Error())
			return nil, err
		}

		width := image.Bounds().Size().X
		height := image.Bounds().Size().Y
		if width <= maxSizeWithoutAThumbnail && height <= maxSizeWithoutAThumbnail {
			log.Printf("Image is not bigger than %d x %d ( %d x %d ), thumbnail will not be created", maxSizeWithoutAThumbnail, maxSizeWithoutAThumbnail, width, height)
			return nil, nil
		}

		thumbnail := resize.Thumbnail(thumbnailSize, thumbnailSize, image, resize.Lanczos3)
		result := bytes.Buffer{}
		err = jpeg.Encode(&result, thumbnail, nil)
		if err != nil {
			log.Printf("Failed to encode thumbnail: '%s'", err.Error())
			return nil, err
		}

		log.Print("Thumbnail was created")
		return result.Bytes(), nil
	}

	log.Printf("Unsupported MIMEType: '%s'", metadata.MIMETypeToString(mimeType))
	return nil, errors.New("unexpected MIMEType")
}
