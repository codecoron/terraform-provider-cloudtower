variable "cluster_name" {
  description = "目标集群名称"
  type        = string
}

variable "vds_name" {
  description = "VDS 名称"
  type        = string
}

variable "vlans" {
  description = "VLAN 配置映射"
  type = map(object({
    name    = string
    vlan_id = number
  }))
}
