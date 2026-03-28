package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/gmpinder/terraform-provider-pangolin/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &resourceResource{}
var _ resource.ResourceWithImportState = &resourceResource{}

func NewResourceResource() resource.Resource {
	return &resourceResource{}
}

type resourceResource struct {
	client *client.Client
}

type resourceResourceModel struct {
	ID                    types.Int64  `tfsdk:"id"`
	OrgID                 types.String `tfsdk:"org_id"`
	Name                  types.String `tfsdk:"name"`
	Protocol              types.String `tfsdk:"protocol"`
	Http                  types.Bool   `tfsdk:"http"`
	Subdomain             types.String `tfsdk:"subdomain"`
	DomainID              types.String `tfsdk:"domain_id"`
	EmailWhitelistEnabled types.Bool   `tfsdk:"email_whitelist_enabled"`
	ApplyRules            types.Bool   `tfsdk:"apply_rules"`
	NiceId                types.String `tfsdk:"nice_id"`
	Ssl                   types.Bool   `tfsdk:"ssl"`
	BlockAccess           types.Bool   `tfsdk:"block_access"`
	Sso                   types.Bool   `tfsdk:"sso"`
	ProxyPort             types.Int32  `tfsdk:"proxy_port"`
	Enabled               types.Bool   `tfsdk:"enabled"`
	StickySession         types.Bool   `tfsdk:"sticky_session"`
	TlsServerName         types.String `tfsdk:"tls_server_name"`
	SetHostHeader         types.String `tfsdk:"host_header"`
	Headers               types.List   `tfsdk:"headers"`
	ProxyProtocol         types.Bool   `tfsdk:"proxy_protocol"`
	ProxyProtocolVersion  types.Int32  `tfsdk:"proxy_protocol_version"`
	PostAuthPath          types.String `tfsdk:"post_auth_path"`
}

type resourceHeader struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func (r *resourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource"
}

