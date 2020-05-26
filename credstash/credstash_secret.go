package credstash

import (
	"github.com/Versent/unicreds"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceCredstashSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecretPut,
		Read:   resourceSecretRead,
		Delete: resourceSecretDelete,
		Exists: resourceSecretExists,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the secret",
			},
			"value": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Required:    true,
				ForceNew:    true,
				Description: "Value of the secret",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the secrets",
			},
			"context": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Encryption context for the secret",
			},
		},
	}
}

func resourceSecretExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	config := meta.(*Config)
	name := d.Get("name").(string)

	log.Printf("[DEBUG] Checking secret name=%q", name)
	_, err := unicreds.GetHighestVersion(&config.TableName, name)
	if err != nil {
		log.Printf("[DEBUG] Error checking secret: %s", err.Error())
		if err == unicreds.ErrSecretNotFound {
			log.Print("[DEBUG] Matched NotFound error, returning no error")
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func resourceSecretPut(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	name := d.Get("name").(string)
	secretData := d.Get("value").(string)

	context := getContext(d)

	version, err := unicreds.ResolveVersion(&config.TableName, name, 0)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Create secret name=%q table=%q version=%q context=%+v", name, config.TableName, version, context)
	err = unicreds.PutSecret(&config.TableName, config.KmsKey, name, secretData, version, context)
	if err != nil {
		return err
	}

	d.SetId(name + "_" + version)
	d.Set("version", version)
	d.Set("context", context)

	return nil
}

func resourceSecretRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	name := d.Get("name").(string)
	version := d.Get("version").(string)

	context := getContext(d)

	if version == "" {
		v, err := unicreds.GetHighestVersion(&config.TableName, name)
		if err != nil {
			return err
		}
		version = v
	}

	log.Printf("[DEBUG] Getting secret for name=%q version=%q context=%+v", name, version, context)
	secretData, err := unicreds.GetSecret(&config.TableName, name, version, context)
	if err != nil {
		return err
	}

	d.Set("value", secretData.Secret)
	d.Set("name", name)
	d.Set("version", version)
	d.Set("context", context)

	return nil
}

func resourceSecretDelete(d *schema.ResourceData, meta interface{}) error {

	// We don't want to delete any secrets so we just remove it from tf state.

	//config := meta.(*Config)
	//name := d.Get("name").(string)
	//err := unicreds.DeleteSecret(&config.TableName, name)
	//
	//if err != nil {
	//	return err
	//}
	d.SetId("")
	return nil
}
