package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProvider *schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = New("")()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"cloudtower": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CLOUDTOWER_SERVER"); v == "" {
		t.Fatal("CLOUDTOWER_SERVER must be set for acceptance tests")
	}
	if v := os.Getenv("CLOUDTOWER_USER"); v == "" {
		t.Fatal("CLOUDTOWER_USER must be set for acceptance tests")
	}
	if v := os.Getenv("CLOUDTOWER_PASSWORD"); v == "" {
		t.Fatal("CLOUDTOWER_PASSWORD must be set for acceptance tests")
	}
}
