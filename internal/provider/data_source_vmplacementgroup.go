package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-cloudtower/internal/cloudtower"
	"github.com/hashicorp/terraform-provider-cloudtower/internal/helper"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/vm_placement_group"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

func dataSourceVmPlacementGroup() *schema.Resource {

	return &schema.Resource{
		Description: "CloudTower vm placement group data source.",

		ReadContext: dataSourceVmPlacementGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name_in"},
				Description:   "filter vm placement group by name.",
			},
			"name_in": {
				Type:          schema.TypeList,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				ConflictsWith: []string{"name"},
				Description:   "filter vm placement group by name in list.",
			},
			"name_contains": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "filter vm placement group by name contains.",
			},
			"cluster_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"cluster_id_in"},
				Description:   "filter vm placement group by cluster id.",
			},
			"cluster_id_in": {
				Type:          schema.TypeList,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				ConflictsWith: []string{"cluster_id"},
				Description:   "filter vm placement group by cluster id in list.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "filter vm placement group by enabled.",
			},
			"vm_placement_group": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of vm placement group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the vm placement group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the vm placement group.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the vm placement group.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the vm placement group is enabled.",
						},
						"local_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The local id of the vm placement group.",
						},
						"cluster": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The cluster of the vm placement group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The cluster name of the vm placement group.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The cluster id of the vm placement group.",
									},
								},
							},
						},
						"vm_host_must_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether vm host must policy is enabled.",
						},
						"vm_host_must_host_uuids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of hosts that VMs must be placed on.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The host id.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The host name.",
									},
									"management_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The host management ip.",
									},
								},
							},
						},
						"vm_host_must_policy": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether vm host must policy is set.",
						},
						"vm_host_prefer_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether vm host prefer policy is enabled.",
						},
						"vm_host_prefer_host_uuids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of hosts that VMs prefer to be placed on.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The host id.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The host name.",
									},
									"management_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The host management ip.",
									},
								},
							},
						},
						"vm_host_prefer_policy": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether vm host prefer policy is set.",
						},
						"vm_vm_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The vm vm policy of the vm placement group.",
						},
						"vm_vm_policy_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether vm vm policy is enabled.",
						},
						"vms": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of VMs in the vm placement group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The vm id.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The vm name.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceVmPlacementGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	ct := meta.(*cloudtower.Client)

	gp := vm_placement_group.NewGetVMPlacementGroupsParams()
	gp.RequestBody = &models.GetVMPlacementGroupsRequestBody{
		Where: &models.VMPlacementGroupWhereInput{},
	}

	where, err := expandVmPlacementGroupWhereInput(d)
	if err != nil {
		return diag.FromErr(err)
	}
	gp.RequestBody.Where = where

	vm_placement_groups, err := ct.Api.VMPlacementGroup.GetVMPlacementGroups(gp)
	if err != nil {
		return diag.FromErr(err)
	}
	output := make([]map[string]interface{}, 0)
	for _, g := range vm_placement_groups.Payload {
		output = append(output, flattenVmPlacementGroup(g))
	}
	err = d.Set("vm_placement_group", output)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func expandVmPlacementGroupWhereInput(d *schema.ResourceData) (*models.VMPlacementGroupWhereInput, error) {
	where := &models.VMPlacementGroupWhereInput{}
	if name := d.Get("name").(string); name != "" {
		where.Name = &name
	} else {
		nameIn, err := helper.SliceInterfacesToTypeSlice[string](d.Get("name_in").([]interface{}))
		if err != nil {
			return nil, err
		} else if len(nameIn) > 0 {
			where.NameIn = nameIn
		}
	}
	if nameContains := d.Get("name_contains").(string); nameContains != "" {
		where.NameContains = &nameContains
	}
	if clusterId := d.Get("cluster_id").(string); clusterId != "" {
		where.Cluster = &models.ClusterWhereInput{
			ID: &clusterId,
		}
	} else {
		clusterIdIn, err := helper.SliceInterfacesToTypeSlice[string](d.Get("cluster_id_in").([]interface{}))
		if err != nil {
			return nil, err
		} else if len(clusterIdIn) > 0 {
			where.Cluster = &models.ClusterWhereInput{
				IDIn: clusterIdIn,
			}
		}
	}
	if enabled, ok := d.GetOkExists("enabled"); ok {
		e := enabled.(bool)
		where.Enabled = &e
	}
	return where, nil
}

func flattenVmPlacementGroup(g *models.VMPlacementGroup) map[string]interface{} {
	res := map[string]interface{}{}

	if g.ID != nil {
		res["id"] = *g.ID
	}
	if g.Name != nil {
		res["name"] = *g.Name
	}
	if g.Description != nil {
		res["description"] = *g.Description
	}
	if g.Enabled != nil {
		res["enabled"] = *g.Enabled
	}
	if g.LocalID != nil {
		res["local_id"] = *g.LocalID
	}
	if g.VMHostMustEnabled != nil {
		res["vm_host_must_enabled"] = *g.VMHostMustEnabled
	}
	if g.VMHostMustPolicy != nil {
		res["vm_host_must_policy"] = *g.VMHostMustPolicy
	}
	if g.VMHostPreferEnabled != nil {
		res["vm_host_prefer_enabled"] = *g.VMHostPreferEnabled
	}
	if g.VMHostPreferPolicy != nil {
		res["vm_host_prefer_policy"] = *g.VMHostPreferPolicy
	}
	if g.VMVMPolicy != nil {
		res["vm_vm_policy"] = string(*g.VMVMPolicy)
	}
	if g.VMVMPolicyEnabled != nil {
		res["vm_vm_policy_enabled"] = *g.VMVMPolicyEnabled
	}

	if g.Cluster != nil {
		cluster := map[string]interface{}{}
		if g.Cluster.ID != nil {
			cluster["id"] = *g.Cluster.ID
		}
		if g.Cluster.Name != nil {
			cluster["name"] = *g.Cluster.Name
		}
		res["cluster"] = []map[string]interface{}{cluster}
	}

	if len(g.VMHostMustHostUuids) > 0 {
		hosts := make([]map[string]interface{}, 0, len(g.VMHostMustHostUuids))
		for _, h := range g.VMHostMustHostUuids {
			host := map[string]interface{}{}
			if h.ID != nil {
				host["id"] = *h.ID
			}
			if h.Name != nil {
				host["name"] = *h.Name
			}
			if h.ManagementIP != nil {
				host["management_ip"] = *h.ManagementIP
			}
			hosts = append(hosts, host)
		}
		res["vm_host_must_host_uuids"] = hosts
	}

	if len(g.VMHostPreferHostUuids) > 0 {
		hosts := make([]map[string]interface{}, 0, len(g.VMHostPreferHostUuids))
		for _, h := range g.VMHostPreferHostUuids {
			host := map[string]interface{}{}
			if h.ID != nil {
				host["id"] = *h.ID
			}
			if h.Name != nil {
				host["name"] = *h.Name
			}
			if h.ManagementIP != nil {
				host["management_ip"] = *h.ManagementIP
			}
			hosts = append(hosts, host)
		}
		res["vm_host_prefer_host_uuids"] = hosts
	}

	if len(g.Vms) > 0 {
		vms := make([]map[string]interface{}, 0, len(g.Vms))
		for _, v := range g.Vms {
			vm := map[string]interface{}{}
			if v.ID != nil {
				vm["id"] = *v.ID
			}
			if v.Name != nil {
				vm["name"] = *v.Name
			}
			vms = append(vms, vm)
		}
		res["vms"] = vms
	}

	return res
}
