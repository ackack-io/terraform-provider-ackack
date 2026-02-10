// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMonitorResource_HTTP(t *testing.T) {
	rName := acctest.RandomWithPrefix("tfacc")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMonitorResourceConfig_HTTP(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ackack_monitor.test", "name", rName),
					resource.TestCheckResourceAttr("ackack_monitor.test", "type", "http"),
					resource.TestCheckResourceAttr("ackack_monitor.test", "url", "https://example.com"),
					resource.TestCheckResourceAttr("ackack_monitor.test", "is_enabled", "true"),
					resource.TestCheckResourceAttrSet("ackack_monitor.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "ackack_monitor.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status", "last_checked", "updated_at", "uptime_percentage"},
			},
			// Update and Read testing
			{
				Config: testAccMonitorResourceConfig_HTTPUpdated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ackack_monitor.test", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("ackack_monitor.test", "frequency_seconds", "120"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMonitorResourceConfig_HTTP(name string) string {
	return fmt.Sprintf(`
resource "ackack_monitor" "test" {
  name              = %[1]q
  type              = "http"
  url               = "https://example.com"
  frequency_seconds = 60
  timeout_ms        = 10000
  is_enabled        = true
}
`, name)
}

func testAccMonitorResourceConfig_HTTPUpdated(name string) string {
	return fmt.Sprintf(`
resource "ackack_monitor" "test" {
  name              = "%[1]s-updated"
  type              = "http"
  url               = "https://example.com"
  frequency_seconds = 120
  timeout_ms        = 10000
  is_enabled        = true
}
`, name)
}
