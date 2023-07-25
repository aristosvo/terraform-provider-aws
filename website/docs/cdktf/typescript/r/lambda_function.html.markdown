---
subcategory: "Lambda"
layout: "aws"
page_title: "AWS: aws_lambda_function"
description: |-
  Provides a Lambda Function resource. Lambda allows you to trigger execution of code in response to events in AWS, enabling serverless backend solutions. The Lambda Function itself includes source code and runtime configuration.
---


<!-- Please do not edit this file, it is generated. -->
# Resource: aws_lambda_function

Provides a Lambda Function resource. Lambda allows you to trigger execution of code in response to events in AWS, enabling serverless backend solutions. The Lambda Function itself includes source code and runtime configuration.

For information about Lambda and how to use it, see [What is AWS Lambda?][1]

For a detailed example of setting up Lambda and API Gateway, see [Serverless Applications with AWS Lambda and API Gateway.][11]

~> **NOTE:** Due to [AWS Lambda improved VPC networking changes that began deploying in September 2019](https://aws.amazon.com/blogs/compute/announcing-improved-vpc-networking-for-aws-lambda-functions/), EC2 subnets and security groups associated with Lambda Functions can take up to 45 minutes to successfully delete. Terraform AWS Provider version 2.31.0 and later automatically handles this increased timeout, however prior versions require setting the customizable deletion timeouts of those Terraform resources to 45 minutes (`delete = "45m"`). AWS and HashiCorp are working together to reduce the amount of time required for resource deletion and updates can be tracked in this [GitHub issue](https://github.com/hashicorp/terraform-provider-aws/issues/10329).

~> **NOTE:** If you get a `KMSAccessDeniedException: Lambda was unable to decrypt the environment variables because KMS access was denied` error when invoking an [`awsLambdaFunction`](/docs/providers/aws/r/lambda_function.html) with environment variables, the IAM role associated with the function may have been deleted and recreated _after_ the function was created. You can fix the problem two ways: 1) updating the function's role to another role and then updating it back again to the recreated role, or 2) by using Terraform to `taint` the function and `apply` your configuration again to recreate the function. (When you create a function, Lambda grants permissions on the KMS key to the function's IAM role. If the IAM role is recreated, the grant is no longer valid. Changing the function's role or recreating the function causes Lambda to update the grant.)

-> To give an external source (like an EventBridge Rule, SNS, or S3) permission to access the Lambda function, use the [`awsLambdaPermission`](lambda_permission.html) resource. See [Lambda Permission Model][4] for more details. On the other hand, the `role` argument of this resource is the function's execution role for identity and access to AWS services and resources.

## Example Usage

### Basic Example

```typescript
// Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { Token, TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { DataArchiveFile } from "./.gen/providers/archive/data-archive-file";
import { DataAwsIamPolicyDocument } from "./.gen/providers/aws/data-aws-iam-policy-document";
import { IamRole } from "./.gen/providers/aws/iam-role";
import { LambdaFunction } from "./.gen/providers/aws/lambda-function";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    /*The following providers are missing schema information and might need manual adjustments to synthesize correctly: archive.
    For a more precise conversion please use the --provider flag in convert.*/
    const lambda = new DataArchiveFile(this, "lambda", {
      output_path: "lambda_function_payload.zip",
      source_file: "lambda.js",
      type: "zip",
    });
    const assumeRole = new DataAwsIamPolicyDocument(this, "assume_role", {
      statement: [
        {
          actions: ["sts:AssumeRole"],
          effect: "Allow",
          principals: [
            {
              identifiers: ["lambda.amazonaws.com"],
              type: "Service",
            },
          ],
        },
      ],
    });
    const iamForLambda = new IamRole(this, "iam_for_lambda", {
      assumeRolePolicy: Token.asString(assumeRole.json),
      name: "iam_for_lambda",
    });
    new LambdaFunction(this, "test_lambda", {
      environment: {
        variables: {
          foo: "bar",
        },
      },
      filename: "lambda_function_payload.zip",
      functionName: "lambda_function_name",
      handler: "index.test",
      role: iamForLambda.arn,
      runtime: "nodejs16.x",
      sourceCodeHash: Token.asString(lambda.outputBase64Sha256),
    });
  }
}

```

### Lambda Layers

~> **NOTE:** The `awsLambdaLayerVersion` attribute values for `arn` and `layerArn` were swapped in version 2.0.0 of the Terraform AWS Provider. For version 1.x, use `layerArn` references. For version 2.x, use `arn` references.

