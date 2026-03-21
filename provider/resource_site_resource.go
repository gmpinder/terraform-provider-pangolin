package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gmpinder/terraform-provider-pangolin/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	// "github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &siteResourceResource{}
var _ resource.ResourceWithImportState = &siteResourceResource{}

func NewSiteResourceResource() resource.Resource {
	return &siteResourceResource{}
}

type siteResourceResource struct {
	client *client.Client
}

type siteResourceResourceModel struct {
	ID                 types.Int64  `tfsdk:"id"`
	NiceID             types.String `tfsdk:"nice_id"`
	OrgID              types.String `tfsdk:"org_id"`
	Name               types.String `tfsdk:"name"`
	Mode               types.String `tfsdk:"mode"`
	SiteID             types.Int64  `tfsdk:"site_id"`
	Destination        types.String `tfsdk:"destination"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	Alias              types.String `tfsdk:"alias"`
	UserIDs            types.List   `tfsdk:"user_ids"`
	RoleIDs            types.List   `tfsdk:"role_ids"`
	ClientIDs          types.List   `tfsdk:"client_ids"`
	TCPPortRangeString types.String `tfsdk:"tcp_port_range_string"`
	UDPPortRangeString types.String `tfsdk:"udp_port_range_string"`
	DisableIcmp        types.Bool   `tfsdk:"disable_icmp"`
}

func (r *siteResourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site_resource"
}

func (r *siteResourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a site resource (Host or CIDR mode).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The ID of the site resource.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"nice_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The human-readable ID of the site resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the organization.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the site resource.",
			},
			"mode": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The mode of the resource (host or cidr).",
				Validators: []validator.String{
					stringvalidator.OneOf("host", "cidr"),
				},
			},
			"site_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the site.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"destination": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The destination address or CIDR.",
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether the resource is enabled.",
			},
			"alias": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The alias for the resource.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(?:[a-zA-Z0-9*?](?:[a-zA-Z0-9*?-]{0,61}[a-zA-Z0-9*?])?\.)+[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`),
						"Alias must be a fully qualified domain name with optional wildcards",
					),
				},
			},
			"user_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				// Default: listdefault.StaticValue(
				// 	types.ListValueMust(types.StringType, make([]attr.Value, 0)),
				// ),
				MarkdownDescription: "The list of user IDs allowed to access this resource.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"role_ids": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
				Computed:    true,
				// Default: listdefault.StaticValue(
				// 	types.ListValueMust(types.Int64Type, make([]attr.Value, 0)),
				// ),
				MarkdownDescription: "The list of role IDs allowed to access this resource.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"client_ids": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
				Computed:    true,
				// Default: listdefault.StaticValue(
				// 	types.ListValueMust(types.Int64Type, make([]attr.Value, 0)),
				// ),
				MarkdownDescription: "The list of client IDs allowed to access this resource.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"tcp_port_range_string": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The TCP port range allowed (e.g., '80,443' or '*'). Defaults to blocking traffic.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(?:(?:[0-9-]+,?)+|\*)$`),
						"Port range must be like 80,43,8000-8500",
					),
				},
			},
			"udp_port_range_string": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The UDP port range allowed (e.g., '53' or '*'). Defaults to blocking traffic.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(?:(?:[0-9-]+,?)+|\*)$`),
						"Port range must be like 80,43,8000-8500",
					),
				},
			},
			"disable_icmp": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to disable ICMP for this resource.",
			},
		},
	}
}

func (r *siteResourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}

	r.client = c
}

func (r *siteResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data siteResourceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res := &client.SiteResource{
		Name:               data.Name.ValueStringPointer(),
		Mode:               data.Mode.ValueStringPointer(),
		SiteID:             data.SiteID.ValueInt64Pointer(),
		Destination:        data.Destination.ValueStringPointer(),
		Enabled:            data.Enabled.ValueBoolPointer(),
		TCPPortRangeString: data.TCPPortRangeString.ValueStringPointer(),
		UDPPortRangeString: data.UDPPortRangeString.ValueStringPointer(),
		DisableIcmp:        data.DisableIcmp.ValueBoolPointer(),
		Alias:              data.Alias.ValueStringPointer(),
		UserIDs:            make([]string, 0),
		RoleIDs:            make([]int64, 0),
		ClientIDs:          make([]int64, 0),
	}

	if !data.UserIDs.IsUnknown() {
		resp.Diagnostics.Append(data.UserIDs.ElementsAs(ctx, &res.UserIDs, false)...)
	}

	if !data.RoleIDs.IsUnknown() {
		resp.Diagnostics.Append(data.RoleIDs.ElementsAs(ctx, &res.RoleIDs, false)...)
	}

	if !data.ClientIDs.IsUnknown() {
		resp.Diagnostics.Append(data.ClientIDs.ElementsAs(ctx, &res.ClientIDs, false)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.client.CreateSiteResource(data.OrgID.ValueString(), res)
	if err != nil {
		resp.Diagnostics.AddError("Error creating site resource", err.Error())
		return
	}

	data.ID = types.Int64PointerValue(created.ID)
	data.NiceID = types.StringPointerValue(created.NiceID)

	roles, diags := types.ListValueFrom(ctx, types.Int64Type, res.RoleIDs)
	resp.Diagnostics.Append(diags...)
	data.RoleIDs = roles

	users, diags := types.ListValueFrom(ctx, types.StringType, res.UserIDs)
	resp.Diagnostics.Append(diags...)
	data.UserIDs = users

	clients, diags := types.ListValueFrom(ctx, types.Int64Type, res.ClientIDs)
	resp.Diagnostics.Append(diags...)
	data.ClientIDs = clients

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *siteResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data siteResourceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetSiteResource(data.OrgID.ValueString(), data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading site resource", err.Error())
		return
	}

	data.Name = types.StringPointerValue(res.Name)
	data.SiteID = types.Int64PointerValue(res.SiteID)
	data.Mode = types.StringPointerValue(res.Mode)
	data.Destination = types.StringPointerValue(res.Destination)
	data.Enabled = types.BoolPointerValue(res.Enabled)
	data.Alias = types.StringPointerValue(res.Alias)
	data.TCPPortRangeString = types.StringPointerValue(res.TCPPortRangeString)
	data.UDPPortRangeString = types.StringPointerValue(res.UDPPortRangeString)
	data.DisableIcmp = types.BoolPointerValue(res.DisableIcmp)

	roleIDs, err := r.client.GetSiteResourceRoles(data.ID.ValueInt64())

	if err != nil {
		resp.Diagnostics.AddError("Failed to get roles for site resource", err.Error())
	} else {
		roles, diags := types.ListValueFrom(ctx, types.Int64Type, roleIDs)
		resp.Diagnostics.Append(diags...)
		data.RoleIDs = roles
	}

	userIDs, err := r.client.GetSiteResourceUsers(data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get roles for site resource", err.Error())
	} else {

		users, diags := types.ListValueFrom(ctx, types.StringType, userIDs)
		resp.Diagnostics.Append(diags...)
		data.UserIDs = users
	}

	clientIDs, err := r.client.GetSiteResourceClients(data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get roles for site resource", err.Error())
	} else {
		clients, diags := types.ListValueFrom(ctx, types.Int64Type, clientIDs)
		resp.Diagnostics.Append(diags...)
		data.ClientIDs = clients
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *siteResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state siteResourceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res := &client.SiteResource{
		Name:               data.Name.ValueStringPointer(),
		Mode:               data.Mode.ValueStringPointer(),
		SiteID:             data.SiteID.ValueInt64Pointer(),
		Destination:        data.Destination.ValueStringPointer(),
		Enabled:            data.Enabled.ValueBoolPointer(),
		TCPPortRangeString: data.TCPPortRangeString.ValueStringPointer(),
		UDPPortRangeString: data.UDPPortRangeString.ValueStringPointer(),
		DisableIcmp:        data.DisableIcmp.ValueBoolPointer(),
	}

	if !data.Alias.IsNull() {
		s := data.Alias.ValueString()
		res.Alias = &s
	}

	resp.Diagnostics.Append(data.UserIDs.ElementsAs(ctx, &res.UserIDs, false)...)
	resp.Diagnostics.Append(data.RoleIDs.ElementsAs(ctx, &res.RoleIDs, false)...)
	resp.Diagnostics.Append(data.ClientIDs.ElementsAs(ctx, &res.ClientIDs, false)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateSiteResource(state.ID.ValueInt64(), res)
	if err != nil {
		resp.Diagnostics.AddError("Error updating site resource", err.Error())
		return
	}

	data.ID = state.ID
	data.NiceID = state.NiceID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *siteResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data siteResourceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSiteResource(data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting site resource", err.Error())
		return
	}
}

func (r *siteResourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: org_id/id
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: org_id/id. Got: %q", req.ID),
		)
		return
	}

	resID, err := strconv.ParseInt(idParts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected id to be an integer. Got: %q", idParts[1]),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("org_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), resID)...)
}
