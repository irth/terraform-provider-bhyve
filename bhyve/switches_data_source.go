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
	client client.Client
}

type switchesDataSourceModel struct {
	Switches map[string]switchesSwitchModel `tfsdk:"switches"`
}

type switchesSwitchModel struct {
	Address types.String `tfsdk:"address"`
	// TODO: support others
}

func (d *switchesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_switches"
}

func (d *switchesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"switches": schema.MapNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
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

	d.client = req.ProviderData.(client.Client)
}

func (d *switchesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state switchesDataSourceModel

	switches, err := d.client.SwitchList()
	if err != nil {
		resp.Diagnostics.AddError("failed to read switches", err.Error())
		return
	}

	state.Switches = map[string]switchesSwitchModel{}

	for _, sw := range switches {
		swM := switchesSwitchModel{}
		if sw.Address != "" {
			swM.Address = types.StringValue(sw.Address)
		}
		state.Switches[sw.Name] = swM
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