```typescript
// Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { LambdaFunction } from "./.gen/providers/aws/lambda-function";
import { LambdaLayerVersion } from "./.gen/providers/aws/lambda-layer-version";
interface MyConfig {
  layerName: any;
  functionName: any;
  role: any;
}
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string, config: MyConfig) {
    super(scope, name);
    const example = new LambdaLayerVersion(this, "example", {
      layerName: config.layerName,
    });
    const awsLambdaFunctionExample = new LambdaFunction(this, "example_1", {
      layers: [example.arn],
      functionName: config.functionName,
      role: config.role,
    });
    /*This allows the Terraform resource name to match the original name. You can remove the call if you don't need them to match.*/
    awsLambdaFunctionExample.overrideLogicalId("example");
  }
}

```

### Lambda Ephemeral Storage

Lambda Function Ephemeral Storage(`/tmp`) allows you to configure the storage upto `10` GB. The default value set to `512` MB.

```typescript
// Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { Token, TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { DataAwsIamPolicyDocument } from "./.gen/providers/aws/data-aws-iam-policy-document";
import { IamRole } from "./.gen/providers/aws/iam-role";
import { LambdaFunction } from "./.gen/providers/aws/lambda-function";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    const assumeRole = new DataAwsIamPolicyDocument(this, "assume_role", {
      statement: [
        {
          actions: ["sts:AssumeRole"],
          effect: "Allow",
          principals: [
            {
              identifiers: ["lambda.amazonaws.com"],
              type: "Service",
            },
          ],
        },
      ],
    });
    const iamForLambda = new IamRole(this, "iam_for_lambda", {
      assumeRolePolicy: Token.asString(assumeRole.json),
      name: "iam_for_lambda",
    });
    new LambdaFunction(this, "test_lambda", {
      ephemeralStorage: {
        size: 10240,
      },
      filename: "lambda_function_payload.zip",
      functionName: "lambda_function_name",
      handler: "index.test",
      role: iamForLambda.arn,
      runtime: "nodejs14.x",
    });
  }
}

```

### Lambda File Systems

Lambda File Systems allow you to connect an Amazon Elastic File System (EFS) file system to a Lambda function to share data across function invocations, access existing data including large files, and save function state.

```typescript
// Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { EfsAccessPoint } from "./.gen/providers/aws/efs-access-point";
import { EfsFileSystem } from "./.gen/providers/aws/efs-file-system";
import { EfsMountTarget } from "./.gen/providers/aws/efs-mount-target";
import { LambdaFunction } from "./.gen/providers/aws/lambda-function";
interface MyConfig {
  functionName: any;
  role: any;
}
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string, config: MyConfig) {
    super(scope, name);
    const efsForLambda = new EfsFileSystem(this, "efs_for_lambda", {
      tags: {
        Name: "efs_for_lambda",
      },
    });
    const alpha = new EfsMountTarget(this, "alpha", {
      fileSystemId: efsForLambda.id,
      securityGroups: [sgForLambda.id],
      subnetId: subnetForLambda.id,
    });
    const accessPointForLambda = new EfsAccessPoint(
      this,
      "access_point_for_lambda",
      {
        fileSystemId: efsForLambda.id,
        posixUser: {
          gid: 1000,
          uid: 1000,
        },
        rootDirectory: {
          creationInfo: {
            ownerGid: 1000,
            ownerUid: 1000,
            permissions: "777",
          },
          path: "/lambda",
        },
      }
    );
    new LambdaFunction(this, "example", {
      dependsOn: [alpha],
      fileSystemConfig: {
        arn: accessPointForLambda.arn,
        localMountPath: "/mnt/efs",
      },
      vpcConfig: {
        securityGroupIds: [sgForLambda.id],
        subnetIds: [subnetForLambda.id],
      },
      functionName: config.functionName,
      role: config.role,
    });
  }
}

```

### Lambda retries

