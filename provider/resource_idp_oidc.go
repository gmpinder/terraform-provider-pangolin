package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gmpinder/terraform-provider-pangolin/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &idpOidcResource{}
var _ resource.ResourceWithImportState = &idpOidcResource{}

func NewIDPOidcResource() resource.Resource {
	return &idpOidcResource{}
}

type idpOidcResource struct {
	client *client.Client
}

type idpOidcResourceModel struct {
	ID             types.Int64  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	ClientID       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
	AuthURL        types.String `tfsdk:"auth_url"`
	TokenURL       types.String `tfsdk:"token_url"`
	IdentifierPath types.String `tfsdk:"identifier_path"`
	EmailPath      types.String `tfsdk:"email_path"`
	NamePath       types.String `tfsdk:"name_path"`
	Scopes         types.String `tfsdk:"scopes"`
	AutoProvision  types.Bool   `tfsdk:"auto_provision"`
	Tags           types.String `tfsdk:"tags"`
	Variant        types.String `tfsdk:"variant"`
}

func (r *idpOidcResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_idp_oidc"
}

func (r *idpOidcResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a global OIDC Identity Provider.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The ID of the IdP.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the IdP.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"client_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The client ID.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"client_secret": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "The client secret.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"auth_url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The authorization URL.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"token_url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The token URL.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"identifier_path": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The identifier path.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"email_path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The email path (nullable).",
			},
			"name_path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The name path (nullable).",
			},
			"scopes": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The scopes.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"auto_provision": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to auto-provision users.",
			},
			"tags": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Tags for the IdP.",
			},
			"variant": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("oidc"),
				MarkdownDescription: "The IdP variant (oidc, google, azure).",
				Validators: []validator.String{
					stringvalidator.OneOf("oidc", "google", "azure"),
				},
			},
		},
	}
}

func (r *idpOidcResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *idpOidcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data idpOidcResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	idp := &client.OIDCIdP{
		Name:           data.Name.ValueStringPointer(),
		ClientID:       data.ClientID.ValueStringPointer(),
		ClientSecret:   data.ClientSecret.ValueStringPointer(),
		AuthURL:        data.AuthURL.ValueStringPointer(),
		TokenURL:       data.TokenURL.ValueStringPointer(),
		IdentifierPath: data.IdentifierPath.ValueStringPointer(),
		EmailPath:      data.EmailPath.ValueStringPointer(),
		NamePath:       data.NamePath.ValueStringPointer(),
		Scopes:         data.Scopes.ValueStringPointer(),
		AutoProvision:  data.AutoProvision.ValueBoolPointer(),
		Tags:           data.Tags.ValueStringPointer(),
		Variant:        data.Variant.ValueStringPointer(),
	}

	created, err := r.client.CreateOIDCIdP(idp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating IdP", err.Error())
		return
	}

	data.ID = types.Int64PointerValue(created.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *idpOidcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data idpOidcResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	idp, err := r.client.GetOIDCIdP(data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading IdP", err.Error())
		return
	}

	data.ID = types.Int64PointerValue(idp.ID)
	data.Name = types.StringPointerValue(idp.Name)
	data.ClientID = types.StringPointerValue(idp.ClientID)
	data.ClientSecret = types.StringPointerValue(idp.ClientSecret)
	data.AuthURL = types.StringPointerValue(idp.AuthURL)
	data.TokenURL = types.StringPointerValue(idp.TokenURL)
	data.IdentifierPath = types.StringPointerValue(idp.IdentifierPath)
	data.EmailPath = types.StringPointerValue(idp.EmailPath)
	data.NamePath = types.StringPointerValue(idp.NamePath)
	data.Scopes = types.StringPointerValue(idp.Scopes)
	data.AutoProvision = types.BoolPointerValue(idp.AutoProvision)
	data.Tags = types.StringPointerValue(idp.Tags)
	data.Variant = types.StringPointerValue(idp.Variant)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *idpOidcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state idpOidcResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	idp := &client.OIDCIdP{
		Name:           data.Name.ValueStringPointer(),
		ClientID:       data.ClientID.ValueStringPointer(),
		ClientSecret:   data.ClientSecret.ValueStringPointer(),
		AuthURL:        data.AuthURL.ValueStringPointer(),
		TokenURL:       data.TokenURL.ValueStringPointer(),
		IdentifierPath: data.IdentifierPath.ValueStringPointer(),
		EmailPath:      data.EmailPath.ValueStringPointer(),
		NamePath:       data.NamePath.ValueStringPointer(),
		Scopes:         data.Scopes.ValueStringPointer(),
		AutoProvision:  data.AutoProvision.ValueBoolPointer(),
		Tags:           data.Tags.ValueStringPointer(),
		Variant:        data.Variant.ValueStringPointer(),
	}

	_, err := r.client.UpdateOIDCIdP(state.ID.ValueInt64(), idp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating IdP", err.Error())
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *idpOidcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data idpOidcResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOIDCIdP(data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting IdP", err.Error())
		return
	}
}

func (r *idpOidcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid IdP ID",
			fmt.Sprintf("Error parsing IDP ID: %s", err),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
