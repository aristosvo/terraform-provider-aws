// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package controltower

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/controltower"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

// @SDKDataSource("aws_controltower_controls")
func DataSourceControls() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: DataSourceControlsRead,

		Schema: map[string]*schema.Schema{
			"enabled_controls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"target_identifier": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidARN,
			},
		},
	}
}

func DataSourceControlsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).ControlTowerClient(ctx)

	targetIdentifier := d.Get("target_identifier").(string)
	input := &controltower.ListEnabledControlsInput{
		TargetIdentifier: aws.String(targetIdentifier),
	}

	var controls []string
	paginator := controltower.NewListEnabledControlsPaginator(conn, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return diag.Errorf("listing ControlTower Controls (%s): %s", targetIdentifier, err)
		}
		for _, v := range page.EnabledControls {
			if v.ControlIdentifier == nil {
				continue
			}
			controls = append(controls, aws.ToString(v.ControlIdentifier))
		}
	}

	d.SetId(targetIdentifier)
	d.Set("enabled_controls", controls)

	return nil
}
