cluster_name = "CustomerService(SuperMicro-4-node)"
vds_name     = "mgmt"

vlans = {
  "web" = {
    name    = "tf-web-network"
    vlan_id = 100
  }
  "db" = {
    name    = "tf-db-network"
    vlan_id = 200
  }
  "cache" = {
    name    = "tf-cache-network"
    vlan_id = 399
  }
}
