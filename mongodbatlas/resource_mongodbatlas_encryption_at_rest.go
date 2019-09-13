package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
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
							Type:     schema.TypeString,
							Required: true,
						},
						"secret_access_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"customer_master_key_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"region": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"azure_key_vault": {
				Type:     schema.TypeMap,
				ForceNew: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							ForceNew: true,
							Required: true,
						},
						"client_id": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"azure_environment": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"subscription_id": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"resource_group_name": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"key_vault_name": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"key_identifier": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"secret": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
					},
				},
			},
			"google_cloud_kms": {
				Type:     schema.TypeMap,
				ForceNew: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							ForceNew: true,
							Required: true,
						},
						"service_account_key": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"key_version_resource_id": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func expandAwsKms(awsKms map[string]interface{}) matlas.AwsKms {
	awsRegion, _ := valRegion(awsKms["region"])
	return matlas.AwsKms{
		Enabled:             pointy.Bool(cast.ToBool(awsKms["enabled"])),
		AccessKeyID:         awsKms["access_key_id"].(string),
		SecretAccessKey:     awsKms["secret_access_key"].(string),
		CustomerMasterKeyID: awsKms["customer_master_key_id"].(string),
		Region:              awsRegion,
	}
}

func resourceMongoDBAtlasEncryptionAtRestCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	encryptionAtRestReq := &matlas.EncryptionAtRest{
		GroupID: d.Get("project_id").(string),
		AwsKms:  expandAwsKms(d.Get("aws_kms").(map[string]interface{})),
		AzureKeyVault: matlas.AzureKeyVault{
			Enabled:           pointy.Bool(cast.ToBool(d.Get("azure_key_vault.enabled"))),
			ClientID:          d.Get("azure_key_vault.client_id").(string),
			AzureEnvironment:  d.Get("azure_key_vault.azure_environment").(string),
			SubscriptionID:    d.Get("azure_key_vault.subscription_id").(string),
			ResourceGroupName: d.Get("azure_key_vault.resource_group_name").(string),
			KeyVaultName:      d.Get("azure_key_vault.key_vault_name").(string),
			KeyIdentifier:     d.Get("azure_key_vault.key_identifier").(string),
			Secret:            d.Get("azure_key_vault.secret").(string),
			TenantID:          d.Get("azure_key_vault.tenant_id").(string),
		},
		GoogleCloudKms: matlas.GoogleCloudKms{
			Enabled:              pointy.Bool(cast.ToBool(d.Get("google_cloud_kms.enabled"))),
			ServiceAccountKey:    d.Get("google_cloud_kms.service_account_key").(string),
			KeyVersionResourceID: d.Get("google_cloud_kms.key_version_resource_id").(string),
		},
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

func resourceMongoDBAtlasEncryptionAtRestDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	_, err := conn.EncryptionsAtRest.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error deleting a encryptionAtRest (%s): %s", d.Id(), err)
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

	if d.HasChange("aws_kms") {
		encrypt.AwsKms = expandAwsKms(d.Get("aws_kms").(map[string]interface{}))
	}

	encrypt.GroupID = projectID

	_, _, err = conn.EncryptionsAtRest.Create(context.Background(), encrypt)

	if err != nil {
		return fmt.Errorf("error updating encryption at rest (%s): %s", projectID, err)
	}

	return resourceMongoDBAtlasEncryptionAtRestRead(d, meta)
}
