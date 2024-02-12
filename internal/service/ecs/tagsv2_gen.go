// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package ecs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	awstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/hashicorp/aws-sdk-go-base/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/logging"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/types/option"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// GetTag fetches an individual ecs service tag for a resource.
// Returns whether the key value and any errors. A NotFoundError is used to signal that no value was found.
// This function will optimise the handling over listTagsV2, if possible.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func GetTag(ctx context.Context, conn *ecs.Client, identifier, key string, optFns ...func(*ecs.Options)) (*string, error) {
	listTags, err := listTagsV2(ctx, conn, identifier, optFns...)

	if err != nil {
		return nil, err
	}

	if !listTags.KeyExists(key) {
		return nil, tfresource.NewEmptyResultError(nil)
	}

	return listTags.KeyValue(key), nil
}

// listTagsV2 lists ecs service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func listTagsV2(ctx context.Context, conn *ecs.Client, identifier string, optFns ...func(*ecs.Options)) (tftags.KeyValueTags, error) {
	input := &ecs.ListTagsForResourceInput{
		ResourceArn: aws.String(identifier),
	}

	output, err := conn.ListTagsForResource(ctx, input, optFns...)

	if tfawserr.ErrMessageContains(err, "InvalidParameterException", "The specified cluster is inactive. Specify an active cluster and try again.") {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return tftags.New(ctx, nil), err
	}

	return keyValueTagsV2(ctx, output.Tags), nil
}

// []*SERVICE.Tag handling

// TagsV2 returns ecs service tags.
func TagsV2(tags tftags.KeyValueTags) []awstypes.Tag {
	result := make([]awstypes.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := awstypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// keyValueTagsV2 creates tftags.KeyValueTags from ecs service tags.
func keyValueTagsV2(ctx context.Context, tags []awstypes.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.ToString(tag.Key)] = tag.Value
	}

	return tftags.New(ctx, m)
}

// getTagsInV2 returns ecs service tags from Context.
// nil is returned if there are no input tags.
func getTagsInV2(ctx context.Context) []awstypes.Tag {
	if inContext, ok := tftags.FromContext(ctx); ok {
		if tags := TagsV2(inContext.TagsIn.UnwrapOrDefault()); len(tags) > 0 {
			return tags
		}
	}

	return nil
}

// setTagsOutV2 sets ecs service tags in Context.
func setTagsOutV2(ctx context.Context, tags []awstypes.Tag) {
	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = option.Some(keyValueTagsV2(ctx, tags))
	}
}

// updateTagsV2 updates ecs service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func updateTagsV2(ctx context.Context, conn *ecs.Client, identifier string, oldTagsMap, newTagsMap any, optFns ...func(*ecs.Options)) error {
	oldTags := tftags.New(ctx, oldTagsMap)
	newTags := tftags.New(ctx, newTagsMap)

	ctx = tflog.SetField(ctx, logging.KeyResourceId, identifier)

	removedTags := oldTags.Removed(newTags)
	removedTags = removedTags.IgnoreSystem(names.ECS)
	if len(removedTags) > 0 {
		input := &ecs.UntagResourceInput{
			ResourceArn: aws.String(identifier),
			TagKeys:     removedTags.Keys(),
		}

		_, err := conn.UntagResource(ctx, input, optFns...)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	updatedTags := oldTags.Updated(newTags)
	updatedTags = updatedTags.IgnoreSystem(names.ECS)
	if len(updatedTags) > 0 {
		input := &ecs.TagResourceInput{
			ResourceArn: aws.String(identifier),
			Tags:        TagsV2(updatedTags),
		}

		_, err := conn.TagResource(ctx, input, optFns...)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}