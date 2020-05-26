package credstash

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccCredstashSecret(t *testing.T) {
	resourceName := "credstash_secret.terraform"
	resourceNameNewVersion := "credstash_secret.terraform_new_version"
	dsResourceName := "data.credstash_secret.terraform"
	dsResourceNameNewVersion := "data.credstash_secret.terraform_new_version"
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configData(testAccCheckCredstashSecretBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "name", dsResourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "value", dsResourceName, "value"),
					resource.TestCheckResourceAttrPair(resourceName, "version", dsResourceName, "version"),
				),
			},
			{
				Config: configData(testAccCheckCredstashSecretResourceNewVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameNewVersion, "name", "terraform-resource-acc1"),
					resource.TestCheckResourceAttr(resourceNameNewVersion, "value", "test1-updated"),
					resource.TestCheckResourceAttr(resourceNameNewVersion, "version", "0000000000000000002"),
				),
			},
			{
				Config: configData(testAccCheckCredstashSecretDatasourceVersions),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameNewVersion, "name", dsResourceNameNewVersion, "name"),
					resource.TestCheckResourceAttrPair(resourceNameNewVersion, "value", dsResourceNameNewVersion, "value"),
					resource.TestCheckResourceAttrPair(resourceNameNewVersion, "version", dsResourceNameNewVersion, "version"),
				),
			},
		},
	})
}

func TestAccCredstashSecretData(t *testing.T) {
	resourceName := "credstash_secret.resource_test1"
	dsResourceName := "data.credstash_secret.data_test1"
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configData(testAccCheckCredstashSecretBasicDatasource),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "name", dsResourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "value", dsResourceName, "value"),
					resource.TestCheckResourceAttrPair(resourceName, "version", dsResourceName, "version"),
				),
			},
		},
	})
}

func TestAccCredstashSecretDataWithContext(t *testing.T) {
	resourceName := "credstash_secret.resource_test2"
	dsResourceName := "data.credstash_secret.data_test2"
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configData(testAccCheckCredstashSecretBasicDatasourceWithContext),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "name", dsResourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "value", dsResourceName, "value"),
					resource.TestCheckResourceAttrPair(resourceName, "version", dsResourceName, "version"),
					resource.TestCheckResourceAttrPair(resourceName, "context.foo", dsResourceName, "context.foo"),
					resource.TestCheckResourceAttrPair(resourceName, "context.baz", dsResourceName, "context.baz"),
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
data "credstash_secret" "terraform" {
	name    = "${credstash_secret.terraform.name}"
}
`

const testAccCheckCredstashSecretResourceNewVersion = `
resource "credstash_secret" "terraform_new_version" {
	name = "terraform-resource-acc1"
	value = "test1-updated"
}
`

const testAccCheckCredstashSecretDatasourceVersions = `
resource "credstash_secret" "terraform_new_version" {
	name = "terraform-resource-acc1"
	value = "test1-updated"
}
data "credstash_secret" "terraform_new_version" {
	name    = "${credstash_secret.terraform_new_version.name}"
	version = "2"
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
const providerConfiguration = `
provider "credstash" {
  region = "%v"
  profile = "%v"
  table = "%v"
  kms_key = "%v"
}
`

func configData(data string) string {
	region := "eu-central-1"
	profile := "staging"
	dynamoDBTable := "credential-store"
	kms_key := "alias/credstash"
	if val, ok := os.LookupEnv("AWS_REGION"); ok {
		region = val
	}
	if val, ok := os.LookupEnv("AWS_PROFILE"); ok {
		profile = val
	}
	if val, ok := os.LookupEnv("CREDSTASH_DYNAMODB_TABLE"); ok {
		dynamoDBTable = val
	}
	if val, ok := os.LookupEnv("AWS_KMS_KEY"); ok {
		kms_key = val
	}
	provider := fmt.Sprintf(providerConfiguration, region, profile, dynamoDBTable, kms_key)

	return fmt.Sprint(provider, data)
}
