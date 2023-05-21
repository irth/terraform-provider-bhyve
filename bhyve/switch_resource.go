package bhyve

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/irth/terraform-provider-bhyve/bhyve/client"
)

type switchResource struct {
	client client.Client
}

var (
	_ resource.Resource              = &switchResource{}
	_ resource.ResourceWithConfigure = &switchResource{}
)

func NewSwitchResource() resource.Resource {
	return &switchResource{}
}

type switchResourceModel struct {
	Name    types.String `tfsdk:"name"`
	Address types.String `tfsdk:"address"`
}

func (r *switchResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_switch"
}

func (r *switchResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"address": schema.StringAttribute{
				Optional: true,
			},
		},
		Blocks:              map[string]schema.Block{},
		Description:         "",
		MarkdownDescription: "",
		DeprecationMessage:  "",
		Version:             0,
	}
}

func (r *switchResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(client.Client)
}

func (r *switchResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan switchResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.SwitchCreate(
		&client.Switch{
			Name:    plan.Name.ValueString(),
			Address: plan.Address.ValueString(),
		},
	)

	if err != nil {
		resp.Diagnostics.AddError("failed to create", err.Error())
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *switchResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state switchResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sw, err := r.client.SwitchInfo(state.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("failed to read", err.Error())
		return
	}

	if sw.Address != "" {
		state.Address = types.StringValue(sw.Address)
	} else {
		state.Address = types.StringNull()
	}

	state.Name = types.StringValue(sw.Name)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *switchResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan switchResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	if plan.Address.IsNull() {
		err = r.client.SwitchAddress(plan.Name.ValueString(), "")
	} else {
		err = r.client.SwitchAddress(plan.Name.ValueString(), plan.Address.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError("failed to update", err.Error())
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *switchResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan switchResourceModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.SwitchDestroy(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to delete", err.Error())
		return
	}
}
