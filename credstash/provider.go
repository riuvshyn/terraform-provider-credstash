package credstash

import (
	"github.com/Versent/unicreds"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var _ terraform.ResourceProvider = Provider()

const defaultKMSKey = "alias/credstash"

type Config struct {
	Region    string
	TableName string
	Profile   string
	KmsKey    string
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"credstash_secret": dataSourceSecret(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"credstash_secret": resourceCredstashSecret(),
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"AWS_REGION",
					"AWS_DEFAULT_REGION",
				}, nil),
				Description:  "The region where AWS operations will take place. Examples are us-east-1, us-west-2, etc.",
				InputDefault: "eu-central-1",
			},
			"table": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The DynamoDB table where the secrets are stored.",
				Default:     "credential-store",
			},
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The profile that should be used to connect to AWS",
			},
			"kms_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     defaultKMSKey,
				Description: "The KMS key to use when storing secrets",
			},
		},
		ConfigureFunc: providerConfig,
	}
}

func providerConfig(d *schema.ResourceData) (interface{}, error) {
	region := d.Get("region").(string)
	profile := d.Get("profile").(string)

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile,
		Config: aws.Config{
			Region: aws.String(region),
		},
	})
	if err != nil {
		return nil, err
	}

	unicreds.SetDynamoDBSession(sess)
	unicreds.SetKMSSession(sess)

	return &Config{
		TableName: d.Get("table").(string),
		KmsKey:    d.Get("kms_key").(string),
		Region:    region,
		Profile:   profile,
	}, nil
}