func (r *resourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	proxyExpressions := path.Expressions{
		path.MatchRoot("proxy_port"),
		path.MatchRoot("proxy_protocol"),
		path.MatchRoot("proxy_protocol_version"),
	}
	httpExpressions := path.Expressions{
		path.MatchRoot("domain_id"),
		path.MatchRoot("subdomain"),
		path.MatchRoot("sticky_session"),
		path.MatchRoot("post_auth_path"),
		path.MatchRoot("sso"),
		path.MatchRoot("tls_server_name"),
		path.MatchRoot("headers"),
		path.MatchRoot("host_header"),
		path.MatchRoot("email_whitelist_enabled"),
		path.MatchRoot("block_access"),
	}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an app-style resource (HTTP/TCP/UDP).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The ID of the resource.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
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
				MarkdownDescription: "The name of the resource.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"nice_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The unique, human-readable ID of the resource.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
					stringvalidator.RegexMatches(
						regexp.MustCompile("^[a-zA-Z0-9-]+$"),
						"String must consist of lower and uppercase letters, numbers, and the '-' character.",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"protocol": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The protocol of the resource (tcp or udp).",
				Validators: []validator.String{
					stringvalidator.OneOf("tcp", "udp"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable the resource.",
			},
			"http": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether the resource is an HTTP resource.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},

			// `http` true
			"subdomain": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The subdomain for the resource. Requires `http` to be true.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^([a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?\.)*[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`),
						"Must be a valid subdomain",
					),
					stringvalidator.ConflictsWith(),
				},
			},
			"domain_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The ID of the domain. Requires `http` to be true.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(proxyExpressions...),
				},
			},
			"sticky_session": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether or not to enable sticky sessions. Requires `http` to be true.",
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(proxyExpressions...),
				},
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"post_auth_path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The HTTP path of the resource to shared. Requires `http` to be true.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(proxyExpressions...),
				},
			},
			"ssl": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Enable SSL for the resource.",
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(proxyExpressions...),
				},
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"block_access": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Block access to the resource.",
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(proxyExpressions...),
				},
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sso": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Allow Pangolin SSO auth.",
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(proxyExpressions...),
				},
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"tls_server_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The hostname expected by the SSL endpoint",
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("ssl"),
					}...),
					stringvalidator.ConflictsWith(proxyExpressions...),
				},
			},
			"host_header": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Set a custom host header to set for requests.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(proxyExpressions...),
				},
			},
			"headers": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Header name.",
						},
						"value": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Header value.",
						},
					},
				},
				Optional:            true,
				MarkdownDescription: "List of headers to set when forwarding requests.",
				Validators: []validator.List{
					listvalidator.ConflictsWith(proxyExpressions...),
				},
			},
			"email_whitelist_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Allow whitelisting access based on email address. Requires SMTP to be setup.",
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(proxyExpressions...),
				},
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"apply_rules": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Apply block list rules.",
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(proxyExpressions...),
				},
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},

			// `http` false
			"proxy_port": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "The port to proxy if `http` is false",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
					int32validator.ConflictsWith(),
				},
			},
			"proxy_protocol": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Enable the proxy protocol.",
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(httpExpressions...),
				},
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"proxy_protocol_version": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Version 1 is text-based and widely supported. Version 2 is binary and more efficient but less compatible. Make sure servers transport is added to dynamic config.",
				Validators: []validator.Int32{
					int32validator.OneOf(1, 2),
					int32validator.ConflictsWith(httpExpressions...),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseNonNullStateForUnknown(),
				},
			},
		},
	}
}

func (r *resourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Http.ValueBool() && data.DomainID.IsNull() {
		resp.Diagnostics.AddError("`domain_id` is null", "`domain_id` must not be null if `http` is true.")
	}

	if !data.Http.ValueBool() && data.ProxyPort.IsNull() {
		resp.Diagnostics.AddError("`proxy_port` is null", "`proxy_port` must not be null if `http` is false.")
	}

	var headers []client.ResourceHeader

	if !data.Headers.IsNull() {
		h := make([]resourceHeader, len(data.Headers.Elements()))
		resp.Diagnostics.Append(data.Headers.ElementsAs(ctx, &h, false)...)

		headers = make([]client.ResourceHeader, len(h))

		for i, header := range h {
			headers[i] = client.ResourceHeader{
				Name:  header.Name.ValueString(),
				Value: header.Value.ValueString(),
			}
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	createResource := &client.Resource{
		Name:     data.Name.ValueStringPointer(),
		Protocol: data.Protocol.ValueStringPointer(),
		Http:     data.Http.ValueBoolPointer(),
	}

	if *createResource.Http {
		createResource.DomainID = data.DomainID.ValueStringPointer()
		createResource.Subdomain = data.Subdomain.ValueStringPointer()
		createResource.StickySession = data.StickySession.ValueBoolPointer()
		createResource.PostAuthPath = data.PostAuthPath.ValueStringPointer()
	} else {
		createResource.ProxyPort = data.ProxyPort.ValueInt32Pointer()
	}

	created, err := r.client.CreateResource(data.OrgID.ValueString(), createResource)
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource", err.Error())
		return
	}

	if data.NiceId.IsUnknown() {
		data.NiceId = types.StringPointerValue(created.NiceID)
	}

	if data.StickySession.IsUnknown() {
		data.StickySession = types.BoolPointerValue(created.StickySession)
	}

	if data.Ssl.IsUnknown() {
		data.Ssl = types.BoolPointerValue(created.Ssl)
	}

	if data.BlockAccess.IsUnknown() {
		data.BlockAccess = types.BoolPointerValue(created.BlockAccess)
	}

	if data.Sso.IsUnknown() {
		data.Sso = types.BoolPointerValue(created.Sso)
	}

	if data.EmailWhitelistEnabled.IsUnknown() {
		data.EmailWhitelistEnabled = types.BoolPointerValue(created.EmailWhitelistEnabled)
	}

	if data.ApplyRules.IsUnknown() {
		data.ApplyRules = types.BoolPointerValue(created.ApplyRules)
	}

	if data.ProxyProtocol.IsUnknown() {
		data.ProxyProtocol = types.BoolPointerValue(created.ProxyProtocol)
	}

	if data.ProxyProtocolVersion.IsUnknown() {
		data.ProxyProtocolVersion = types.Int32PointerValue(created.ProxyProtocolVersion)
	}

	updateResource := &client.Resource{
		Enabled: data.Enabled.ValueBoolPointer(),
		NiceID:  data.NiceId.ValueStringPointer(),
	}
	if *created.Http {
		updateResource.EmailWhitelistEnabled = data.EmailWhitelistEnabled.ValueBoolPointer()
		updateResource.ApplyRules = data.ApplyRules.ValueBoolPointer()
		updateResource.Ssl = data.Ssl.ValueBoolPointer()
		updateResource.BlockAccess = data.BlockAccess.ValueBoolPointer()
		updateResource.Sso = data.Sso.ValueBoolPointer()
		updateResource.TlsServerName = data.TlsServerName.ValueStringPointer()
		updateResource.SetHostHeader = data.SetHostHeader.ValueStringPointer()
		updateResource.Headers = headers
	} else {
		updateResource.ProxyProtocol = data.ProxyProtocol.ValueBoolPointer()
		updateResource.ProxyProtocolVersion = data.ProxyProtocolVersion.ValueInt32Pointer()
	}

	updated, err := r.client.UpdateResource(*created.ID, updateResource)

	if err != nil {
		resp.Diagnostics.AddError("Error adding extra properties to resource during creation", err.Error())
	} else {
		data.NiceId = types.StringPointerValue(updated.NiceID)
		data.EmailWhitelistEnabled = types.BoolPointerValue(updated.EmailWhitelistEnabled)
		data.ApplyRules = types.BoolPointerValue(updated.ApplyRules)
		data.StickySession = types.BoolPointerValue(updated.StickySession)
		data.Ssl = types.BoolPointerValue(updated.Ssl)
		data.BlockAccess = types.BoolPointerValue(updated.BlockAccess)
		data.Sso = types.BoolPointerValue(updated.Sso)
		data.ProxyProtocol = types.BoolPointerValue(updated.ProxyProtocol)
	}

	data.ID = types.Int64PointerValue(created.ID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetResource(data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading resource", err.Error())
		return
	}

	data.ApplyRules = types.BoolPointerValue(res.ApplyRules)
	data.BlockAccess = types.BoolPointerValue(res.BlockAccess)
	data.DomainID = types.StringPointerValue(res.DomainID)
	data.EmailWhitelistEnabled = types.BoolPointerValue(res.EmailWhitelistEnabled)
	data.Enabled = types.BoolPointerValue(res.Enabled)
	data.Name = types.StringPointerValue(res.Name)
	data.NiceId = types.StringPointerValue(res.NiceID)
	data.PostAuthPath = types.StringPointerValue(res.PostAuthPath)
	data.Protocol = types.StringPointerValue(res.Protocol)
	data.ID = types.Int64PointerValue(res.ID)
	data.OrgID = types.StringPointerValue(res.OrgID)
	data.ProxyPort = types.Int32PointerValue(res.ProxyPort)
	data.ProxyProtocol = types.BoolPointerValue(res.ProxyProtocol)
	data.ProxyProtocolVersion = types.Int32PointerValue(res.ProxyProtocolVersion)
	data.SetHostHeader = types.StringPointerValue(res.SetHostHeader)
	data.Ssl = types.BoolPointerValue(res.Ssl)
	data.Sso = types.BoolPointerValue(res.Sso)
	data.StickySession = types.BoolPointerValue(res.StickySession)
	data.Subdomain = types.StringPointerValue(res.Subdomain)
	data.TlsServerName = types.StringPointerValue(res.TlsServerName)

	if len(res.Headers) > 0 {
		rHeaders := make([]resourceHeader, len(res.Headers))

		for i, header := range res.Headers {
			rHeaders[i] = resourceHeader{
				Name:  types.StringValue(header.Name),
				Value: types.StringValue(header.Value),
			}
		}
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("headers"), rHeaders)...)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state resourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	var headers []client.ResourceHeader

	if !data.Headers.IsNull() {
		h := make([]resourceHeader, len(data.Headers.Elements()))
		resp.Diagnostics.Append(data.Headers.ElementsAs(ctx, &h, false)...)

		headers = make([]client.ResourceHeader, len(h))

		for i, header := range h {
			headers[i] = client.ResourceHeader{
				Name:  header.Name.ValueString(),
				Value: header.Value.ValueString(),
			}
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	res := &client.Resource{
		Name:    data.Name.ValueStringPointer(),
		NiceID:  data.NiceId.ValueStringPointer(),
		Enabled: data.Enabled.ValueBoolPointer(),
	}

	if data.Http.ValueBool() {
		res.Subdomain = data.Subdomain.ValueStringPointer()
		res.DomainID = data.DomainID.ValueStringPointer()
		res.EmailWhitelistEnabled = data.EmailWhitelistEnabled.ValueBoolPointer()
		res.ApplyRules = data.ApplyRules.ValueBoolPointer()
		res.Ssl = data.Ssl.ValueBoolPointer()
		res.BlockAccess = data.BlockAccess.ValueBoolPointer()
		res.Sso = data.Sso.ValueBoolPointer()
		res.StickySession = data.StickySession.ValueBoolPointer()
		res.TlsServerName = data.TlsServerName.ValueStringPointer()
		res.SetHostHeader = data.SetHostHeader.ValueStringPointer()
		res.PostAuthPath = data.PostAuthPath.ValueStringPointer()
		res.Headers = headers
	} else {
		res.ProxyPort = data.ProxyPort.ValueInt32Pointer()
		res.ProxyProtocol = data.ProxyProtocol.ValueBoolPointer()
		res.ProxyProtocolVersion = data.ProxyProtocolVersion.ValueInt32Pointer()
	}

	_, err := r.client.UpdateResource(state.ID.ValueInt64(), res)
	if err != nil {
		resp.Diagnostics.AddError("Error updating resource", err.Error())
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteResource(data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting resource", err.Error())
		return
	}
}

func (r *resourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected id to be an integer. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), resID)...)
}
