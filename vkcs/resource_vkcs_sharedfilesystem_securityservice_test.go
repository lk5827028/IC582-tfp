package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/securityservices"
)

func TestAccSFSSecurityService_basic(t *testing.T) {
	var securityservice securityservices.SecurityService

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckSFS(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSSecurityServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSSecurityServiceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSSecurityServiceExists("vkcs_sharedfilesystem_securityservice.securityservice_1", &securityservice),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_securityservice.securityservice_1", "name", "security"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_securityservice.securityservice_1", "description", "created by terraform"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_securityservice.securityservice_1", "type", "active_directory"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_securityservice.securityservice_1", "server", "192.168.199.10"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_securityservice.securityservice_1", "dns_ip", "192.168.199.10"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_securityservice.securityservice_1", "domain", "example.com"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_securityservice.securityservice_1", "ou", "CN=Computers,DC=example,DC=com"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_securityservice.securityservice_1", "user", "joinDomainUser"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_securityservice.securityservice_1", "password", "s8cret"),
				),
			},
			{
				Config: testAccSFSSecurityServiceConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSSecurityServiceExists("vkcs_sharedfilesystem_securityservice.securityservice_1", &securityservice),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_securityservice.securityservice_1", "name", "security_through_obscurity"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_securityservice.securityservice_1", "description", ""),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_securityservice.securityservice_1", "type", "kerberos"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_securityservice.securityservice_1", "server", "192.168.199.11"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_securityservice.securityservice_1", "dns_ip", "192.168.199.11"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_securityservice.securityservice_1", "domain", ""),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_securityservice.securityservice_1", "ou", ""),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_securityservice.securityservice_1", "user", ""),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_securityservice.securityservice_1", "password", ""),
				),
			},
		},
	})
}

func testAccCheckSFSSecurityServiceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config)
	sfsClient, err := config.SharedfilesystemV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_sharedfilesystem_securityservice" {
			continue
		}

		_, err := securityservices.Get(sfsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Manila securityservice still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckSFSSecurityServiceExists(n string, securityservice *securityservices.SecurityService) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config)
		sfsClient, err := config.SharedfilesystemV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
		}

		found, err := securityservices.Get(sfsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*securityservice = *found

		return nil
	}
}

const testAccSFSSecurityServiceConfigBasic = `
resource "vkcs_sharedfilesystem_securityservice" "securityservice_1" {
  name        = "security"
  description = "created by terraform"
  type        = "active_directory"
  server      = "192.168.199.10"
  dns_ip      = "192.168.199.10"
  domain      = "example.com"
  ou          = "CN=Computers,DC=example,DC=com"
  user        = "joinDomainUser"
  password    = "s8cret"
}
`

const testAccSFSSecurityServiceConfigUpdate = `
resource "vkcs_sharedfilesystem_securityservice" "securityservice_1" {
  name        = "security_through_obscurity"
  description = ""
  type        = "kerberos"
  server      = "192.168.199.11"
  dns_ip      = "192.168.199.11"
}
`
