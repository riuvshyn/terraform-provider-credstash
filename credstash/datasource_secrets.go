package credstash

import (
	"log"
	"strconv"

	"github.com/Versent/unicreds"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceSecret() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecretRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the secret",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Version of the secrets",
				Default:     "",
			},
			"table": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of DynamoDB table where the secrets are stored",
				Default:     "",
			},
			"context": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Encryption context for the secret",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Value of the secret",
			},
		},
	}
}

func dataSourceSecretRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	name := d.Get("name").(string)
	version := d.Get("version").(string)

	context := getContext(d)

	// get latest version if version is not set
	if version == "" {
		v, err := unicreds.GetHighestVersion(&config.TableName, name)
		if err != nil {
			return err
		}
		version = v
	} else {
		v, err := strconv.Atoi(version)
		if err != nil {
			return err
		}
		version = unicreds.PaddedInt(v)
	}

	log.Printf("[DEBUG] Getting secret for name=%q table=%q version=%q context=%+v", name, config.TableName, version, context)
	out, err := unicreds.GetSecret(&config.TableName, name, version, context)
	if err != nil {
		return err
	}

	d.Set("value", out.Secret)
	d.Set("version", out.Version)
	d.Set("context", context)
	d.SetId(name)

	return nil
}
