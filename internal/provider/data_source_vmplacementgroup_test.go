package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVmPlacementGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmPlacementGroupConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_placement_group.all", "vm_placement_group.0.id"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_placement_group.all", "vm_placement_group.0.name"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_placement_group.all", "vm_placement_group.0.description"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_placement_group.all", "vm_placement_group.0.enabled"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_placement_group.all", "vm_placement_group.0.local_id"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_placement_group.all", "vm_placement_group.0.cluster.0.id"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_placement_group.all", "vm_placement_group.0.cluster.0.name"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_placement_group.all", "vm_placement_group.0.vm_host_must_enabled"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_placement_group.all", "vm_placement_group.0.vm_host_must_policy"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_placement_group.all", "vm_placement_group.0.vm_host_prefer_enabled"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_placement_group.all", "vm_placement_group.0.vm_host_prefer_policy"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_placement_group.all", "vm_placement_group.0.vm_vm_policy"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_placement_group.all", "vm_placement_group.0.vm_vm_policy_enabled"),
				),
			},
		},
	})
}

func TestAccDataSourceVmPlacementGroup_filterByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmPlacementGroupConfigFilterByName("web-app"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cloudtower_vm_placement_group.by_name", "vm_placement_group.0.name", "web-app"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_placement_group.by_name", "vm_placement_group.0.id"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_placement_group.by_name", "vm_placement_group.0.cluster.0.id"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_placement_group.by_name", "vm_placement_group.0.cluster.0.name"),
				),
			},
		},
	})
}

func TestAccDataSourceVmPlacementGroup_filterByNameContains(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmPlacementGroupConfigFilterByNameContains("web"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cloudtower_vm_placement_group.by_name_contains", "vm_placement_group.0.name", "web-app"),
				),
			},
		},
	})
}

func TestAccDataSourceVmPlacementGroup_filterByClusterId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmPlacementGroupConfigFilterByClusterId(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_placement_group.by_cluster", "vm_placement_group.0.id"),
					resource.TestCheckResourceAttrSet("data.cloudtower_vm_placement_group.by_cluster", "vm_placement_group.0.name"),
					resource.TestCheckResourceAttrPair(
						"data.cloudtower_vm_placement_group.by_cluster", "vm_placement_group.0.cluster.0.id",
						"data.cloudtower_cluster.test", "clusters.0.id",
					),
				),
			},
		},
	})
}

func testAccDataSourceVmPlacementGroupConfigBasic() string {
	return fmt.Sprintf(`
provider "cloudtower" {
  username          = "%s"
  password          = "%s"
  user_source       = "%s"
  cloudtower_server = "%s"
}

data "cloudtower_vm_placement_group" "all" {}
`, os.Getenv("CLOUDTOWER_USER"), os.Getenv("CLOUDTOWER_PASSWORD"), os.Getenv("CLOUDTOWER_USER_SOURCE"), os.Getenv("CLOUDTOWER_SERVER"))
}

func testAccDataSourceVmPlacementGroupConfigFilterByName(name string) string {
	return fmt.Sprintf(`
provider "cloudtower" {
  username          = "%s"
  password          = "%s"
  user_source       = "%s"
  cloudtower_server = "%s"
}

data "cloudtower_vm_placement_group" "by_name" {
  name = "%s"
}
`, os.Getenv("CLOUDTOWER_USER"), os.Getenv("CLOUDTOWER_PASSWORD"), os.Getenv("CLOUDTOWER_USER_SOURCE"), os.Getenv("CLOUDTOWER_SERVER"), name)
}

func testAccDataSourceVmPlacementGroupConfigFilterByNameContains(nameContains string) string {
	return fmt.Sprintf(`
provider "cloudtower" {
  username          = "%s"
  password          = "%s"
  user_source       = "%s"
  cloudtower_server = "%s"
}

data "cloudtower_vm_placement_group" "by_name_contains" {
  name_contains = "%s"
}
`, os.Getenv("CLOUDTOWER_USER"), os.Getenv("CLOUDTOWER_PASSWORD"), os.Getenv("CLOUDTOWER_USER_SOURCE"), os.Getenv("CLOUDTOWER_SERVER"), nameContains)
}

func testAccDataSourceVmPlacementGroupConfigFilterByClusterId() string {
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

data "cloudtower_vm_placement_group" "by_cluster" {
  cluster_id = data.cloudtower_cluster.test.clusters[0].id
}
`, os.Getenv("CLOUDTOWER_USER"), os.Getenv("CLOUDTOWER_PASSWORD"), os.Getenv("CLOUDTOWER_USER_SOURCE"), os.Getenv("CLOUDTOWER_SERVER"))
}
