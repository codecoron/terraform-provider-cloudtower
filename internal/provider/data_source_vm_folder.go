package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-cloudtower/internal/cloudtower"
	"github.com/hashicorp/terraform-provider-cloudtower/internal/helper"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/vm_folder"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

func dataSourceVmFolder() *schema.Resource {

	return &schema.Resource{
		Description: "CloudTower vm folder data source.",

		ReadContext: dataSourceVmFolderRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name_in"},
				Description:   "filter vm folder by name.",
			},
			"name_in": {
				Type:          schema.TypeList,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				ConflictsWith: []string{"name"},
				Description:   "filter vm folder by name in list.",
			},
			"name_contains": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "filter vm folder by name contains.",
			},
			"cluster_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"cluster_id_in"},
				Description:   "filter vm folder by cluster id.",
			},
			"cluster_id_in": {
				Type:          schema.TypeList,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				ConflictsWith: []string{"cluster_id"},
				Description:   "filter vm folder by cluster id in list.",
			},
			"vm_folder": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of vm folder",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the vm folder.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the vm folder.",
						},
						"local_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The local id of the vm folder.",
						},
						"vm_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of VMs in the vm folder.",
						},
						"cluster": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The cluster of the vm folder.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The cluster name of the vm folder.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The cluster id of the vm folder.",
									},
								},
							},
						},
						"vms": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of VMs in the vm folder.",
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

func dataSourceVmFolderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	ct := meta.(*cloudtower.Client)

	gp := vm_folder.NewGetVMFoldersParams()
	gp.RequestBody = &models.GetVMFoldersRequestBody{
		Where: &models.VMFolderWhereInput{},
	}

	where, err := expandVmFolderWhereInput(d)
	if err != nil {
		return diag.FromErr(err)
	}
	gp.RequestBody.Where = where

	vm_folders, err := ct.Api.VMFolder.GetVMFolders(gp)
	if err != nil {
		return diag.FromErr(err)
	}
	output := make([]map[string]interface{}, 0)
	for _, f := range vm_folders.Payload {
		output = append(output, flattenVmFolder(f))
	}
	err = d.Set("vm_folder", output)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func expandVmFolderWhereInput(d *schema.ResourceData) (*models.VMFolderWhereInput, error) {
	where := &models.VMFolderWhereInput{}
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
	return where, nil
}

func flattenVmFolder(f *models.VMFolder) map[string]interface{} {
	res := map[string]interface{}{}

	if f.ID != nil {
		res["id"] = *f.ID
	}
	if f.Name != nil {
		res["name"] = *f.Name
	}
	if f.LocalID != nil {
		res["local_id"] = *f.LocalID
	}
	if f.VMNum != nil {
		res["vm_num"] = *f.VMNum
	}

	if f.Cluster != nil {
		cluster := map[string]interface{}{}
		if f.Cluster.ID != nil {
			cluster["id"] = *f.Cluster.ID
		}
		if f.Cluster.Name != nil {
			cluster["name"] = *f.Cluster.Name
		}
		res["cluster"] = []map[string]interface{}{cluster}
	}

	if len(f.Vms) > 0 {
		vms := make([]map[string]interface{}, 0, len(f.Vms))
		for _, v := range f.Vms {
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
