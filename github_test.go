package unleash_test

import (
	"testing"

	"github.com/GoPex/unleash"
)

func TestEvaluateUrl(t *testing.T) {
	url := "https://api.github.com/repos/GoPex/unleash_test_repository/{archive_format}{/ref}"
	expected := "https://api.github.com/repos/GoPex/unleash_test_repository/tarball/master"
	evaluatedUrl := unleash.EvaluateUrl(url, "master")
	if evaluatedUrl != expected {
		t.Errorf("Github url evaluation failed, expected '%s', actual '%s'", expected, evaluatedUrl)
	}
}
