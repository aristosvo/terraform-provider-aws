// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identitystore_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccIdentityStoreGroupsDataSource_filterDisplayName(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName2 := "aws_identitystore_group.test.1"
	dataSourceName := "data.aws_identitystore_groups.test"
	name := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckSSOAdminInstances(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.IdentityStoreServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccGroupsDataSourceConfig_filterDisplayName(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "groups.1.description", resourceName2, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "groups.1.display_name", resourceName2, "display_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "groups.1.group_id", resourceName2, "group_id"),
					resource.TestCheckResourceAttr(dataSourceName, "groups.1.external_ids.#", "0"),
				),
			},
		},
	})
}

func testAccGroupsDataSourceConfig_base(name string) string {
	return fmt.Sprintf(`
data "aws_ssoadmin_instances" "test" {}

resource "aws_identitystore_group" "test" {
  count = 2
  identity_store_id = tolist(data.aws_ssoadmin_instances.test.identity_store_ids)[0]
  display_name      = "%[1]q-${count.index}"
  description       = "Acceptance Test ${count.index}"
}

`, name)
}

func testAccGroupsDataSourceConfig_filterDisplayName(name string) string {
	return acctest.ConfigCompose(testAccGroupsDataSourceConfig_base(name), `
data "aws_identitystore_group" "test" {
  filter {
    attribute_path  = "DisplayName"
    attribute_value = aws_identitystore_group.test.1.display_name
  }

  identity_store_id = tolist(data.aws_ssoadmin_instances.test.identity_store_ids)[0]
}
`)
}
