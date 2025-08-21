package helper_test

import (
	"os"
	"photosync/src/helper"
	"testing"
)

var ENV_NAME string = "ENV_NAME"
var ENV_VALUE string = "ENV_VALUE"

func TestEnvGetterShouldReturnEnviornmentVariable(t *testing.T) {
	os.Setenv(ENV_NAME, ENV_VALUE)
	sut := helper.NewEnvGetter()

	result := sut.Get(ENV_NAME)

	if result != ENV_VALUE {
		t.Errorf("'%s' != '%s'\n", result, ENV_VALUE)
	}
}
