package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorCreateEncryptionAtRest       = "error creating Encryption At Rest: %s"
	errorReadEncryptionAtRest         = "error getting Encryption At Rest: %s"
	errorDeleteEncryptionAtRest       = "error deleting Encryption At Rest: (%s): %s"
	errorUpdateEncryptionAtRest       = "error updating Encryption At Rest: %s"
	errorAlertEncryptionAtRestSetting = "error setting `%s` for Encryption At Rest (%s): %s"
)

func resourceMongoDBAtlasEncryptionAtRest() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceMongoDBAtlasEncryptionAtRestCreate,
		ReadWithoutTimeout:   resourceMongoDBAtlasEncryptionAtRestRead,
		DeleteWithoutTimeout: resourceMongoDBAtlasEncryptionAtRestDelete,
		UpdateWithoutTimeout: resourceMongoDBAtlasEncryptionAtRestUpdate,
		Importer:             &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"aws_kms": {
				Type:      schema.TypeList,
				MaxItems:  1,
				Optional:  true,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"access_key_id": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"secret_access_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"customer_master_key_id": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"role_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"azure_key_vault": {
				Type:      schema.TypeList,
				MaxItems:  1,
				Optional:  true,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"client_id": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"azure_environment": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subscription_id": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"resource_group_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"key_vault_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"key_identifier": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"secret": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"tenant_id": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
			"google_cloud_kms": {
				Type:      schema.TypeList,
				MaxItems:  1,
				Optional:  true,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"service_account_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"key_version_resource_id": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

func resourceMongoDBAtlasEncryptionAtRestCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	encryptionAtRestReq := &matlas.EncryptionAtRest{
		GroupID: d.Get("project_id").(string),
	}

	aws, awsOk := d.GetOk("aws_kms")
	if awsOk {
		err := validateAwsKms(aws.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		encryptionAtRestReq.AwsKms = expandAwsKms(aws.([]interface{}))
	}
	azure, azureOk := d.GetOk("azure_key_vault")
	if azureOk {
		encryptionAtRestReq.AzureKeyVault = expandAzureKeyVault(azure.([]interface{}))
	}
	gcp, gcpOk := d.GetOk("google_cloud_kms")
	if gcpOk {
		encryptionAtRestReq.GoogleCloudKms = expandGCPKms(gcp.([]interface{}))
	}

	_, _, err := conn.EncryptionsAtRest.Create(ctx, encryptionAtRestReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCreateEncryptionAtRest, err))
	}

	d.SetId(d.Get("project_id").(string))

	return resourceMongoDBAtlasEncryptionAtRestRead(ctx, d, meta)
}

func resourceMongoDBAtlasEncryptionAtRestRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	resp, response, err := conn.EncryptionsAtRest.Get(context.Background(), d.Id())
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorReadEncryptionAtRest, err))
	}

	values := flattenAWSKMS(&resp.AwsKms)
	if len(values) != 0 {
		value := values[0].(map[string]interface{})
		if !counterEmptyValues(value) {
			aws, awsOk := d.GetOk("aws_kms")
			if awsOk {
				aws2 := aws.(map[string]interface{})
				value["secret_access_key"] = cast.ToString(aws2["secret_access_key"])
				if v, sa := value["role_id"]; sa {
					if v.(string) == "" {
						delete(value, "role_id")
					}
				}
				if v, sa := value["access_key_id"]; sa {
					if v.(string) == "" {
						delete(value, "access_key_id")
						delete(value, "secret_access_key")
					}
				}
			}
			values[0] = value

			if err = d.Set("aws_kms", values); err != nil {
				return diag.FromErr(fmt.Errorf(errorAlertEncryptionAtRestSetting, "aws_kms", d.Id(), err))
			}
		}
	}

	values = flattenAzureVault(&resp.AzureKeyVault)
	if len(values) != 0 {
		value := values[0].(map[string]interface{})

		if !counterEmptyValues(value) {
			azure, azureOk := d.GetOk("azure_key_vault")
			if azureOk {
				azure2 := azure.(map[string]interface{})
				value["secret"] = cast.ToString(azure2["secret"])
				if v, sa := value["secret"]; sa {
					if v.(string) == "" {
						delete(value, "secret")
					}
				}
			}
			values[0] = value

			if err = d.Set("azure_key_vault", values); err != nil {
				return diag.FromErr(fmt.Errorf(errorAlertEncryptionAtRestSetting, "azure_key_vault", d.Id(), err))
			}
		}
	}

	values = flattenGCPKms(&resp.GoogleCloudKms)
	if len(values) != 0 {
		value := values[0].(map[string]interface{})
		if !counterEmptyValues(value) {
			gcp, gcpOk := d.GetOk("google_cloud_kms")
			if gcpOk {
				gcp2 := gcp.(map[string]interface{})
				value["service_account_key"] = cast.ToString(gcp2["service_account_key"])
				if v, sa := value["service_account_key"]; sa {
					if v.(string) == "" {
						delete(value, "service_account_key")
					}
				}
			}
			values[0] = value

			if err = d.Set("google_cloud_kms", values); err != nil {
				return diag.FromErr(fmt.Errorf(errorAlertEncryptionAtRestSetting, "google_cloud_kms", d.Id(), err))
			}
		}
	}

	if err = d.Set("project_id", d.Id()); err != nil {
		return diag.FromErr(fmt.Errorf(errorAlertEncryptionAtRestSetting, "project_id", d.Id(), err))
	}

	return nil
}

