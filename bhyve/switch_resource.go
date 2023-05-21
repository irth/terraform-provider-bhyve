package bhyve

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/irth/terraform-provider-bhyve/bhyve/client"
)

type switchResource struct {
	client client.Executor
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
			},
			"address": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *switchResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(client.Executor)
}

func (r *switchResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan switchResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: validate address is cidr
	// TODO: move this code to client

	params := []string{"switch", "create", plan.Name.ValueString()}
	if !plan.Address.IsNull() {
		params = append(params, "-a", plan.Address.ValueString())
	}

	_, err := r.client.RunCmd("vm", params...)
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

	// TODO: use vm switch info isntead
	switches := client.Switches{}
	err := switches.LoadFromSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError("failed to read", err.Error())
		return
	}

	switchesMap := switches.AsMap()
	if sw, ok := switchesMap[state.Name.ValueString()]; ok {
		if sw.Address != "" {
			state.Address = types.StringValue(sw.Address)
		}
	} else {
		resp.Diagnostics.AddError("failed to read", "switch not found")
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *switchResource) Update(_ context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	panic("implement me")
}

func (r *switchResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	panic("implement me")
}
