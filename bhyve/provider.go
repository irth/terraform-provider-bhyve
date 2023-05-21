package bhyve

import (
	"context"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/irth/terraform-provider-bhyve/bhyve/client"
)

var (
	_ provider.Provider = (*bhyveProvider)(nil)
)

type bhyveProvider struct{}

type bhyveProviderModel struct {
	Host      types.String `tfsdk:"host"`
	Port      types.Int64  `tfsdk:"port"`
	User      types.String `tfsdk:"user"`
	SSH       types.String `tfsdk:"ssh"`
	SSHParams types.List   `tfsdk:"ssh_params"`
}

func New() provider.Provider {
	return &bhyveProvider{}
}

func (p *bhyveProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bhyve"
}

func (p *bhyveProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Required: true,
			},
			"port": schema.Int64Attribute{
				Optional: true,
			},
			"user": schema.StringAttribute{
				Optional: true,
			},
			"ssh": schema.StringAttribute{
				Optional: true,
			},
			"ssh_params": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (p *bhyveProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config bhyveProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"bhyve host address cannot be unknown",
			"The vm-bhyve provider cannot work without a bhyve host.",
		)
	}

	// TODO: check others

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("BHYVE_HOST")
	portStr := os.Getenv("BHYVE_PORT")
	port := 22
	if portStr != "" {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			resp.Diagnostics.AddError("BHYVE_PORT env invalid", "BHYVE_PORT environment variable is not a valid port number")
			return
		}
	}
	user := os.Getenv("BHYVE_USER")

	ssh := os.Getenv("BHYVE_SSH")

	sshParams := []string{}
	// TODO: support ssh params from env

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Port.IsNull() {
		port = int(config.Port.ValueInt64())
	}

	if !config.SSH.IsNull() {
		ssh = config.SSH.ValueString()
	}

	if !config.SSHParams.IsNull() {
		err := config.SSHParams.ElementsAs(ctx, sshParams, false)
		if err != nil {
			resp.Diagnostics.AddError("ssh_params is not a valid list of strings", "ssh_params must be a list of strings")
			return
		}
	}

	if user == "" {
		user = "root"
	}

	if ssh == "" {
		ssh = "ssh"
	}

	if port == 0 {
		port = 22
	}

	executor := client.SSHExecutor{
		Host:   host,
		Port:   port,
		User:   user,
		Rsh:    ssh,
		Params: sshParams,
	}

	resp.DataSourceData = executor
	resp.ResourceData = executor
}

func (p *bhyveProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSwitchesDataSource,
	}
}

func (p *bhyveProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSwitchResource,
	}
}
