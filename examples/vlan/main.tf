terraform {
  required_providers {
    cloudtower = {
      source = "registry.terraform.io/smartxworks/cloudtower"
    }
  }
}

provider "cloudtower" {
  username          = "terraform"
  user_source       = "LOCAL"
  cloudtower_server = "192.68.190.1"
}

# 查询目标集群
data "cloudtower_cluster" "target_cluster" {
  name = var.cluster_name
}

# 查询 VDS
data "cloudtower_vds" "test" {
  name       = var.vds_name
  cluster_id = data.cloudtower_cluster.target_cluster.clusters[0].id
}

# 用例 1：基本创建，只传 name + vds_id，系统自动分配 vlan_id
resource "cloudtower_vlan" "basic" {
  name   = "tf-example-vlan-basic"
  vds_id = data.cloudtower_vds.test.vds[0].id
}

# 用例 2：显式指定 vlan_id
resource "cloudtower_vlan" "with_vlan_id" {
  name    = "tf-example-vlan-with-id"
  vds_id  = data.cloudtower_vds.test.vds[0].id
  vlan_id = 200
}

# 用例 3：批量创建多个 VLAN
resource "cloudtower_vlan" "multi_vlan_ids" {
  for_each = var.vlans

  name    = each.value.name
  vds_id  = data.cloudtower_vds.test.vds[0].id
  vlan_id = each.value.vlan_id
}

output "vlan_basic" {
  value = {
    id       = cloudtower_vlan.basic.id
    name     = cloudtower_vlan.basic.name
    local_id = cloudtower_vlan.basic.local_id
    type     = cloudtower_vlan.basic.type
    vlan_id  = cloudtower_vlan.basic.vlan_id
  }
}

output "vlan_with_vlan_id" {
  value = {
    id       = cloudtower_vlan.with_vlan_id.id
    name     = cloudtower_vlan.with_vlan_id.name
    local_id = cloudtower_vlan.with_vlan_id.local_id
    type     = cloudtower_vlan.with_vlan_id.type
    vlan_id  = cloudtower_vlan.with_vlan_id.vlan_id
  }
}

output "multi_vlan_ids" {
  value = { for k, v in cloudtower_vlan.multi_vlan_ids : k => {
    id       = v.id
    name     = v.name
    local_id = v.local_id
    vlan_id  = v.vlan_id
  } }
}
