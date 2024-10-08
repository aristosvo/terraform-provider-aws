// Code generated by internal/generate/servicepackage/main.go; DO NOT EDIT.

package identitystore

import (
	"context"

	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	identitystore_sdkv2 "github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{
		{
			Factory: newGroupsDataSource,
			Name:    "Groups",
		},
	}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  dataSourceGroup,
			TypeName: "aws_identitystore_group",
			Name:     "Group",
		},
		{
			Factory:  dataSourceUser,
			TypeName: "aws_identitystore_user",
			Name:     "User",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  resourceGroup,
			TypeName: "aws_identitystore_group",
			Name:     "Group",
		},
		{
			Factory:  resourceGroupMembership,
			TypeName: "aws_identitystore_group_membership",
			Name:     "Group Membership",
		},
		{
			Factory:  resourceUser,
			TypeName: "aws_identitystore_user",
			Name:     "User",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.IdentityStore
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*identitystore_sdkv2.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws_sdkv2.Config))

	return identitystore_sdkv2.NewFromConfig(cfg,
		identitystore_sdkv2.WithEndpointResolverV2(newEndpointResolverSDKv2()),
		withBaseEndpoint(config[names.AttrEndpoint].(string)),
	), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
