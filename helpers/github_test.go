package helpers_test

import (
	"testing"

	"github.com/GoPex/unleash/helpers"
)

func TestEvaluateURL(t *testing.T) {
	url := "https://api.github.com/repos/GoPex/unleash_test_repository/{archive_format}{/ref}"
	expected := "https://api.github.com/repos/GoPex/unleash_test_repository/tarball/master"
	evaluatedUrl := helpers.EvaluateURL(url, "master")
	if evaluatedUrl != expected {
		t.Errorf("Github url evaluation failed, expected '%s', actual '%s'", expected, evaluatedUrl)
	}
}
