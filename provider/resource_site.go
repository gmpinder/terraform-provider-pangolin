package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gmpinder/terraform-provider-pangolin/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &siteResource{}
var _ resource.ResourceWithImportState = &siteResource{}

func NewSiteResource() resource.Resource {
	return &siteResource{}
}

type siteResource struct {
	client *client.Client
}

type siteResourceModel struct {
	ID         types.Int64  `tfsdk:"id"`
	OrgID      types.String `tfsdk:"org_id"`
	Name       types.String `tfsdk:"name"`
	NiceId     types.String `tfsdk:"nice_id"`
	NewtID     types.String `tfsdk:"newt_id"`
	PubKey     types.String `tfsdk:"pub_key"`
	NewtSecret types.String `tfsdk:"newt_secret"`
	Address    types.String `tfsdk:"address"`
	Subnet     types.String `tfsdk:"subnet"`
	Type       types.String `tfsdk:"type"`
}

func (r *siteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (r *siteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an organization role.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The ID of the site.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},

			// Inputs
			"org_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The org ID the site is deployed under.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the resource.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"nice_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The human readable string ID of the site.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-z0-9-]$`), "Should be lowercase with hyphens"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("newt"),
				Validators: []validator.String{
					stringvalidator.OneOf("newt", "wireguard", "local"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			// read-only
			"newt_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID to use for setting up a Newt instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pub_key": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The public key used by the Newt Instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"newt_secret": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "The secret key used by the Newt Instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"address": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The address CIDR of the Newt Instance",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subnet": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The subnet address CIDR of the Newt Instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *siteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *siteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data siteResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	site := &client.Site{
		Name: data.Name.ValueStringPointer(),
		Type: data.Type.ValueStringPointer(),
	}

	if data.Type.ValueString() == "newt" {
		defaults, err := r.client.GetSiteDefaults(data.OrgID.ValueString())

		if err != nil {
			resp.Diagnostics.AddError("Error retrieving site defaults", err.Error())
			return
		}

		site.PubKey = &defaults.PublicKey
		site.Subnet = &defaults.Subnet
		site.Address = &defaults.ClientAddress
		site.NewtID = &defaults.NewtId
		site.Secret = &defaults.NewtSecret
	}

	created, err := r.client.CreateSite(data.OrgID.ValueString(), site)
	if err != nil {
		resp.Diagnostics.AddError("Error creating site", err.Error())
		return
	}

	data.ID = types.Int64PointerValue(created.ID)
	data.Address = types.StringPointerValue(created.Address)
	data.Subnet = types.StringPointerValue(created.Subnet)
	data.NewtID = types.StringPointerValue(site.NewtID)
	data.NewtSecret = types.StringPointerValue(site.Secret)
	data.PubKey = types.StringPointerValue(site.PubKey)

	if !data.NiceId.IsUnknown() {
		su := &client.Site{
			Name:   data.Name.ValueStringPointer(),
			NiceId: data.NiceId.ValueStringPointer(),
		}

		updated, err := r.client.UpdateSite(*created.ID, su)

		if err != nil {
			resp.Diagnostics.AddError("Error updating `nice_id`", err.Error())
		} else {
			data.NiceId = types.StringPointerValue(updated.NiceId)
		}
	} else {
		data.NiceId = types.StringPointerValue(created.NiceId)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *siteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data siteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	site, err := r.client.GetSite(data.OrgID.ValueString(), data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading role", err.Error())
		return
	}

	data.Name = types.StringPointerValue(site.Name)
	data.ID = types.Int64PointerValue(site.ID)
	data.Address = types.StringPointerValue(site.Address)
	data.NewtID = types.StringPointerValue(site.NewtID)
	data.NiceId = types.StringPointerValue(site.NiceId)
	data.OrgID = types.StringPointerValue(site.OrgID)
	data.PubKey = types.StringPointerValue(site.PubKey)
	data.Subnet = types.StringPointerValue(site.Subnet)
	data.Type = types.StringPointerValue(site.Type)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *siteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state siteResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	site := &client.Site{
		Name: data.Name.ValueStringPointer(),
	}

	if !data.NiceId.IsNull() {
		site.NiceId = data.NiceId.ValueStringPointer()
	}

	_, err := r.client.UpdateSite(state.ID.ValueInt64(), site)
	if err != nil {
		resp.Diagnostics.AddError("Error updating site", err.Error())
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *siteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data siteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSite(data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting site", err.Error())
		return
	}
}

func (r *siteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: org_id/role_id
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: org_id/role_id. Got: %q", req.ID),
		)
		return
	}

	siteID, err := strconv.ParseInt(idParts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected role_id to be an integer. Got: %q", idParts[1]),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("org_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), siteID)...)
}
