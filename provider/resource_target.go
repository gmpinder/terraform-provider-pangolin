package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gmpinder/terraform-provider-pangolin/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ resource.Resource = &targetResource{}
var _ resource.ResourceWithImportState = &targetResource{}

func NewTargetResource() resource.Resource {
	return &targetResource{}
}

type targetResource struct {
	client *client.Client
}

type targetResourceModel struct {
	ID              types.Int64  `tfsdk:"id"`
	ResourceID      types.Int64  `tfsdk:"resource_id"`
	SiteID          types.Int64  `tfsdk:"site_id"`
	IP              types.String `tfsdk:"ip"`
	Port            types.Int32  `tfsdk:"port"`
	Method          types.String `tfsdk:"method"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	Path            types.String `tfsdk:"path"`
	PathMatchType   types.String `tfsdk:"path_match_type"`
	RewritePath     types.String `tfsdk:"rewrite_path"`
	RewritePathType types.String `tfsdk:"rewrite_path_type"`
	Priority        types.Int32  `tfsdk:"priority"`
	HealthCheck     types.Object `tfsdk:"health_check"`
}

type targetHealthCheck struct {
	Enabled           types.Bool   `tfsdk:"enabled"`
	Path              types.String `tfsdk:"path"`
	Scheme            types.String `tfsdk:"scheme"`
	Mode              types.String `tfsdk:"mode"`
	Hostname          types.String `tfsdk:"hostname"`
	Port              types.Int32  `tfsdk:"port"`
	Interval          types.Int32  `tfsdk:"interval"`
	UnhealthyInterval types.Int32  `tfsdk:"unhealthy_interval"`
	Timeout           types.Int32  `tfsdk:"timeout"`
	Headers           types.List   `tfsdk:"headers"`
	FollowRedirects   types.Bool   `tfsdk:"follow_redirects"`
	Method            types.String `tfsdk:"method"`
	Status            types.Int32  `tfsdk:"status"`
	TlsServerName     types.String `tfsdk:"tls_server_name"`
}

type targetHCHeader struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func (r *targetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_target"
}

func (r *targetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a backend target for a resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The ID of the target.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"resource_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the resource this target belongs to.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"site_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the site.",
			},
			"ip": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The IP address of the target.",
			},
			"port": schema.Int32Attribute{
				Required:            true,
				MarkdownDescription: "The port of the target.",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
			},
			"method": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("http"),
				MarkdownDescription: "The load balancing method.",
				Validators: []validator.String{
					stringvalidator.OneOf("http", "https", "h2c"),
				},
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether the target is enabled.",
			},
			"path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The path for the target.",
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("path_match_type")),
				},
			},
			"path_match_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The path match type.",
				Validators: []validator.String{
					stringvalidator.OneOf("prefix", "exact", "regex"),
					stringvalidator.AlsoRequires(path.MatchRoot("path")),
				},
			},
			"rewrite_path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The rewrite path.",
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("rewrite_path_type")),
				},
			},
			"rewrite_path_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The rewrite path type.",
				Validators: []validator.String{
					stringvalidator.OneOf("prefix", "exact", "regex", "stripPrefix"),
					stringvalidator.AlsoRequires(path.MatchRoot("rewrite_path")),
				},
			},
			"priority": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The priority of the target.",
				Default:             int32default.StaticInt32(100),
				Validators: []validator.Int32{
					int32validator.Between(1, 100),
				},
			},
			"health_check": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "Whether health checks are enabled.",
					},
					"path": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
						MarkdownDescription: "The health check path.",
					},
					"scheme": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("http"),
						MarkdownDescription: "The health check scheme (http or https).",
						Validators: []validator.String{
							stringvalidator.OneOf("http", "https"),
						},
					},
					"mode": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("http"),
						MarkdownDescription: "The health check mode.",
						Validators: []validator.String{
							stringvalidator.OneOf("http"),
						},
					},
					"hostname": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The health check hostname or IP.",
					},
					"port": schema.Int32Attribute{
						Required:            true,
						MarkdownDescription: "The health check port.",
						Validators: []validator.Int32{
							int32validator.Between(1, 65535),
						},
					},
					"interval": schema.Int32Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The health check interval.",
						Default:             int32default.StaticInt32(5),
						Validators: []validator.Int32{
							int32validator.AtLeast(5),
						},
					},
					"unhealthy_interval": schema.Int32Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The health check unhealthy interval.",
						Default:             int32default.StaticInt32(5),
						Validators: []validator.Int32{
							int32validator.AtLeast(5),
						},
					},
					"timeout": schema.Int32Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The health check timeout.",
						Default:             int32default.StaticInt32(1),
						Validators: []validator.Int32{
							int32validator.AtLeast(1),
						},
					},
					"follow_redirects": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
						MarkdownDescription: "Whether to follow redirects during health checks.",
					},
					"method": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The health check method.",
						Default:             stringdefault.StaticString("GET"),
						Validators: []validator.String{
							stringvalidator.OneOf("GET", "PUT", "POST", "HEAD", "DELETE"),
						},
					},
					"status": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "The expected health check status code.",
						Validators: []validator.Int32{
							int32validator.AtLeast(1),
						},
					},
					"tls_server_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The TLS server name for health checks.",
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
					},
				},
				Optional:            true,
				MarkdownDescription: "Health Check options",
			},
		},
	}
}

func (r *targetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *targetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data targetResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var healthcheck *targetHealthCheck
	var hc_headers []client.TargetHeader
	resp.Diagnostics.Append(data.HealthCheck.As(ctx, &healthcheck, basetypes.ObjectAsOptions{})...)

	if healthcheck != nil && !healthcheck.Headers.IsNull() {
		h := make([]targetHCHeader, len(healthcheck.Headers.Elements()))
		resp.Diagnostics.Append(healthcheck.Headers.ElementsAs(ctx, &h, false)...)

		hc_headers = make([]client.TargetHeader, len(h))

		for i, header := range h {
			hc_headers[i] = client.TargetHeader{
				Name:  header.Name.ValueString(),
				Value: header.Value.ValueString(),
			}
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	target := &client.Target{
		SiteID:          data.SiteID.ValueInt64Pointer(),
		IP:              data.IP.ValueStringPointer(),
		Port:            data.Port.ValueInt32Pointer(),
		Method:          data.Method.ValueStringPointer(),
		Enabled:         data.Enabled.ValueBoolPointer(),
		Path:            data.Path.ValueStringPointer(),
		PathMatchType:   data.PathMatchType.ValueStringPointer(),
		RewritePath:     data.RewritePath.ValueStringPointer(),
		RewritePathType: data.RewritePathType.ValueStringPointer(),
		Priority:        data.Priority.ValueInt32Pointer(),
	}

	if healthcheck != nil {
		target.HCEnabled = healthcheck.Enabled.ValueBoolPointer()
		target.HCPath = healthcheck.Path.ValueStringPointer()
		target.HCScheme = healthcheck.Scheme.ValueStringPointer()
		target.HCMode = healthcheck.Mode.ValueStringPointer()
		target.HCHostname = healthcheck.Hostname.ValueStringPointer()
		target.HCPort = healthcheck.Port.ValueInt32Pointer()
		target.HCInterval = healthcheck.Interval.ValueInt32Pointer()
		target.HCUnhealthyInterval = healthcheck.UnhealthyInterval.ValueInt32Pointer()
		target.HCTimeout = healthcheck.Timeout.ValueInt32Pointer()
		target.HCHeaders = hc_headers
		target.HCFollowRedirects = healthcheck.FollowRedirects.ValueBoolPointer()
		target.HCMethod = healthcheck.Method.ValueStringPointer()
		target.HCStatus = healthcheck.Status.ValueInt32Pointer()
		target.HCTlsServerName = healthcheck.TlsServerName.ValueStringPointer()
	}

	created, err := r.client.CreateTarget(data.ResourceID.ValueInt64(), target)
	if err != nil {
		resp.Diagnostics.AddError("Error creating target", err.Error())
		return
	}

	data.ID = types.Int64PointerValue(created.ID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *targetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data targetResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	target, err := r.client.GetTarget(data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading target", err.Error())
		return
	}

	data.ID = types.Int64PointerValue(target.ID)
	data.Enabled = types.BoolPointerValue(target.Enabled)
	data.IP = types.StringPointerValue(target.IP)
	data.Method = types.StringPointerValue(target.Method)
	data.Path = types.StringPointerValue(target.Path)
	data.PathMatchType = types.StringPointerValue(target.PathMatchType)
	data.Port = types.Int32PointerValue(target.Port)
	data.Priority = types.Int32PointerValue(target.Priority)
	data.ResourceID = types.Int64PointerValue(target.ResourceID)
	data.RewritePath = types.StringPointerValue(target.RewritePath)
	data.RewritePathType = types.StringPointerValue(target.RewritePathType)
	data.SiteID = types.Int64PointerValue(target.SiteID)

	if target.HCHostname != nil && target.Port != nil {
		healthcheck := targetHealthCheck{
			Enabled:           types.BoolPointerValue(target.Enabled),
			Path:              types.StringPointerValue(target.HCPath),
			Scheme:            types.StringPointerValue(target.HCScheme),
			Mode:              types.StringPointerValue(target.HCMode),
			Hostname:          types.StringPointerValue(target.HCHostname),
			Port:              types.Int32PointerValue(target.HCPort),
			Interval:          types.Int32PointerValue(target.HCInterval),
			UnhealthyInterval: types.Int32PointerValue(target.HCUnhealthyInterval),
			Timeout:           types.Int32PointerValue(target.HCTimeout),
			FollowRedirects:   types.BoolPointerValue(target.HCFollowRedirects),
			Method:            types.StringPointerValue(target.HCMethod),
			Status:            types.Int32PointerValue(target.HCStatus),
			TlsServerName:     types.StringPointerValue(target.HCTlsServerName),
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("health_check"), healthcheck)...)

		if target.HCHeaders != nil {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("health_check.headers"), target.HCHeaders)...)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *targetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state targetResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var healthcheck *targetHealthCheck
	var hc_headers []client.TargetHeader
	resp.Diagnostics.Append(data.HealthCheck.As(ctx, &healthcheck, basetypes.ObjectAsOptions{})...)

	if healthcheck != nil && !healthcheck.Headers.IsNull() {
		h := make([]targetHCHeader, len(healthcheck.Headers.Elements()))
		resp.Diagnostics.Append(healthcheck.Headers.ElementsAs(ctx, &h, false)...)

		hc_headers = make([]client.TargetHeader, len(h))

		for i, header := range h {
			hc_headers[i] = client.TargetHeader{
				Name:  header.Name.ValueString(),
				Value: header.Value.ValueString(),
			}
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	target := &client.Target{
		SiteID:          data.SiteID.ValueInt64Pointer(),
		IP:              data.IP.ValueStringPointer(),
		Port:            data.Port.ValueInt32Pointer(),
		Method:          data.Method.ValueStringPointer(),
		Enabled:         data.Enabled.ValueBoolPointer(),
		Path:            data.Path.ValueStringPointer(),
		PathMatchType:   data.PathMatchType.ValueStringPointer(),
		RewritePath:     data.RewritePath.ValueStringPointer(),
		RewritePathType: data.RewritePathType.ValueStringPointer(),
		Priority:        data.Priority.ValueInt32Pointer(),
	}

	if healthcheck != nil {
		target.HCEnabled = healthcheck.Enabled.ValueBoolPointer()
		target.HCPath = healthcheck.Path.ValueStringPointer()
		target.HCScheme = healthcheck.Scheme.ValueStringPointer()
		target.HCMode = healthcheck.Mode.ValueStringPointer()
		target.HCHostname = healthcheck.Hostname.ValueStringPointer()
		target.HCPort = healthcheck.Port.ValueInt32Pointer()
		target.HCInterval = healthcheck.Interval.ValueInt32Pointer()
		target.HCUnhealthyInterval = healthcheck.UnhealthyInterval.ValueInt32Pointer()
		target.HCTimeout = healthcheck.Timeout.ValueInt32Pointer()
		target.HCHeaders = hc_headers
		target.HCFollowRedirects = healthcheck.FollowRedirects.ValueBoolPointer()
		target.HCMethod = healthcheck.Method.ValueStringPointer()
		target.HCStatus = healthcheck.Status.ValueInt32Pointer()
		target.HCTlsServerName = healthcheck.TlsServerName.ValueStringPointer()
	}

	_, err := r.client.UpdateTarget(state.ID.ValueInt64(), target)
	if err != nil {
		resp.Diagnostics.AddError("Error updating target", err.Error())
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *targetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data targetResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTarget(data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting target", err.Error())
		return
	}
}

func (r *targetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected id to be an integer. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
