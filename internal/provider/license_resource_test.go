// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAidboxLicenseResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAidboxLicenseResourceConfig("license-one", "development"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_license.test", "name", "license-one"),
					resource.TestCheckResourceAttr("aidbox_license.test", "type", "development"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aidbox_license.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccAidboxLicenseResourceConfig("license-two", "development"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_license.test", "name", "license-two"),
					resource.TestCheckResourceAttr("aidbox_license.test", "type", "development"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAidboxLicenseResourceConfig(name string, licenseType string) string {
	return fmt.Sprintf(`
resource "aidbox_license" "test" {
  name = %[1]q
  type = %[2]q
}
`, name, licenseType)
}
