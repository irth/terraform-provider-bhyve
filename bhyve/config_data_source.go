package bhyve

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/irth/terraform-provider-bhyve/bhyve/client"
)

func NewConfigDataSource() datasource.DataSource {
	return &configDataSource{}
}

var (
	_ datasource.DataSource              = &configDataSource{}
	_ datasource.DataSourceWithConfigure = &configDataSource{}
)

type configDataSource struct {
	client client.Client
}

type configDataSourceModel struct {
	BhyveEnabled types.Bool   `tfsdk:"bhyve_enabled"`
	VMEnabled    types.Bool   `tfsdk:"vm_enabled"`
	VMDir        types.String `tfsdk:"vm_dir"`
}

func (d *configDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config"
}

func (d *configDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bhyve_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"vm_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"vm_dir": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *configDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(client.Client)
}

func (d *configDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state configDataSourceModel

	config, err := d.client.Config()
	if err != nil {
		resp.Diagnostics.AddError("failed to read", err.Error())
		return
	}

	state.BhyveEnabled = types.BoolValue(config.BhyveEnabled)
	state.VMEnabled = types.BoolValue(config.VMEnabled)
	state.VMDir = types.StringValue(config.VMDir)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
