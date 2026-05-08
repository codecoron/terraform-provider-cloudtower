package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVmFolder_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmFolderConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_folder.all", "vm_folder.0.id"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_folder.all", "vm_folder.0.name"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_folder.all", "vm_folder.0.cluster.0.id"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_folder.all", "vm_folder.0.cluster.0.name"),
				),
			},
		},
	})
}

func TestAccDataSourceVmFolder_filterByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmFolderConfigFilterByName("jinlong.li"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cloudtower_vm_folder.by_name", "vm_folder.0.name", "jinlong.li"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_folder.by_name", "vm_folder.0.id"),
				),
			},
		},
	})
}

func TestAccDataSourceVmFolder_filterByNameContains(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmFolderConfigFilterByNameContains("test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_folder.by_name_contains", "vm_folder.0.id"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_folder.by_name_contains", "vm_folder.0.name"),
				),
			},
		},
	})
}

func TestAccDataSourceVmFolder_filterByClusterId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmFolderConfigFilterByClusterId(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_folder.by_cluster", "vm_folder.0.id"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_folder.by_cluster", "vm_folder.0.name"),
					resource.TestCheckResourceAttrPair(
						"data.cloudtower_vm_folder.by_cluster", "vm_folder.0.cluster.0.id",
						"data.cloudtower_cluster.test", "clusters.0.id",
					),
				),
			},
		},
	})
}

func testAccDataSourceVmFolderConfigBasic() string {
	return fmt.Sprintf(`
provider "cloudtower" {
  username          = "%s"
  password          = "%s"
  user_source       = "%s"
  cloudtower_server = "%s"
}

data "cloudtower_vm_folder" "all" {}
`, os.Getenv("CLOUDTOWER_USER"), os.Getenv("CLOUDTOWER_PASSWORD"), os.Getenv("CLOUDTOWER_USER_SOURCE"), os.Getenv("CLOUDTOWER_SERVER"))
}

func testAccDataSourceVmFolderConfigFilterByName(name string) string {
	return fmt.Sprintf(`
provider "cloudtower" {
  username          = "%s"
  password          = "%s"
  user_source       = "%s"
  cloudtower_server = "%s"
}

data "cloudtower_vm_folder" "by_name" {
  name = "%s"
}
`, os.Getenv("CLOUDTOWER_USER"), os.Getenv("CLOUDTOWER_PASSWORD"), os.Getenv("CLOUDTOWER_USER_SOURCE"), os.Getenv("CLOUDTOWER_SERVER"), name)
}

func testAccDataSourceVmFolderConfigFilterByNameContains(nameContains string) string {
	return fmt.Sprintf(`
provider "cloudtower" {
  username          = "%s"
  password          = "%s"
  user_source       = "%s"
  cloudtower_server = "%s"
}

data "cloudtower_vm_folder" "by_name_contains" {
  name_contains = "%s"
}
`, os.Getenv("CLOUDTOWER_USER"), os.Getenv("CLOUDTOWER_PASSWORD"), os.Getenv("CLOUDTOWER_USER_SOURCE"), os.Getenv("CLOUDTOWER_SERVER"), nameContains)
}

func testAccDataSourceVmFolderConfigFilterByClusterId() string {
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

data "cloudtower_vm_folder" "by_cluster" {
  cluster_id = data.cloudtower_cluster.test.clusters[0].id
}
`, os.Getenv("CLOUDTOWER_USER"), os.Getenv("CLOUDTOWER_PASSWORD"), os.Getenv("CLOUDTOWER_USER_SOURCE"), os.Getenv("CLOUDTOWER_SERVER"))
}
