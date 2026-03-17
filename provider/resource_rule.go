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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ruleResource{}
var _ resource.ResourceWithImportState = &ruleResource{}

func NewRuleResource() resource.Resource {
	return &ruleResource{}
}

type ruleResource struct {
	client *client.Client
}

type ruleResourceModel struct {
	ID         types.Int64  `tfsdk:"id"`
	ResourceID types.Int64  `tfsdk:"resource_id"`
	Action     types.String `tfsdk:"action"`
	Match      types.String `tfsdk:"match"`
	Value      types.String `tfsdk:"value"`
	Priority   types.Int32  `tfsdk:"priority"`
	Enabled    types.Bool   `tfsdk:"enabled"`
}

func (r *ruleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rule"
}

func (r *ruleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a rule for a resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
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
			"action": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The action to perform for the match.",
				Validators: []validator.String{
					stringvalidator.OneOf("ACCEPT", "DROP", "PASS"),
				},
			},
			"match": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "What category to match on.",
				Validators: []validator.String{
					stringvalidator.OneOf("CIDR", "IP", "PATH", "COUNTRY", "ASN"),
				},
			},
			"value": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The value to match.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"priority": schema.Int32Attribute{
				Required:            true,
				MarkdownDescription: "The priority to evaluate the rules. Lower number is evaluated first.",
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			"enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
		},
	}
}

func (r *ruleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ruleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ruleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule := &client.Rule{
		Action:   data.Action.ValueString(),
		Match:    data.Match.ValueString(),
		Value:    data.Value.ValueString(),
		Priority: data.Priority.ValueInt32(),
		Enabled:  data.Enabled.ValueBool(),
	}

	createdRule, err := r.client.CreateRule(data.ResourceID.ValueInt64(), rule)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create rule", err.Error())
		return
	}

	data.ID = types.Int64PointerValue(createdRule.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ruleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ruleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.GetRule(data.ResourceID.ValueInt64(), data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get rule", err.Error())
		return
	}

	data.ID = types.Int64PointerValue(rule.ID)
	data.Action = types.StringValue(rule.Action)
	data.Enabled = types.BoolValue(rule.Enabled)
	data.Match = types.StringValue(rule.Match)
	data.Priority = types.Int32Value(rule.Priority)
	data.ResourceID = types.Int64PointerValue(rule.ResourceID)
	data.Value = types.StringValue(rule.Value)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ruleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state ruleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule := &client.Rule{
		Action:   data.Action.ValueString(),
		Match:    data.Match.ValueString(),
		Value:    data.Value.ValueString(),
		Priority: data.Priority.ValueInt32(),
		Enabled:  data.Enabled.ValueBool(),
	}

	_, err := r.client.UpdateRule(data.ResourceID.ValueInt64(), state.ID.ValueInt64(), rule)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update rule", err.Error())
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ruleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ruleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRule(data.ResourceID.ValueInt64(), data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting rule", err.Error())
		return
	}
}

func (r *ruleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
