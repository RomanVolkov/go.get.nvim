package validations_test

import (
	"main/validations"
	"testing"
)

func TestIsValidPackageURL(t *testing.T) {
	validPackage := "github.com/go-chi/chi"
	invalidPackage := "github.com/GO-chi/chi"

	isValid, err := validations.IsPackageURLValid(validPackage)
	if !isValid {
		t.Errorf("Failed to validate valid package: %v\n", err)
	}

	isValid, err = validations.IsPackageURLValid(invalidPackage)
	if isValid {
		t.Errorf("Failed to validate invalid package: %v\n", err)
	}
}

func TestCleanupInvalidPackageURLs(t *testing.T) {
	uniqueURLs := map[string]bool{}
	uniqueURLs["github.com/go-chi/chi"] = false
	uniqueURLs["github.com/Go-chi/chi"] = false

	validations.CleanupInvalidPackageURLs(&uniqueURLs)

	if len(uniqueURLs) != 1 {
		t.Errorf("Failed to CleanupInvalidPackageURLs: %v\n", uniqueURLs)
	}

	for url := range uniqueURLs {
		if url != "github.com/go-chi/chi" {
			t.Errorf("Incorrect url:\n%v", url)
		}
	}
}
