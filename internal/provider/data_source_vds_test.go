package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVds_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVdsConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cloudtower_vds.all", "vds.0.id"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vds.all", "vds.0.name"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vds.all", "vds.0.type"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vds.all", "vds.0.bond_mode"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vds.all", "vds.0.ovsbr_name"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vds.all", "vds.0.internal"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vds.all", "vds.0.vlans_num"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vds.all", "vds.0.cluster.0.id"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vds.all", "vds.0.cluster.0.name"),
				),
			},
		},
	})
}

func TestAccDataSourceVds_filterByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVdsConfigFilterByName("management"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cloudtower_vds.by_name", "vds.0.name", "management"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vds.by_name", "vds.0.id"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vds.by_name", "vds.0.cluster.0.id"),
				),
			},
		},
	})
}

func TestAccDataSourceVds_filterByNameContains(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVdsConfigFilterByNameContains("manage"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cloudtower_vds.by_name_contains", "vds.0.id"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vds.by_name_contains", "vds.0.name"),
				),
			},
		},
	})
}

func TestAccDataSourceVds_filterByClusterId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVdsConfigFilterByClusterId(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cloudtower_vds.by_cluster", "vds.0.id"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vds.by_cluster", "vds.0.name"),
					resource.TestCheckResourceAttrPair(
						"data.cloudtower_vds.by_cluster", "vds.0.cluster.0.id",
						"data.cloudtower_cluster.test", "clusters.0.id",
					),
				),
			},
		},
	})
}

func testAccDataSourceVdsConfigBasic() string {
	return fmt.Sprintf(`
provider "cloudtower" {
  username          = "%s"
  password          = "%s"
  user_source       = "%s"
  cloudtower_server = "%s"
}

data "cloudtower_vds" "all" {}
`, os.Getenv("CLOUDTOWER_USER"), os.Getenv("CLOUDTOWER_PASSWORD"), os.Getenv("CLOUDTOWER_USER_SOURCE"), os.Getenv("CLOUDTOWER_SERVER"))
}

func testAccDataSourceVdsConfigFilterByName(name string) string {
	return fmt.Sprintf(`
provider "cloudtower" {
  username          = "%s"
  password          = "%s"
  user_source       = "%s"
  cloudtower_server = "%s"
}

data "cloudtower_vds" "by_name" {
  name = "%s"
}
`, os.Getenv("CLOUDTOWER_USER"), os.Getenv("CLOUDTOWER_PASSWORD"), os.Getenv("CLOUDTOWER_USER_SOURCE"), os.Getenv("CLOUDTOWER_SERVER"), name)
}

func testAccDataSourceVdsConfigFilterByNameContains(nameContains string) string {
	return fmt.Sprintf(`
provider "cloudtower" {
  username          = "%s"
  password          = "%s"
  user_source       = "%s"
  cloudtower_server = "%s"
}

data "cloudtower_vds" "by_name_contains" {
  name_contains = "%s"
}
`, os.Getenv("CLOUDTOWER_USER"), os.Getenv("CLOUDTOWER_PASSWORD"), os.Getenv("CLOUDTOWER_USER_SOURCE"), os.Getenv("CLOUDTOWER_SERVER"), nameContains)
}

func testAccDataSourceVdsConfigFilterByClusterId() string {
	return fmt.Sprintf(`
provider "cloudtower" {
  username          = "%s"
  password          = "%s"
  user_source       = "%s"
  cloudtower_server = "%s"
}

data "cloudtower_cluster" "test" {
  name = "CustomerService(SuperMicro-4-node)"
}

data "cloudtower_vds" "by_cluster" {
  cluster_id = data.cloudtower_cluster.test.clusters[0].id
}
`, os.Getenv("CLOUDTOWER_USER"), os.Getenv("CLOUDTOWER_PASSWORD"), os.Getenv("CLOUDTOWER_USER_SOURCE"), os.Getenv("CLOUDTOWER_SERVER"))
}