Lambda Functions allow you to configure error handling for asynchronous invocation. The settings that it supports are `Maximum age of event` and `Retry attempts` as stated in [Lambda documentation for Configuring error handling for asynchronous invocation](https://docs.aws.amazon.com/lambda/latest/dg/invocation-async.html#invocation-async-errors). To configure these settings, refer to the [aws_lambda_function_event_invoke_config resource](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function_event_invoke_config).

## CloudWatch Logging and Permissions

For more information about CloudWatch Logs for Lambda, see the [Lambda User Guide](https://docs.aws.amazon.com/lambda/latest/dg/monitoring-functions-logs.html).

```typescript
// Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformVariable, Token, TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { CloudwatchLogGroup } from "./.gen/providers/aws/cloudwatch-log-group";
import { DataAwsIamPolicyDocument } from "./.gen/providers/aws/data-aws-iam-policy-document";
import { IamPolicy } from "./.gen/providers/aws/iam-policy";
import { IamRolePolicyAttachment } from "./.gen/providers/aws/iam-role-policy-attachment";
import { LambdaFunction } from "./.gen/providers/aws/lambda-function";
interface MyConfig {
  role: any;
}
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string, config: MyConfig) {
    super(scope, name);
    /*Terraform Variables are not always the best fit for getting inputs in the context of Terraform CDK.
    You can read more about this at https://cdk.tf/variables*/
    const lambdaFunctionName = new TerraformVariable(
      this,
      "lambda_function_name",
      {
        default: "lambda_function_name",
      }
    );
    const example = new CloudwatchLogGroup(this, "example", {
      name: "/aws/lambda/${" + lambdaFunctionName.value + "}",
      retentionInDays: 14,
    });
    const lambdaLogging = new DataAwsIamPolicyDocument(this, "lambda_logging", {
      statement: [
        {
          actions: [
            "logs:CreateLogGroup",
            "logs:CreateLogStream",
            "logs:PutLogEvents",
          ],
          effect: "Allow",
          resources: ["arn:aws:logs:*:*:*"],
        },
      ],
    });
    const awsIamPolicyLambdaLogging = new IamPolicy(this, "lambda_logging_3", {
      description: "IAM policy for logging from a lambda",
      name: "lambda_logging",
      path: "/",
      policy: Token.asString(lambdaLogging.json),
    });
    /*This allows the Terraform resource name to match the original name. You can remove the call if you don't need them to match.*/
    awsIamPolicyLambdaLogging.overrideLogicalId("lambda_logging");
    const lambdaLogs = new IamRolePolicyAttachment(this, "lambda_logs", {
      policyArn: Token.asString(awsIamPolicyLambdaLogging.arn),
      role: iamForLambda.name,
    });
    new LambdaFunction(this, "test_lambda", {
      dependsOn: [lambdaLogs, example],
      functionName: lambdaFunctionName.stringValue,
      role: config.role,
    });
  }
}

```

## Specifying the Deployment Package

AWS Lambda expects source code to be provided as a deployment package whose structure varies depending on which `runtime` is in use. See [Runtimes][6] for the valid values of `runtime`. The expected structure of the deployment package can be found in [the AWS Lambda documentation for each runtime][8].

Once you have created your deployment package you can specify it either directly as a local file (using the `filename` argument) or indirectly via Amazon S3 (using the `s3Bucket`, `s3Key` and `s3ObjectVersion` arguments). When providing the deployment package via S3 it may be useful to use [the `awsS3Object` resource](s3_object.html) to upload it.

For larger deployment packages it is recommended by Amazon to upload via S3, since the S3 API has better support for uploading large files efficiently.

## Argument Reference

The following arguments are required:

* `functionName` - (Required) Unique name for your Lambda Function.
* `role` - (Required) Amazon Resource Name (ARN) of the function's execution role. The role provides the function's identity and access to AWS services and resources.

The following arguments are optional:

* `architectures` - (Optional) Instruction set architecture for your Lambda function. Valid values are `["x8664"]` and `["arm64"]`. Default is `["x8664"]`. Removing this attribute, function's architecture stay the same.
* `codeSigningConfigArn` - (Optional) To enable code signing for this function, specify the ARN of a code-signing configuration. A code-signing configuration includes a set of signing profiles, which define the trusted publishers for this function.
* `deadLetterConfig` - (Optional) Configuration block. Detailed below.
* `description` - (Optional) Description of what your Lambda Function does.
* `environment` - (Optional) Configuration block. Detailed below.
* `ephemeralStorage` - (Optional) The amount of Ephemeral storage(`/tmp`) to allocate for the Lambda Function in MB. This parameter is used to expand the total amount of Ephemeral storage available, beyond the default amount of `512`MB. Detailed below.
* `fileSystemConfig` - (Optional) Configuration block. Detailed below.
* `filename` - (Optional) Path to the function's deployment package within the local filesystem. Exactly one of `filename`, `imageUri`, or `s3Bucket` must be specified.
* `handler` - (Optional) Function [entrypoint][3] in your code.
* `imageConfig` - (Optional) Configuration block. Detailed below.
* `imageUri` - (Optional) ECR image URI containing the function's deployment package. Exactly one of `filename`, `imageUri`,  or `s3Bucket` must be specified.
* `kmsKeyArn` - (Optional) Amazon Resource Name (ARN) of the AWS Key Management Service (KMS) key that is used to encrypt environment variables. If this configuration is not provided when environment variables are in use, AWS Lambda uses a default service key. If this configuration is provided when environment variables are not in use, the AWS Lambda API does not save this configuration and Terraform will show a perpetual difference of adding the key. To fix the perpetual difference, remove this configuration.
* `layers` - (Optional) List of Lambda Layer Version ARNs (maximum of 5) to attach to your Lambda Function. See [Lambda Layers][10]
* `memorySize` - (Optional) Amount of memory in MB your Lambda Function can use at runtime. Defaults to `128`. See [Limits][5]
* `packageType` - (Optional) Lambda deployment package type. Valid values are `zip` and `image`. Defaults to `zip`.
* `publish` - (Optional) Whether to publish creation/change as new Lambda Function Version. Defaults to `false`.
* `reservedConcurrentExecutions` - (Optional) Amount of reserved concurrent executions for this lambda function. A value of `0` disables lambda from being triggered and `1` removes any concurrency limitations. Defaults to Unreserved Concurrency Limits `1`. See [Managing Concurrency][9]
* `replaceSecurityGroupsOnDestroy` - (Optional, **Deprecated**) **AWS no longer supports this operation. This attribute now has no effect and will be removed in a future major version.** Whether to replace the security groups on associated lambda network interfaces upon destruction. Removing these security groups from orphaned network interfaces can speed up security group deletion times by avoiding a dependency on AWS's internal cleanup operations. By default, the ENI security groups will be replaced with the `default` security group in the function's VPC. Set the `replacementSecurityGroupIds` attribute to use a custom list of security groups for replacement.
* `replacementSecurityGroupIds` - (Optional, **Deprecated**) List of security group IDs to assign to orphaned Lambda function network interfaces upon destruction. `replaceSecurityGroupsOnDestroy` must be set to `true` to use this attribute.
* `runtime` - (Optional) Identifier of the function's runtime. See [Runtimes][6] for valid values.
* `s3Bucket` - (Optional) S3 bucket location containing the function's deployment package. This bucket must reside in the same AWS region where you are creating the Lambda function. Exactly one of `filename`, `imageUri`, or `s3Bucket` must be specified. When `s3Bucket` is set, `s3Key` is required.
* `s3Key` - (Optional) S3 key of an object containing the function's deployment package. When `s3Bucket` is set, `s3Key` is required.
* `s3ObjectVersion` - (Optional) Object version containing the function's deployment package. Conflicts with `filename` and `imageUri`.
* `skipDestroy` - (Optional) Set to true if you do not wish the function to be deleted at destroy time, and instead just remove the function from the Terraform state.
* `sourceCodeHash` - (Optional) Used to trigger updates. Must be set to a base64-encoded SHA256 hash of the package file specified with either `filename` or `s3Key`. The usual way to set this is `filebase64Sha256("fileZip")` (Terraform 0.11.12 and later) or `base64Sha256(file("fileZip"))` (Terraform 0.11.11 and earlier), where "file.zip" is the local filename of the lambda function source archive.
* `snapStart` - (Optional) Snap start settings block. Detailed below.
* `tags` - (Optional) Map of tags to assign to the object. If configured with a provider [`defaultTags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.
* `timeout` - (Optional) Amount of time your Lambda Function has to run in seconds. Defaults to `3`. See [Limits][5].
* `tracingConfig` - (Optional) Configuration block. Detailed below.
* `vpcConfig` - (Optional) Configuration block. Detailed below.

### dead_letter_config

Dead letter queue configuration that specifies the queue or topic where Lambda sends asynchronous events when they fail processing. For more information, see [Dead Letter Queues](https://docs.aws.amazon.com/lambda/latest/dg/invocation-async.html#dlq).

* `targetArn` - (Required) ARN of an SNS topic or SQS queue to notify when an invocation fails. If this option is used, the function's IAM role must be granted suitable access to write to the target object, which means allowing either the `sns:publish` or `sqs:sendMessage` action on this ARN, depending on which service is targeted.

### environment

* `variables` - (Optional) Map of environment variables that are accessible from the function code during execution. If provided at least one key must be present.

### ephemeral_storage

* `size` - (Required) The size of the Lambda function Ephemeral storage(`/tmp`) represented in MB. The minimum supported `ephemeralStorage` value defaults to `512`MB and the maximum supported value is `10240`MB.

### file_system_config

Connection settings for an EFS file system. Before creating or updating Lambda functions with `fileSystemConfig`, EFS mount targets must be in available lifecycle state. Use `dependsOn` to explicitly declare this dependency. See [Using Amazon EFS with Lambda][12].

* `arn` - (Required) Amazon Resource Name (ARN) of the Amazon EFS Access Point that provides access to the file system.
* `localMountPath` - (Required) Path where the function can access the file system, starting with /mnt/.

### image_config

Container image configuration values that override the values in the container image Dockerfile.

* `command` - (Optional) Parameters that you want to pass in with `entryPoint`.
* `entryPoint` - (Optional) Entry point to your application, which is typically the location of the runtime executable.
* `workingDirectory` - (Optional) Working directory.

### snap_start

Snap start settings for low-latency startups. This feature is currently only supported for `java11` runtimes. Remove this block to delete the associated settings (rather than setting `apply_on = "None"`).

* `applyOn` - (Required) Conditions where snap start is enabled. Valid values are `publishedVersions`.

### tracing_config

* `mode` - (Required) Whether to sample and trace a subset of incoming requests with AWS X-Ray. Valid values are `passThrough` and `active`. If `passThrough`, Lambda will only trace the request from an upstream service if it contains a tracing header with "sampled=1". If `active`, Lambda will respect any tracing header it receives from an upstream service. If no tracing header is received, Lambda will call X-Ray for a tracing decision.

### vpc_config

For network connectivity to AWS resources in a VPC, specify a list of security groups and subnets in the VPC. When you connect a function to a VPC, it can only access resources and the internet through that VPC. See [VPC Settings][7].

~> **NOTE:** If both `subnetIds` and `securityGroupIds` are empty then `vpcConfig` is considered to be empty or unset.

* `securityGroupIds` - (Required) List of security group IDs associated with the Lambda function.
* `subnetIds` - (Required) List of subnet IDs associated with the Lambda function.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) identifying your Lambda Function.
* `invokeArn` - ARN to be used for invoking Lambda Function from API Gateway - to be used in [`awsApiGatewayIntegration`](/docs/providers/aws/r/api_gateway_integration.html)'s `uri`.
* `lastModified` - Date this resource was last modified.
* `qualifiedArn` - ARN identifying your Lambda Function Version (if versioning is enabled via `publish = true`).
* `qualifiedInvokeArn` - Qualified ARN (ARN with lambda version number) to be used for invoking Lambda Function from API Gateway - to be used in [`awsApiGatewayIntegration`](/docs/providers/aws/r/api_gateway_integration.html)'s `uri`.
* `signingJobArn` - ARN of the signing job.
* `signingProfileVersionArn` - ARN of the signing profile version.
* `snapStartOptimizationStatus` - Optimization status of the snap start configuration. Valid values are `on` and `off`.
* `sourceCodeSize` - Size in bytes of the function .zip file.
* `tagsAll` - A map of tags assigned to the resource, including those inherited from the provider [`defaultTags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block).
* `version` - Latest published version of your Lambda Function.
* `vpcConfigVpcId` - ID of the VPC.

[1]: https://docs.aws.amazon.com/lambda/latest/dg/welcome.html
[3]: https://docs.aws.amazon.com/lambda/latest/dg/walkthrough-custom-events-create-test-function.html
[4]: https://docs.aws.amazon.com/lambda/latest/dg/intro-permission-model.html
[5]: https://docs.aws.amazon.com/lambda/latest/dg/limits.html
[6]: https://docs.aws.amazon.com/lambda/latest/dg/API_CreateFunction.html#SSS-CreateFunction-request-Runtime
[7]: https://docs.aws.amazon.com/lambda/latest/dg/configuration-vpc.html
[8]: https://docs.aws.amazon.com/lambda/latest/dg/deployment-package-v2.html
[9]: https://docs.aws.amazon.com/lambda/latest/dg/concurrent-executions.html
[10]: https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html
[11]: https://learn.hashicorp.com/terraform/aws/lambda-api-gateway
[12]: https://docs.aws.amazon.com/lambda/latest/dg/services-efs.html

## Timeouts

[Configuration options](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts):

* `create` - (Default `10M`)
* `update` - (Default `10M`)
* `delete` - (Default `10M`)

## Import

Lambda Functions can be imported using the `functionName`, e.g.,

```
$ terraform import aws_lambda_function.test_lambda my_test_lambda_function
```

<!-- cache-key: cdktf-0.17.1 input-f83d1fe4cff176090ce1a1c27506e7b709b89d4b703aee61a87c906d4d71c81e -->