func resourceMongoDBAtlasEncryptionAtRestUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Id()

	encrypt, _, err := conn.EncryptionsAtRest.Get(ctx, projectID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorUpdateEncryptionAtRest, err))
	}

	encrypt.GroupID = projectID

	if d.HasChange("aws_kms") {
		encrypt.AwsKms = expandAwsKms(d.Get("aws_kms").([]interface{}))
	}

	if d.HasChange("azure_key_vault") {
		encrypt.AzureKeyVault = expandAzureKeyVault(d.Get("azure_key_vault").([]interface{}))
	}

	if d.HasChange("google_cloud_kms") {
		encrypt.GoogleCloudKms = expandGCPKms(d.Get("google_cloud_kms").([]interface{}))
	}

	_, _, err = conn.EncryptionsAtRest.Create(ctx, encrypt)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating encryption at rest (%s): %s", projectID, err))
	}

	return resourceMongoDBAtlasEncryptionAtRestRead(ctx, d, meta)
}

func resourceMongoDBAtlasEncryptionAtRestDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	_, err := conn.EncryptionsAtRest.Delete(ctx, d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDeleteEncryptionAtRest, d.Id(), err))
	}

	return nil
}

func expandAwsKms(awsKms []interface{}) matlas.AwsKms {
	if len(awsKms) == 0 {
		return matlas.AwsKms{}
	}

	awsKmsObj := awsKms[0].(map[string]interface{})

	awsRegion, _ := valRegion(awsKmsObj["region"])

	return matlas.AwsKms{
		Enabled:             pointy.Bool(cast.ToBool(awsKmsObj["enabled"])),
		AccessKeyID:         cast.ToString(awsKmsObj["access_key_id"]),
		SecretAccessKey:     cast.ToString(awsKmsObj["secret_access_key"]),
		CustomerMasterKeyID: cast.ToString(awsKmsObj["customer_master_key_id"]),
		Region:              awsRegion,
		RoleID:              cast.ToString(awsKmsObj["role_id"]),
	}
}

func validateAwsKms(awsKms []interface{}) error {
	if len(awsKms) == 0 {
		return fmt.Errorf("empty aws_kms")
	}

	v := awsKms[0].(map[string]interface{})

	_, akOk := v["access_key_id"]
	_, saOk := v["secret_access_key"]
	_, rOk := v["role_id"]

	if (akOk && saOk && rOk) || (akOk && rOk) || (saOk && rOk) {
		return fmt.Errorf("for credentials: `access_key_id` and `secret_access_key` are allowed but not `role_id`." +
			" For roles: `access_key_id` and `secret_access_key` are not allowed but `role_id` is allowed")
	}

	return nil
}

func expandAzureKeyVault(azure []interface{}) matlas.AzureKeyVault {
	if len(azure) == 0 {
		return matlas.AzureKeyVault{}
	}

	azureObj := azure[0].(map[string]interface{})

	return matlas.AzureKeyVault{
		Enabled:           pointy.Bool(cast.ToBool(azureObj["enabled"])),
		ClientID:          cast.ToString(azureObj["client_id"]),
		AzureEnvironment:  cast.ToString(azureObj["azure_environment"]),
		SubscriptionID:    cast.ToString(azureObj["subscription_id"]),
		ResourceGroupName: cast.ToString(azureObj["resource_group_name"]),
		KeyVaultName:      cast.ToString(azureObj["key_vault_name"]),
		KeyIdentifier:     cast.ToString(azureObj["key_identifier"]),
		Secret:            cast.ToString(azureObj["secret"]),
		TenantID:          cast.ToString(azureObj["tenant_id"]),
	}
}

func expandGCPKms(gcpKms []interface{}) matlas.GoogleCloudKms {
	if len(gcpKms) == 0 {
		return matlas.GoogleCloudKms{}
	}

	gcpKmsObj := gcpKms[0].(map[string]interface{})

	return matlas.GoogleCloudKms{
		Enabled:              pointy.Bool(cast.ToBool(gcpKmsObj["enabled"])),
		ServiceAccountKey:    cast.ToString(gcpKmsObj["service_account_key"]),
		KeyVersionResourceID: cast.ToString(gcpKmsObj["key_version_resource_id"]),
	}
}

func flattenAWSKMS(m *matlas.AwsKms) []interface{} {
	return []interface{}{map[string]interface{}{
		"enabled":                cast.ToString(m.Enabled),
		"access_key_id":          m.AccessKeyID,
		"customer_master_key_id": m.CustomerMasterKeyID,
		"region":                 m.Region,
		"role_id":                m.RoleID,
	}}
}

func flattenAzureVault(m *matlas.AzureKeyVault) []interface{} {
	return []interface{}{map[string]interface{}{
		"enabled":             cast.ToString(m.Enabled),
		"client_id":           m.ClientID,
		"azure_environment":   m.AzureEnvironment,
		"subscription_id":     m.SubscriptionID,
		"resource_group_name": m.ResourceGroupName,
		"key_vault_name":      m.KeyVaultName,
		"key_identifier":      m.KeyIdentifier,
		"secret":              m.Secret,
		"tenant_id":           m.TenantID,
	}}
}

func flattenGCPKms(m *matlas.GoogleCloudKms) []interface{} {
	return []interface{}{map[string]interface{}{
		"enabled":                 cast.ToString(m.Enabled),
		"service_account_key":     m.ServiceAccountKey,
		"key_version_resource_id": m.KeyVersionResourceID,
	}}
}

func counterEmptyValues(values map[string]interface{}) bool {
	count := 0
	for i := range values {
		if val, ok := values[i]; ok {
			strval, okT := val.(string)
			if okT && strval == "" || strval == "false" {
				count++
			}
		}
	}

	return len(values) == count
}
