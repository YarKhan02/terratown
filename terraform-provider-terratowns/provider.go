package main

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type homesProvider struct{}

type homesProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

type homesClient struct {
	endpoint string
	token    string
	client   *http.Client
}

func NewProvider() provider.Provider {
	return &homesProvider{}
}

func (p *homesProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "terratowns"
}

func (p *homesProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Required: true,
			},
			"token": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *homesProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config homesProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Endpoint.IsUnknown() || config.Token.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown Provider Configuration",
			"Cannot create API client with unknown endpoint or token.",
		)
		return
	}

	resp.ResourceData = &homesClient{
		endpoint: config.Endpoint.ValueString(),
		token:    config.Token.ValueString(),
		client:   &http.Client{},
	}
}

func (p *homesProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewHomeResource,
	}
}

func (p *homesProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}
