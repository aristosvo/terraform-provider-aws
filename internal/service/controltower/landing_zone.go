// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package controltower

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/controltower"
	"github.com/aws/aws-sdk-go-v2/service/controltower/document"
	types "github.com/aws/aws-sdk-go-v2/service/controltower/types"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// @SDKResource("aws_controltower_landing_zone")
func ResourceLandingZone() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceLandingZoneCreate,
		ReadWithoutTimeout:   resourceLandingZoneRead,
		DeleteWithoutTimeout: resourceLandingZoneDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"manifest": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

const (
	ARNSeparator = "/"
	ARNService   = "controltower"

	LandingZoneResourcePrefix = "landingzone"
)

// LandingZoneARNToIdentifier converts Amazon Resource Name (ARN) to the Identifier.
func LandingZoneARNToIdentifier(inputARN string) (string, error) {
	parsedARN, err := arn.Parse(inputARN)

	if err != nil {
		return "", fmt.Errorf("parsing ARN (%s): %w", inputARN, err)
	}

	if actual, expected := parsedARN.Service, ARNService; actual != expected {
		return "", fmt.Errorf("expected service %s in ARN (%s), got: %s", expected, inputARN, actual)
	}

	resourceParts := strings.Split(parsedARN.Resource, ARNSeparator)

	if actual, expected := len(resourceParts), 2; actual < expected {
		return "", fmt.Errorf("expected at least %d resource parts in ARN (%s), got: %d", expected, inputARN, actual)
	}

	if actual, expected := resourceParts[0], LandingZoneResourcePrefix; actual != expected {
		return "", fmt.Errorf("expected resource prefix %s in ARN (%s), got: %s", expected, inputARN, actual)
	}

	return resourceParts[len(resourceParts)-1], nil
}

func resourceLandingZoneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).ControlTowerClient(ctx)

	manifest := d.Get("manifest").(string)
	version := d.Get("version").(string)
	input := &controltower.CreateLandingZoneInput{
		Manifest: document.NewLazyDocument(manifest),
		Version:  aws.String(version),
	}

	output, err := conn.CreateLandingZone(ctx, input)
	if err != nil {
		return diag.Errorf("creating ControlTower Landing Zone: %s", err)
	}
	id, err := LandingZoneARNToIdentifier(aws.ToString(output.Arn))
	if err != nil {
		return diag.Errorf("parsing ControlTower Landing Zone Arn: %s", err)
	}

	d.SetId(id)

	if _, err := waitLandingZoneOperationSucceeded(ctx, conn, aws.ToString(output.OperationIdentifier), d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("waiting for ControlTower Landing Zone (%s) create: %s", d.Id(), err)
	}

	return resourceLandingZoneRead(ctx, d, meta)
}

func resourceLandingZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).ControlTowerClient(ctx)
	output, err := conn.GetLandingZone(ctx, &controltower.GetLandingZoneInput{
		LandingZoneIdentifier: aws.String(d.Id()),
	})

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] ControlTower Landing Zone %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading ControlTower Control (%s): %s", d.Id(), err)
	}

	d.Set("manifest", document.output.LandingZone.Manifest)
	d.Set("version", output.LandingZone.Version)

	return nil
}

func resourceLandingZoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).ControlTowerClient(ctx)

	log.Printf("[DEBUG] Deleting ControlTower Landing Zone: %s", d.Id())
	output, err := conn.DeleteLandingZone(ctx, &controltower.DeleteLandingZoneInput{
		LandingZoneIdentifier: aws.String(d.Id()),
	})

	if err != nil {
		return diag.Errorf("deleting ControlTower Landing Zone (%s): %s", d.Id(), err)
	}

	if _, err := waitLandingZoneOperationSucceeded(ctx, conn, aws.ToString(output.OperationIdentifier), d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("waiting for ControlTower Landing Zone (%s) delete: %s", d.Id(), err)
	}

	return nil
}

func findLandingZoneOperationByID(ctx context.Context, conn *controltower.Client, id string) (*types.LandingZoneOperationDetail, error) {
	input := &controltower.GetLandingZoneOperationInput{
		OperationIdentifier: aws.String(id),
	}

	output, err := conn.GetLandingZoneOperation(ctx, input)

	if tfresource.NotFound(err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.OperationDetails == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.OperationDetails, nil
}

func statusLandingZoneOperation(ctx context.Context, conn *controltower.Client, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := findLandingZoneOperationByID(ctx, conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, string(output.Status), nil
	}
}

func waitLandingZoneOperationSucceeded(ctx context.Context, conn *controltower.Client, id string, timeout time.Duration) (*types.LandingZoneOperationDetail, error) { //nolint:unparam
	stateConf := &retry.StateChangeConf{
		Pending: []string{string(types.LandingZoneOperationStatusInProgress)},
		Target:  []string{string(types.LandingZoneOperationStatusSucceeded)},
		Refresh: statusLandingZoneOperation(ctx, conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*types.LandingZoneOperationDetail); ok {
		if status := output.Status; status == types.LandingZoneOperationStatusFailed {
			tfresource.SetLastError(err, errors.New(aws.ToString(output.StatusMessage)))
		}

		return output, err
	}

	return nil, err
}
