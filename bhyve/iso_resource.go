package bhyve

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/irth/terraform-provider-bhyve/bhyve/client"
)

type isoResource struct {
	client client.Client
}

var (
	_ resource.Resource              = &isoResource{}
	_ resource.ResourceWithConfigure = &isoResource{}
)

func NewIsoResource() resource.Resource {
	return &isoResource{}
}

type isoResourceModel struct {
	Image    types.Bool   `tfsdk:"image"`
	Name     types.String `tfsdk:"name"`
	URL      types.String `tfsdk:"url"`
	Checksum types.String `tfsdk:"sha256sum"`
}

func (r *isoResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iso"
}

func (r *isoResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"image": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sha256sum": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *isoResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(client.Client)
}

func (r *isoResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan isoResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	if !plan.Image.ValueBool() {
		err = r.client.ISO(plan.URL.ValueString(), plan.Name.ValueString(), plan.Checksum.ValueString())
	} else {
		err = r.client.IMG(plan.URL.ValueString(), plan.Name.ValueString(), plan.Checksum.ValueString())
	}

	if err != nil {
		var cmdExecError client.CommandExecutionError
		if errors.As(err, &cmdExecError) {
			tflog.Error(ctx, "command execution error", map[string]any{
				"stderr":     cmdExecError.Stderr,
				"cmd":        cmdExecError.Cmd,
				"returnCode": cmdExecError.ReturnCode,
			})
		}
		tflog.Error(ctx, "error", map[string]any{
			"error": err,
		})
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *isoResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state isoResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var checksum string
	var err error
	if state.Image.ValueBool() {
		checksum, err = r.client.ChecksumIMG(state.Name.ValueString())
	} else {
		checksum, err = r.client.ChecksumISO(state.Name.ValueString())
	}
	if err != nil {
		if errors.As(err, &client.ErrFileNotFound{}) {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("failed to read", err.Error())
		}
		return
	}

	state.Checksum = types.StringValue(checksum)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *isoResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("should not happen", "all attributes force replacement, so update should never be called")
}

func (r *isoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state isoResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	if !state.Image.ValueBool() {
		err = r.client.RemoveISO(state.Name.ValueString())
	} else {
		err = r.client.RemoveIMG(state.Name.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError("failed to delete", err.Error())
		return
	}
}
