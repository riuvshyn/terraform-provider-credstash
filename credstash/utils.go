package credstash

import (
	"fmt"
	"github.com/Versent/unicreds"
	"github.com/hashicorp/terraform/helper/schema"
)

func getContext(d *schema.ResourceData) *unicreds.EncryptionContextValue {
	context := unicreds.NewEncryptionContextValue()
	for k, v := range d.Get("context").(map[string]interface{}) {
		context.Set(fmt.Sprintf("%s:%v", k, v))
	}
	return context
}
