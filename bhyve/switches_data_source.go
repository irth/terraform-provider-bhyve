package bhyve

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/irth/terraform-provider-bhyve/bhyve/client"
)

func NewSwitchesDataSource() datasource.DataSource {
	return &switchesDataSource{}
}

var (
	_ datasource.DataSource              = &switchesDataSource{}
	_ datasource.DataSourceWithConfigure = &switchesDataSource{}
)

type switchesDataSource struct {
	client client.Executor
}

type switchesDataSourceModel struct {
	Switches []switchesSwitchModel `tfsdk:"switches"`
}

type switchesSwitchModel struct {
	Name    types.String `tfsdk:"name"`
	Address types.String `tfsdk:"address"`
	// TODO: support others
}

func (d *switchesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_switches"
}

func (d *switchesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"switches": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"address": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *switchesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(client.Executor)
}

func (d *switchesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state switchesDataSourceModel

	switches := client.Switches{}
	err := switches.LoadFromSystem(d.client)
	if err != nil {
		resp.Diagnostics.AddError("failed to read switches", err.Error())
		return
	}

	state.Switches = []switchesSwitchModel{}

	for _, sw := range switches {
		swM := switchesSwitchModel{
			Name: types.StringValue(sw.Name),
		}
		if sw.Address != "" {
			swM.Address = types.StringValue(sw.Address)
		}
		state.Switches = append(state.Switches, swM)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
