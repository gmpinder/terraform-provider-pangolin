package provider

import (
	"context"
	"fmt"

	"github.com/gmpinder/terraform-provider-pangolin/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &organizationResource{}
var _ resource.ResourceWithImportState = &organizationResource{}

func NewOrganizationResource() resource.Resource {
	return &organizationResource{}
}

type organizationResource struct {
	client *client.Client
}

type organizationResourceModel struct {
	ID                              types.String `tfsdk:"id"`
	Name                            types.String `tfsdk:"name"`
	Subnet                          types.String `tfsdk:"subnet"`
	UtilitySubnet                   types.String `tfsdk:"utility_subnet"`
	RequireTwoFactor                types.Bool   `tfsdk:"require_two_factor"`
	MaxSessionLengthHours           types.Int32  `tfsdk:"max_session_length_hours"`
	PasswordExpiryDays              types.Int32  `tfsdk:"password_expiry_days"`
	SettingsLogRetentionDaysRequest types.Int32  `tfsdk:"settings_log_retention_days_request"`
	SettingsLogRetentionDaysAccess  types.Int32  `tfsdk:"settings_log_retention_days_access"`
	SettingsLogRetentionDaysAction  types.Int32  `tfsdk:"settings_log_retention_days_action"`
}

func (o *organizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (o *organizationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an organization organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the organization.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the organization.",
			},
			"subnet": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("100.90.128.0/24"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"utility_subnet": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("100.96.128.0/24"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"require_two_factor": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"max_session_length_hours": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "The maximum length in hours for a valid session.",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"password_expiry_days": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "The number of days before a password expires.",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"settings_log_retention_days_request": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The number of days to retain request logs",
				Validators: []validator.Int32{
					int32validator.AtLeast(-1),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"settings_log_retention_days_access": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The number of days to retain access logs",
				Validators: []validator.Int32{
					int32validator.AtLeast(-1),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"settings_log_retention_days_action": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The number of days to retain action logs",
				Validators: []validator.Int32{
					int32validator.AtLeast(-1),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (o *organizationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}

	o.client = c
}

func (o *organizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data organizationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org := &client.Organization{
		ID:            data.ID.ValueStringPointer(),
		Name:          data.Name.ValueStringPointer(),
		Subnet:        data.Subnet.ValueStringPointer(),
		UtilitySubnet: data.UtilitySubnet.ValueStringPointer(),
	}

	created, err := o.client.CreateOrganization(org)

	if err != nil {
		resp.Diagnostics.AddError("error creating organization", err.Error())
		return
	}

	org = &client.Organization{
		RequireTwoFactor:                data.RequireTwoFactor.ValueBoolPointer(),
		MaxSessionLengthHours:           data.MaxSessionLengthHours.ValueInt32Pointer(),
		PasswordExpiryDays:              data.PasswordExpiryDays.ValueInt32Pointer(),
		SettingsLogRetentionDaysRequest: data.SettingsLogRetentionDaysRequest.ValueInt32Pointer(),
		SettingsLogRetentionDaysAccess:  data.SettingsLogRetentionDaysAccess.ValueInt32Pointer(),
		SettingsLogRetentionDaysAction:  data.SettingsLogRetentionDaysAction.ValueInt32Pointer(),
	}

	updated, err := o.client.UpdateOrganization(*created.ID, org)

	if err != nil {
		resp.Diagnostics.AddError("error updating organization during create", err.Error())
	} else {
		if data.RequireTwoFactor.IsUnknown() {
			data.RequireTwoFactor = types.BoolPointerValue(updated.RequireTwoFactor)
		}

		if data.SettingsLogRetentionDaysAccess.IsUnknown() {
			data.SettingsLogRetentionDaysAccess = types.Int32PointerValue(updated.SettingsLogRetentionDaysAccess)
		}

		if data.SettingsLogRetentionDaysAction.IsUnknown() {
			data.SettingsLogRetentionDaysAction = types.Int32PointerValue(updated.SettingsLogRetentionDaysAction)
		}

		if data.SettingsLogRetentionDaysRequest.IsUnknown() {
			data.SettingsLogRetentionDaysRequest = types.Int32PointerValue(updated.SettingsLogRetentionDaysRequest)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (o *organizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data organizationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org, err := o.client.GetOrganization(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("error getting organization", err.Error())
		return
	}

	data.MaxSessionLengthHours = types.Int32PointerValue(org.MaxSessionLengthHours)
	data.Name = types.StringPointerValue(org.Name)
	data.PasswordExpiryDays = types.Int32PointerValue(org.PasswordExpiryDays)
	data.RequireTwoFactor = types.BoolPointerValue(org.RequireTwoFactor)
	data.SettingsLogRetentionDaysAccess = types.Int32PointerValue(org.SettingsLogRetentionDaysAccess)
	data.SettingsLogRetentionDaysAction = types.Int32PointerValue(org.SettingsLogRetentionDaysAction)
	data.SettingsLogRetentionDaysRequest = types.Int32PointerValue(org.SettingsLogRetentionDaysRequest)
	data.Subnet = types.StringPointerValue(org.Subnet)
	data.UtilitySubnet = types.StringPointerValue(org.UtilitySubnet)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (o *organizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state organizationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org := &client.Organization{
		Name:                            data.Name.ValueStringPointer(),
		RequireTwoFactor:                data.RequireTwoFactor.ValueBoolPointer(),
		MaxSessionLengthHours:           data.MaxSessionLengthHours.ValueInt32Pointer(),
		PasswordExpiryDays:              data.PasswordExpiryDays.ValueInt32Pointer(),
		SettingsLogRetentionDaysRequest: data.SettingsLogRetentionDaysRequest.ValueInt32Pointer(),
		SettingsLogRetentionDaysAccess:  data.SettingsLogRetentionDaysAccess.ValueInt32Pointer(),
		SettingsLogRetentionDaysAction:  data.SettingsLogRetentionDaysAction.ValueInt32Pointer(),
	}

	updated, err := o.client.UpdateOrganization(state.ID.ValueString(), org)

	if err != nil {
		resp.Diagnostics.AddError("error updating organization", err.Error())
		return
	}

	data.MaxSessionLengthHours = types.Int32PointerValue(updated.MaxSessionLengthHours)
	data.Name = types.StringPointerValue(updated.Name)
	data.PasswordExpiryDays = types.Int32PointerValue(updated.PasswordExpiryDays)
	data.RequireTwoFactor = types.BoolPointerValue(updated.RequireTwoFactor)
	data.SettingsLogRetentionDaysAccess = types.Int32PointerValue(updated.SettingsLogRetentionDaysAccess)
	data.SettingsLogRetentionDaysAction = types.Int32PointerValue(updated.SettingsLogRetentionDaysAction)
	data.SettingsLogRetentionDaysRequest = types.Int32PointerValue(updated.SettingsLogRetentionDaysRequest)
	data.Subnet = types.StringPointerValue(updated.Subnet)
	data.UtilitySubnet = types.StringPointerValue(updated.UtilitySubnet)

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (o *organizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data organizationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := o.client.DeleteOrganization(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error deleting organization", err.Error())
		return
	}
}

func (o *organizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
