package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-provider-cloudtower/internal/cloudtower"
	"github.com/hashicorp/terraform-provider-cloudtower/internal/helper"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/vm_placement_group"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVmPlacementGroup() *schema.Resource {

	return &schema.Resource{
		Description: "CloudTower vm placement group data source.",

		ReadContext: dataSourceVmPlacementGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "filter vm placement group by name.",
			},
			"name_in": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "filter vm placement group by name in list.",
			},
			"name_contains": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "filter vm placement group by name contains.",
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

	if name := d.Get("name").(string); name != "" {
		gp.RequestBody.Where.Name = &name
	} else {
		nameIn, err := helper.SliceInterfacesToTypeSlice[string](d.Get("name_in").([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		} else if len(nameIn) > 0 {
			gp.RequestBody.Where.NameIn = nameIn
		}
	}
	if nameContains := d.Get("name_contains").(string); nameContains != "" {
		gp.RequestBody.Where.NameContains = &nameContains
	}

	vm_placement_groups, err := ct.Api.VMPlacementGroup.GetVMPlacementGroups(gp)
	if err != nil {
		return diag.FromErr(err)
	}
	output := make([]map[string]interface{}, 0)
	for _, d := range vm_placement_groups.Payload {
		output = append(output, map[string]interface{}{
			"id":   d.ID,
			"name": d.Name,
			"cluster": []map[string]interface{}{
				{
					"name": d.Cluster.Name,
					"id":   d.Cluster.ID,
				},
			},
		})
	}
	err = d.Set("vm_placement_groups", output)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
