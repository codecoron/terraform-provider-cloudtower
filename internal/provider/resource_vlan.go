package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-cloudtower/internal/cloudtower"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/vlan"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

func resourceVlan() *schema.Resource {
	return &schema.Resource{
		Description: "CloudTower VM VLAN resource.",

		CreateContext: resourceVlanCreate,
		ReadContext:   resourceVlanRead,
		UpdateContext: resourceVlanUpdate,
		DeleteContext: resourceVlanDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VLAN's name",
			},
			"vds_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the VDS this VLAN belongs to",
			},
			"vlan_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 4095),
				Description:  "VLAN ID (0-4095)",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VLAN's ID",
			},
			"local_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VLAN's local ID",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VLAN's network type",
			},
		},
	}
}

func resourceVlanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ct := meta.(*cloudtower.Client)
	cvp := vlan.NewCreateVMVlanParams()
	name := d.Get("name").(string)
	vdsID := d.Get("vds_id").(string)
	vid := models.VlanID(d.Get("vlan_id").(int))
	params := &models.VMVlanCreationParams{
		Name:   &name,
		VdsID:  &vdsID,
		VlanID: &vid,
	}
	cvp.RequestBody = []*models.VMVlanCreationParams{params}
	vlans, err := ct.Api.Vlan.CreateVMVlan(cvp)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(*vlans.Payload[0].Data.ID)
	taskIds := make([]string, 0)
	for _, v := range vlans.Payload {
		if v.TaskID != nil {
			taskIds = append(taskIds, *v.TaskID)
		}
	}
	_, err = ct.WaitTasksFinish(ctx, taskIds)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceVlanRead(ctx, d, meta)
}

func resourceVlanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ct := meta.(*cloudtower.Client)

	id := d.Id()
	gvp := vlan.NewGetVlansParams()
	gvp.RequestBody = &models.GetVlansRequestBody{
		Where: &models.VlanWhereInput{
			ID: &id,
		},
	}
	vlans, err := ct.Api.Vlan.GetVlans(gvp)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(vlans.Payload) < 1 {
		d.SetId("")
		return diags
	}
	v := vlans.Payload[0]
	if err = d.Set("name", v.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("local_id", v.LocalID); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("type", string(*v.Type)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("vlan_id", int(*v.VlanID)); err != nil {
		return diag.FromErr(err)
	}
	if v.Vds != nil && v.Vds.ID != nil {
		if err = d.Set("vds_id", *v.Vds.ID); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceVlanUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ct := meta.(*cloudtower.Client)
	uvp := vlan.NewUpdateVlanParams()
	id := d.Id()
	data := &models.VMVlanUpdationParamsData{}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		data.Name = &name
	}
	if d.HasChange("vlan_id") {
		vid := models.VlanID(d.Get("vlan_id").(int))
		data.VlanID = &vid
	}
	uvp.RequestBody = &models.VMVlanUpdationParams{
		Where: &models.VlanWhereInput{
			ID: &id,
		},
		Data: data,
	}
	vlans, err := ct.Api.Vlan.UpdateVlan(uvp)
	if err != nil {
		return diag.FromErr(err)
	}
	taskIds := make([]string, 0)
	for _, v := range vlans.Payload {
		if v.TaskID != nil {
			taskIds = append(taskIds, *v.TaskID)
		}
	}
	_, err = ct.WaitTasksFinish(ctx, taskIds)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceVlanRead(ctx, d, meta)
}

func resourceVlanDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ct := meta.(*cloudtower.Client)
	dvp := vlan.NewDeleteVlanParams()
	id := d.Id()
	dvp.RequestBody = &models.VlanDeletionParams{
		Where: &models.VlanWhereInput{
			ID: &id,
		},
	}
	vlans, err := ct.Api.Vlan.DeleteVlan(dvp)
	if err != nil {
		return diag.FromErr(err)
	}
	taskIds := make([]string, 0)
	for _, v := range vlans.Payload {
		if v.TaskID != nil {
			taskIds = append(taskIds, *v.TaskID)
		}
	}
	_, err = ct.WaitTasksFinish(ctx, taskIds)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}

func resourceVlanImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ct := meta.(*cloudtower.Client)
	id := d.Id()
	gvp := vlan.NewGetVlansParams()
	gvp.RequestBody = &models.GetVlansRequestBody{
		Where: &models.VlanWhereInput{
			ID: &id,
		},
	}
	vlans, err := ct.Api.Vlan.GetVlans(gvp)
	if err != nil {
		return nil, err
	}
	if len(vlans.Payload) < 1 {
		return nil, fmt.Errorf("vlan with id %s not found", id)
	}
	return []*schema.ResourceData{d}, nil
}
