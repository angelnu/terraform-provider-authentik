package provider

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSourcePlex() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSourcePlexCreate,
		ReadContext:   resourceSourcePlexRead,
		UpdateContext: resourceSourcePlexUpdate,
		DeleteContext: resourceSourcePlexDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
			"authentication_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enrollment_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"policy_engine_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.POLICYENGINEMODE_ANY,
			},

			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"allowed_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"allow_friends": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"plex_token": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceSourcePlexSchemaToSource(d *schema.ResourceData) (*api.PlexSourceRequest, diag.Diagnostics) {
	r := api.PlexSourceRequest{
		Name:    d.Get("name").(string),
		Slug:    d.Get("slug").(string),
		Enabled: boolToPointer(d.Get("enabled").(bool)),

		ClientId:     stringToPointer(d.Get("client_id").(string)),
		AllowFriends: boolToPointer(d.Get("allow_friends").(bool)),
		PlexToken:    d.Get("plex_token").(string),
	}

	r.AuthenticationFlow.Set(stringToPointer(d.Get("authentication_flow").(string)))
	r.EnrollmentFlow.Set(stringToPointer(d.Get("enrollment_flow").(string)))

	pm := api.PolicyEngineMode(d.Get("policy_engine_mode").(string))
	r.PolicyEngineMode = &pm

	allowedServers := d.Get("allowed_servers").([]interface{})
	as := make([]string, len(allowedServers))
	for i, prov := range allowedServers {
		as[i] = prov.(string)
	}
	r.AllowedServers = &as

	return &r, nil
}

func resourceSourcePlexCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceSourcePlexSchemaToSource(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.SourcesApi.SourcesPlexCreate(ctx).PlexSourceRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourcePlexRead(ctx, d, m)
}

func resourceSourcePlexRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	res, hr, err := c.client.SourcesApi.SourcesPlexRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("name", res.Name)
	d.Set("slug", res.Slug)
	d.Set("uuid", res.Pk)

	if res.AuthenticationFlow.IsSet() {
		d.Set("authentication_flow", res.AuthenticationFlow.Get())
	}
	if res.EnrollmentFlow.IsSet() {
		d.Set("enrollment_flow", res.EnrollmentFlow.Get())
	}
	d.Set("enabled", res.Enabled)
	d.Set("policy_engine_mode", res.PolicyEngineMode)

	d.Set("client_id", res.ClientId)
	d.Set("allowed_servers", res.AllowedServers)
	d.Set("allow_friends", res.AllowFriends)
	d.Set("plex_token", res.PlexToken)
	return diags
}

func resourceSourcePlexUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	app, di := resourceSourcePlexSchemaToSource(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.SourcesApi.SourcesPlexUpdate(ctx, d.Id()).PlexSourceRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourcePlexRead(ctx, d, m)
}

func resourceSourcePlexDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.SourcesApi.SourcesPlexDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
