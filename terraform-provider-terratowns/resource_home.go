package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type homeResource struct {
	client *homesClient
}

type homeModel struct {
	ID             types.String `tfsdk:"id"`
	UserUUID       types.String `tfsdk:"user_uuid"`
	Town           types.String `tfsdk:"town"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	DomainName     types.String `tfsdk:"domain_name"`
	ContentVersion types.Int64  `tfsdk:"content_version"`
}

type createHomeRequest struct {
	Town           string `json:"town"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	DomainName     string `json:"domain_name"`
	ContentVersion int64  `json:"content_version"`
}

type createHomeResponse struct {
	UUID string `json:"uuid"`
}

type readHomeResponse struct {
	UUID           string `json:"uuid"`
	Town           string `json:"town"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	DomainName     string `json:"domain_name"`
	ContentVersion int64  `json:"content_version"`
}

func NewHomeResource() resource.Resource {
	return &homeResource{}
}

func (r *homeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_home"
}

func (r *homeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"user_uuid": schema.StringAttribute{
				Required: true,
			},
			"town": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Required: true,
			},
			"domain_name": schema.StringAttribute{
				Required: true,
			},
			"content_version": schema.Int64Attribute{
				Required: true,
			},
		},
	}
}

func (r *homeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*homesClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *homesClient, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *homeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan homeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected provider configuration to initialize API client")
		return
	}

	createPayload := createHomeRequest{
		Town:           plan.Town.ValueString(),
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		DomainName:     plan.DomainName.ValueString(),
		ContentVersion: plan.ContentVersion.ValueInt64(),
	}

	body, err := json.Marshal(createPayload)
	if err != nil {
		resp.Diagnostics.AddError("Failed to serialize create payload", err.Error())
		return
	}

	url := fmt.Sprintf("%s/api/u/%s/homes", strings.TrimRight(r.client.endpoint, "/"), plan.UserUUID.ValueString())
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		resp.Diagnostics.AddError("Failed to build create request", err.Error())
		return
	}

	addCommonHeaders(httpReq, r.client.token)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create home", err.Error())
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		resp.Diagnostics.AddError("Create failed", fmt.Sprintf("status=%d body=%s", httpResp.StatusCode, string(bodyBytes)))
		return
	}

	var createResp createHomeResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&createResp); err != nil {
		resp.Diagnostics.AddError("Failed to parse create response", err.Error())
		return
	}

	plan.ID = types.StringValue(createResp.UUID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *homeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state homeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected provider configuration to initialize API client")
		return
	}

	url := fmt.Sprintf("%s/api/u/%s/homes/%s", strings.TrimRight(r.client.endpoint, "/"), state.UserUUID.ValueString(), state.ID.ValueString())
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build read request", err.Error())
		return
	}

	addCommonHeaders(httpReq, r.client.token)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read home", err.Error())
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		resp.Diagnostics.AddError("Read failed", fmt.Sprintf("status=%d body=%s", httpResp.StatusCode, string(bodyBytes)))
		return
	}

	var readResp readHomeResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&readResp); err != nil {
		resp.Diagnostics.AddError("Failed to parse read response", err.Error())
		return
	}

	state.ID = types.StringValue(readResp.UUID)
	state.Town = types.StringValue(readResp.Town)
	state.Name = types.StringValue(readResp.Name)
	state.Description = types.StringValue(readResp.Description)
	state.DomainName = types.StringValue(readResp.DomainName)
	state.ContentVersion = types.Int64Value(readResp.ContentVersion)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *homeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan homeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected provider configuration to initialize API client")
		return
	}

	updatePayload := createHomeRequest{
		Town:           plan.Town.ValueString(),
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		DomainName:     plan.DomainName.ValueString(),
		ContentVersion: plan.ContentVersion.ValueInt64(),
	}

	body, err := json.Marshal(updatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Failed to serialize update payload", err.Error())
		return
	}

	url := fmt.Sprintf("%s/api/u/%s/homes/%s", strings.TrimRight(r.client.endpoint, "/"), plan.UserUUID.ValueString(), plan.ID.ValueString())
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		resp.Diagnostics.AddError("Failed to build update request", err.Error())
		return
	}

	addCommonHeaders(httpReq, r.client.token)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update home", err.Error())
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		resp.Diagnostics.AddError("Update failed", fmt.Sprintf("status=%d body=%s", httpResp.StatusCode, string(bodyBytes)))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *homeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state homeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected provider configuration to initialize API client")
		return
	}

	url := fmt.Sprintf("%s/api/u/%s/homes/%s", strings.TrimRight(r.client.endpoint, "/"), state.UserUUID.ValueString(), state.ID.ValueString())
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build delete request", err.Error())
		return
	}

	addCommonHeaders(httpReq, r.client.token)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete home", err.Error())
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusNotFound {
		return
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		resp.Diagnostics.AddError("Delete failed", fmt.Sprintf("status=%d body=%s", httpResp.StatusCode, string(bodyBytes)))
		return
	}
}

func addCommonHeaders(req *http.Request, token string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
}

var _ resource.Resource = &homeResource{}
var _ resource.ResourceWithConfigure = &homeResource{}
