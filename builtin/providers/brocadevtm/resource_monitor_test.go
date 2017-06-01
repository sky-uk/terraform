package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm"
	"github.com/sky-uk/go-brocade-vtm/api/monitor"
	"testing"
)

func TestAccBrocadeVTMMonitorBasic(t *testing.T) {

	randomInt := acctest.RandInt()
	monitorName := fmt.Sprintf("acctest_brocadevtm_monitor-%d", randomInt)
	monitorResourceName := "brocadevtm_monitor.acctest"

	fmt.Printf("\n\nMonitor Name is %s.\n\n", monitorName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccBrocadeVTMMonitorCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBrocadeVTMMonitorBasicTemplate(monitorName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMMonitorExists(monitorName, monitorResourceName),
					resource.TestCheckResourceAttr(monitorResourceName, "name", monitorName),
					resource.TestCheckResourceAttr(monitorResourceName, "delay", "6"),
					resource.TestCheckResourceAttr(monitorResourceName, "timeout", "3"),
					resource.TestCheckResourceAttr(monitorResourceName, "failures", "3"),
					resource.TestCheckResourceAttr(monitorResourceName, "verbose", "true"),
					resource.TestCheckResourceAttr(monitorResourceName, "use_ssl", "false"),
					resource.TestCheckResourceAttr(monitorResourceName, "http_host_header", ""),
					resource.TestCheckResourceAttr(monitorResourceName, "http_authentication", ""),
					resource.TestCheckResourceAttr(monitorResourceName, "http_body_regex", ""),
					resource.TestCheckResourceAttr(monitorResourceName, "http_path", "/"),
				),
			},
			{
				Config: testAccBrocadeVTMMonitorAllTemplate(monitorName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMMonitorExists(monitorName, monitorResourceName),
					resource.TestCheckResourceAttr(monitorResourceName, "name", monitorName),
					resource.TestCheckResourceAttr(monitorResourceName, "delay", "5"),
					resource.TestCheckResourceAttr(monitorResourceName, "timeout", "3"),
					resource.TestCheckResourceAttr(monitorResourceName, "failures", "9"),
					resource.TestCheckResourceAttr(monitorResourceName, "verbose", "false"),
					resource.TestCheckResourceAttr(monitorResourceName, "use_ssl", "false"),
					resource.TestCheckResourceAttr(monitorResourceName, "http_host_header", "some_header"),
					resource.TestCheckResourceAttr(monitorResourceName, "http_authentication", "some_authentication"),
					resource.TestCheckResourceAttr(monitorResourceName, "http_body_regex", "^ok"),
					resource.TestCheckResourceAttr(monitorResourceName, "http_path", "/some/status/page"),
				),
			},
		},
	})
}

func testAccBrocadeVTMMonitorCheckDestroy(state *terraform.State) error {

	vtmClient := testAccProvider.Meta().(*brocadevtm.VTMClient)
	var name string

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_monitor" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}

		api := monitor.NewGetAll()
		err := vtmClient.Do(api)
		if err != nil {
			return nil
		}
		if api.GetResponse().FilterByName(name).Name == "acctest_brocade_monitor" {
			return fmt.Errorf("Brocade vTM monitor %s still exists.", name)
		}
	}
	return nil
}

func testAccBrocadeVTMMonitorExists(monitorName, monitorResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[monitorResourceName]
		if !ok {
			return fmt.Errorf("\nBrocade vTM Monitor resource %s not found in resources\n", monitorResourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\nBrocade vTM Monitor ID not set in resources\n")
		}

		vtmClient := testAccProvider.Meta().(*brocadevtm.VTMClient)
		getAllAPI := monitor.NewGetAll()

		err := vtmClient.Do(getAllAPI)
		if err != nil {
			return fmt.Errorf("Error: %+v", err)
		}
		foundMonitor := getAllAPI.GetResponse().FilterByName(monitorName)
		if foundMonitor.Name != monitorName {
			return fmt.Errorf("Brocade vTM Monitor %s not found on remote vTM", monitorName)
		}
		return nil
	}
}

func testAccBrocadeVTMMonitorBasicTemplate(monitorName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_monitor" "acctest" {
  name = "%s"
  delay = 6
  verbose = true
}
`, monitorName)
}

func testAccBrocadeVTMMonitorAllTemplate(monitorName string) string {
	return fmt.Sprintf(`
resource "brocadevtm_monitor" "acctest" {
  name = "%s"
  delay = 5
  timeout = 3
  failures = 9
  verbose = false
  use_ssl = false
  http_host_header = "some_header"
  http_authentication = "some_authentication"
  http_body_regex = "^ok"
  http_path = "/some/status/page"
}
`, monitorName)
}
