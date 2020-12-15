package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func resourceMongoDBAtlasEncryptionAtRest() *schema.Resource {
	return &schema.Resource{
		Create:   resourceMongoDBAtlasEncryptionAtRestCreate,
		Read:     resourceMongoDBAtlasEncryptionAtRestRead,
		Delete:   resourceMongoDBAtlasEncryptionAtRestDelete,
		Update:   resourceMongoDBAtlasEncryptionAtRestUpdate,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"aws_kms": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"access_key_id": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"secret_access_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"customer_master_key_id": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"azure_key_vault": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"client_id": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"azure_environment": {
							Type:     schema.TypeString,
							Required: true,
						},
						"subscription_id": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"resource_group_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key_vault_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key_identifier": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"secret": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"tenant_id": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
			"google_cloud_kms": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"service_account_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"key_version_resource_id": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

func resourceMongoDBAtlasEncryptionAtRestCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	encryptionAtRestReq := &matlas.EncryptionAtRest{
		GroupID:        d.Get("project_id").(string),
		AwsKms:         expandAwsKms(d.Get("aws_kms").(map[string]interface{})),
		AzureKeyVault:  expandAzureKeyVault(d.Get("azure_key_vault").(map[string]interface{})),
		GoogleCloudKms: expandGCPKms(d.Get("google_cloud_kms").(map[string]interface{})),
	}

	_, _, err := conn.EncryptionsAtRest.Create(context.Background(), encryptionAtRestReq)
	if err != nil {
		return fmt.Errorf("error creating Encryption at Rest: %s", err)
	}

	d.SetId(d.Get("project_id").(string))

	return resourceMongoDBAtlasEncryptionAtRestRead(d, meta)
}

func resourceMongoDBAtlasEncryptionAtRestRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	_, _, err := conn.EncryptionsAtRest.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error getting Encryption at Rest information: %s", err)
	}

	return nil
}

func resourceMongoDBAtlasEncryptionAtRestUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	projectID := d.Id()

	encrypt, _, err := conn.EncryptionsAtRest.Get(context.Background(), projectID)
	if err != nil {
		return fmt.Errorf("error getting encryption at rest information: %s", err)
	}

	encrypt.GroupID = projectID

	if d.HasChange("aws_kms") {
		encrypt.AwsKms = expandAwsKms(d.Get("aws_kms").(map[string]interface{}))
	}

	if d.HasChange("azure_key_vault") {
		encrypt.AzureKeyVault = expandAzureKeyVault(d.Get("azure_key_vault").(map[string]interface{}))
	}

	if d.HasChange("google_cloud_kms") {
		encrypt.GoogleCloudKms = expandGCPKms(d.Get("google_cloud_kms").(map[string]interface{}))
	}

	_, _, err = conn.EncryptionsAtRest.Create(context.Background(), encrypt)
	if err != nil {
		return fmt.Errorf("error updating encryption at rest (%s): %s", projectID, err)
	}

	return resourceMongoDBAtlasEncryptionAtRestRead(d, meta)
}

func resourceMongoDBAtlasEncryptionAtRestDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	_, err := conn.EncryptionsAtRest.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error removing encryption at rest (%s): %s", d.Id(), err)
	}

	return nil
}

func expandAwsKms(awsKms map[string]interface{}) matlas.AwsKms {
	awsRegion, _ := valRegion(awsKms["region"])

	return matlas.AwsKms{
		Enabled:             pointy.Bool(cast.ToBool(awsKms["enabled"])),
		AccessKeyID:         cast.ToString(awsKms["access_key_id"]),
		SecretAccessKey:     cast.ToString(awsKms["secret_access_key"]),
		CustomerMasterKeyID: cast.ToString(awsKms["customer_master_key_id"]),
		Region:              awsRegion,
		RoleID:              cast.ToString(awsKms["role_id"]),
	}
}

func expandAzureKeyVault(azure map[string]interface{}) matlas.AzureKeyVault {
	return matlas.AzureKeyVault{
		Enabled:           pointy.Bool(cast.ToBool(azure["enabled"])),
		ClientID:          cast.ToString(azure["client_id"]),
		AzureEnvironment:  cast.ToString(azure["azure_environment"]),
		SubscriptionID:    cast.ToString(azure["subscription_id"]),
		ResourceGroupName: cast.ToString(azure["resource_group_name"]),
		KeyVaultName:      cast.ToString(azure["key_vault_name"]),
		KeyIdentifier:     cast.ToString(azure["key_identifier"]),
		Secret:            cast.ToString(azure["secret"]),
		TenantID:          cast.ToString(azure["tenant_id"]),
	}
}

func expandGCPKms(gcpKms map[string]interface{}) matlas.GoogleCloudKms {
	return matlas.GoogleCloudKms{
		Enabled:              pointy.Bool(cast.ToBool(gcpKms["enabled"])),
		ServiceAccountKey:    cast.ToString(gcpKms["service_account_key"]),
		KeyVersionResourceID: cast.ToString(gcpKms["key_version_resource_id"]),
	}
}
