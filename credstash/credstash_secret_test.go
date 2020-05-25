package credstash

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccCredstashSecret(t *testing.T) {
	resourceName := "credstash_secret.terraform"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCredstashSecretBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "terraform-resource-acc1"),
					resource.TestCheckResourceAttr(resourceName, "value", "test1"),
					resource.TestCheckResourceAttr(resourceName, "version", "0000000000000000001"),
				),
			},
			{
				Config: testAccCheckCredstashSecretBasicNewVersion,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "terraform-resource-acc1"),
					resource.TestCheckResourceAttr(resourceName, "value", "test1-updated"),
					resource.TestCheckResourceAttr(resourceName, "version", "0000000000000000002"),
				),
			},
		},
	})
}

func TestAccCredstashSecretData(t *testing.T)  {
	resourceName := "credstash_secret.resource_test1"
	dsResourceName := "data.credstash_secret.data_test1"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCredstashSecretBasicDatasource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "terraform-data-acc1"),
					resource.TestCheckResourceAttr(resourceName, "value", "secret1"),
					resource.TestCheckResourceAttr(dsResourceName, "value", "secret1"),
				),
			},
		},
	})
}

func TestAccCredstashSecretDataWithContext(t *testing.T)  {
	resourceName := "credstash_secret.resource_test2"
	dsResourceName := "data.credstash_secret.data_test2"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCredstashSecretBasicDatasourceWithContext,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "terraform-data-acc2"),
					resource.TestCheckResourceAttr(resourceName, "value", "secret2"),
					resource.TestCheckResourceAttr(dsResourceName, "value", "secret2"),
					resource.TestCheckResourceAttrPair(resourceName, "context.foo", dsResourceName,"context.foo"),
					resource.TestCheckResourceAttrPair(resourceName, "context.baz", dsResourceName,"context.baz"),
				),
			},
		},
	})
}

const testAccCheckCredstashSecretBasic = `
resource "credstash_secret" "terraform" {
	name = "terraform-resource-acc1"
	value = "test1"
}
`
const testAccCheckCredstashSecretBasicNewVersion = `
resource "credstash_secret" "terraform" {
	name = "terraform-resource-acc1"
	value = "test1-updated"
}
`

const testAccCheckCredstashSecretBasicDatasource = `
resource "credstash_secret" "resource_test1" {
	name = "terraform-data-acc1"
	value = "secret1"
}
data "credstash_secret" "data_test1" {
	name    = "${credstash_secret.resource_test1.name}"
}
`

const testAccCheckCredstashSecretBasicDatasourceWithContext = `
resource "credstash_secret" "resource_test2" {
	name = "terraform-data-acc2"
	value = "secret2"
	context = {
		"foo" = "bar"
		"baz" = "qux"
	 }
}
data "credstash_secret" "data_test2" {
	name    = "${credstash_secret.resource_test2.name}"
	context = {
		"foo" = "bar"
		"baz" = "qux"
	 }
}
`
