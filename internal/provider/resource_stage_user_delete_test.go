package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStageUserDelete(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageUserDelete(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_user_delete.name", "name", rName),
				),
			},
		},
	})
}

func testAccResourceStageUserDelete(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_user_delete" "name" {
  name              = "%s"
}
`, name)
}