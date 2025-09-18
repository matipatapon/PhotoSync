package helper_test

import (
	"fmt"
	"os"
	"photosync/src/helper"
	"photosync/src/metadata"
	"reflect"
	"testing"
)

func TestThumbnailCreatorShouldCreateThumbnailForJpg(t *testing.T) {
	files := []string{"600_x_601", "601_x_600"}
	for _, file := range files {
		image, err := os.ReadFile(fmt.Sprintf("../../test/images/%s.jpg", file))
		if err != nil {
			t.FailNow()
		}
		sut := helper.ThumbnailCreator{}

		result, err := sut.Create(image, metadata.JPG)
		if err != nil || result == nil {
			t.FailNow()
		}

		thumbnail, err := os.ReadFile(fmt.Sprintf("../../test/images/%s_thumbnail.jpg", file))
		if err != nil {
			t.FailNow()
		}

		if !reflect.DeepEqual(thumbnail, result) {
			t.FailNow()
		}
	}
}

func TestThumbnailCreatorShouldNotCreateThumbnailForImageNotBiggerThan600x600(t *testing.T) {
	paths := []string{"../../test/images/600_x_600.jpg", "../../test/images/exif.jpg"}
	for _, path := range paths {
		image, err := os.ReadFile(path)
		if err != nil {
			t.FailNow()
		}
		sut := helper.ThumbnailCreator{}

		result, err := sut.Create(image, metadata.JPG)
		if err != nil || result != nil {
			t.FailNow()
		}
	}
}

func TestThumbnailCreatorShouldReturnErrorWhenFailedToDecodeFile(t *testing.T) {
	sut := helper.ThumbnailCreator{}

	result, err := sut.Create([]byte{0, 0, 0}, metadata.JPG)
	if err == nil || result != nil {
		t.FailNow()
	}
}

func TestThumbnailCreatorShouldReturnErrorForUnsupportedMIMEType(t *testing.T) {
	image, err := os.ReadFile("../../test/images/600_x_601.jpg")
	if err != nil {
		t.FailNow()
	}
	sut := helper.ThumbnailCreator{}

	result, err := sut.Create(image, metadata.UNKNOWN)
	if err == nil || result != nil {
		t.FailNow()
	}
}
