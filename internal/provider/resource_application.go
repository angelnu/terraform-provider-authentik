package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationCreate,
		ReadContext:   resourceApplicationRead,
		UpdateContext: resourceApplicationUpdate,
		DeleteContext: resourceApplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol_provider": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"meta_launch_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"meta_icon": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"meta_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"meta_publisher": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_engine_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.POLICYENGINEMODE_ANY,
			},
			"open_in_new_tab": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceApplicationSchemaToModel(d *schema.ResourceData) *api.ApplicationRequest {
	m := api.ApplicationRequest{
		Name:         d.Get("name").(string),
		Slug:         d.Get("slug").(string),
		Provider:     api.NullableInt32{},
		OpenInNewTab: boolToPointer(d.Get("open_in_new_tab").(bool)),
	}

	if p, pSet := d.GetOk("protocol_provider"); pSet {
		m.Provider.Set(intToPointer(p.(int)))
	} else {
		m.Provider.Set(nil)
	}

	if l, ok := d.Get("group").(string); ok {
		m.Group = &l
	}
	if l, ok := d.Get("meta_launch_url").(string); ok {
		m.MetaLaunchUrl = &l
	}
	if l, ok := d.Get("meta_description").(string); ok {
		m.MetaDescription = &l
	}
	if l, ok := d.Get("meta_publisher").(string); ok {
		m.MetaPublisher = &l
	}

	pm := api.PolicyEngineMode(d.Get("policy_engine_mode").(string))
	m.PolicyEngineMode = &pm
	return &m
}

func resourceApplicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceApplicationSchemaToModel(d)

	res, hr, err := c.client.CoreApi.CoreApplicationsCreate(ctx).ApplicationRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)

	if i, iok := d.GetOk("meta_icon"); iok {
		hr, err := c.client.CoreApi.CoreApplicationsSetIconUrlCreate(ctx, res.Slug).FilePathRequest(api.FilePathRequest{
			Url: i.(string),
		}).Execute()
		if err != nil {
			return httpToDiag(d, hr, err)
		}
	}
	return resourceApplicationRead(ctx, d, m)
}

func resourceApplicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.CoreApi.CoreApplicationsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	d.Set("uuid", res.Pk)
	d.Set("name", res.Name)
	d.Set("group", res.Group)
	d.Set("slug", res.Slug)
	d.Set("open_in_new_tab", res.OpenInNewTab)
	d.Set("protocol_provider", 0)
	if prov := res.Provider.Get(); prov != nil {
		d.Set("protocol_provider", int(*prov))
	}
	d.Set("meta_launch_url", res.MetaLaunchUrl)
	if res.MetaIcon.IsSet() {
		d.Set("meta_icon", res.MetaIcon.Get())
	}
	d.Set("meta_description", res.MetaDescription)
	d.Set("meta_publisher", res.MetaPublisher)
	d.Set("policy_engine_mode", res.PolicyEngineMode)
	return diags
}

func resourceApplicationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceApplicationSchemaToModel(d)

	res, hr, err := c.client.CoreApi.CoreApplicationsUpdate(ctx, d.Id()).ApplicationRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	if i, iok := d.GetOk("meta_icon"); iok {
		hr, err := c.client.CoreApi.CoreApplicationsSetIconUrlCreate(ctx, res.Slug).FilePathRequest(api.FilePathRequest{
			Url: i.(string),
		}).Execute()
		if err != nil {
			return httpToDiag(d, hr, err)
		}
	}

	d.SetId(res.Slug)
	return resourceApplicationRead(ctx, d, m)
}

func resourceApplicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.CoreApi.CoreApplicationsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
