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
		CreateContext: resourceMongoDBAtlasEncryptionAtRestCreate,
		ReadContext:   resourceMongoDBAtlasEncryptionAtRestRead,
		DeleteContext: resourceMongoDBAtlasEncryptionAtRestDelete,
		UpdateContext: resourceMongoDBAtlasEncryptionAtRestUpdate,
		Importer:      &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"aws_kms": {
				Type:          schema.TypeMap,
				Optional:      true,
				Sensitive:     true,
				Deprecated:    "use aws_kms_config instead",
				ConflictsWith: []string{"aws_kms_config"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(map[string]interface{})

					_, akOk := v["access_key_id"]
					_, saOk := v["secret_access_key"]
					_, rOk := v["role_id"]

					if (akOk && saOk && rOk) || (akOk && rOk) || (saOk && rOk) {
						errs = append(errs, fmt.Errorf("%q For credentials: `access_key_id` and `secret_access_key` are allowed but not `role_id`."+
							" For roles: `access_key_id` and `secret_access_key` are not allowed but `role_id` is allowed", key))
					}

					return
				},
			},
			"azure_key_vault": {
				Type:          schema.TypeMap,
				Optional:      true,
				Sensitive:     true,
				Deprecated:    "use azure_key_vault_config instead",
				ConflictsWith: []string{"azure_key_vault_config"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"google_cloud_kms": {
				Type:          schema.TypeMap,
				Optional:      true,
				Sensitive:     true,
				Deprecated:    "use google_cloud_kms_config instead",
				ConflictsWith: []string{"google_cloud_kms_config"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"aws_kms_config": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"aws_kms"},
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
			"azure_key_vault_config": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"azure_key_vault"},
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
			"google_cloud_kms_config": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"google_cloud_kms"},
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

	// Deprecated workflows
	aws, awsOk := d.GetOk("aws_kms")
	if awsOk {
		encryptionAtRestReq.AwsKms = expandAwsKms(aws.(map[string]interface{}))
	}
	azure, azureOk := d.GetOk("azure_key_vault")
	if azureOk {
		encryptionAtRestReq.AzureKeyVault = expandAzureKeyVault(azure.(map[string]interface{}))
	}
	gcp, gcpOk := d.GetOk("google_cloud_kms")
	if gcpOk {
		encryptionAtRestReq.GoogleCloudKms = expandGCPKms(gcp.(map[string]interface{}))
	}

	// End Depprecated workflows

	awsC, awsCOk := d.GetOk("aws_kms_config")
	if awsCOk {
		err := validateAwsKms(awsC.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		encryptionAtRestReq.AwsKms = expandAwsKmsConfig(awsC.([]interface{}))
	}
	azureC, azureCOk := d.GetOk("azure_key_vault_config")
	if azureCOk {
		encryptionAtRestReq.AzureKeyVault = expandAzureKeyVaultConfig(azureC.([]interface{}))
	}
	gcpC, gcpCOk := d.GetOk("google_cloud_kms_config")
	if gcpCOk {
		encryptionAtRestReq.GoogleCloudKms = expandGCPKmsConfig(gcpC.([]interface{}))
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

	// Deprecated workflows
	values := flattenAWSKMS(&resp.AwsKms)
	if !counterEmptyValues(values) {
		aws, awsOk := d.GetOk("aws_kms")
		if awsOk {
			aws2 := aws.(map[string]interface{})
			values["secret_access_key"] = cast.ToString(aws2["secret_access_key"])
			if v, sa := values["role_id"]; sa {
				if v.(string) == "" {
					delete(values, "role_id")
				}
			}
			if v, sa := values["access_key_id"]; sa {
				if v.(string) == "" {
					delete(values, "access_key_id")
					delete(values, "secret_access_key")
				}
			}
			if err = d.Set("aws_kms", values); err != nil {
				return diag.FromErr(fmt.Errorf(errorAlertEncryptionAtRestSetting, "aws_kms", d.Id(), err))
			}
		}
	}

	values = flattenAzureVault(&resp.AzureKeyVault)
	if !counterEmptyValues(values) {
		azure, azureOk := d.GetOk("azure_key_vault")
		if azureOk {
			azure2 := azure.(map[string]interface{})
			values["secret"] = cast.ToString(azure2["secret"])
			if v, sa := values["secret"]; sa {
				if v.(string) == "" {
					delete(values, "secret")
				}
			}
			if err = d.Set("azure_key_vault", values); err != nil {
				return diag.FromErr(fmt.Errorf(errorAlertEncryptionAtRestSetting, "azure_key_vault", d.Id(), err))
			}
		}
	}

	values = flattenGCPKms(&resp.GoogleCloudKms)
	if !counterEmptyValues(values) {
		gcp, gcpOk := d.GetOk("google_cloud_kms")
		if gcpOk {
			gcp2 := gcp.(map[string]interface{})
			values["service_account_key"] = cast.ToString(gcp2["service_account_key"])
			if v, sa := values["service_account_key"]; sa {
				if v.(string) == "" {
					delete(values, "service_account_key")
				}
			}
			if err = d.Set("google_cloud_kms", values); err != nil {
				return diag.FromErr(fmt.Errorf(errorAlertEncryptionAtRestSetting, "google_cloud_kms", d.Id(), err))
			}
		}
	}
	// End Deprecated workflows

	values2 := flattenAWSKMSConfig(&resp.AwsKms)
	if len(values2) != 0 {
		value := values2[0].(map[string]interface{})
		if !counterEmptyValues(value) {
			aws, awsOk := d.GetOk("aws_kms_config")
			if awsOk {
				awsObj := aws.([]interface{})
				if len(awsObj) > 0 {
					aws2 := awsObj[0].(map[string]interface{})
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
			}
			values2[0] = value

			if err = d.Set("aws_kms_config", values2); err != nil {
				return diag.FromErr(fmt.Errorf(errorAlertEncryptionAtRestSetting, "aws_kms_config", d.Id(), err))
			}
		}
	}

	values2 = flattenAzureVaultConfig(&resp.AzureKeyVault)
	if len(values2) != 0 {
		value := values2[0].(map[string]interface{})

		if !counterEmptyValues(value) {
			azure, azureOk := d.GetOk("azure_key_vault_config")
			if azureOk {
				azureObj := azure.([]interface{})
				if len(azureObj) > 0 {
					azure2 := azureObj[0].(map[string]interface{})
					value["secret"] = cast.ToString(azure2["secret"])
					if v, sa := value["secret"]; sa {
						if v.(string) == "" {
							delete(value, "secret")
						}
					}
				}
			}
			values2[0] = value

			if err = d.Set("azure_key_vault_config", values2); err != nil {
				return diag.FromErr(fmt.Errorf(errorAlertEncryptionAtRestSetting, "azure_key_vault_config", d.Id(), err))
			}
		}
	}

	values2 = flattenGCPKmsConfig(&resp.GoogleCloudKms)
	if len(values2) != 0 {
		value := values2[0].(map[string]interface{})
		if !counterEmptyValues(value) {
			gcp, gcpOk := d.GetOk("google_cloud_kms_config")
			if gcpOk {
				gcpObj := gcp.([]interface{})
				if len(gcpObj) > 0 {
					gcp2 := gcpObj[0].(map[string]interface{})
					value["service_account_key"] = cast.ToString(gcp2["service_account_key"])
					if v, sa := value["service_account_key"]; sa {
						if v.(string) == "" {
							delete(value, "service_account_key")
						}
					}
				}
			}
			values2[0] = value

			if err = d.Set("google_cloud_kms_config", values2); err != nil {
				return diag.FromErr(fmt.Errorf(errorAlertEncryptionAtRestSetting, "google_cloud_kms_config", d.Id(), err))
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
		encrypt.AwsKms = expandAwsKms(d.Get("aws_kms").(map[string]interface{}))
	}

	if d.HasChange("azure_key_vault") {
		encrypt.AzureKeyVault = expandAzureKeyVault(d.Get("azure_key_vault").(map[string]interface{}))
	}

	if d.HasChange("google_cloud_kms") {
		encrypt.GoogleCloudKms = expandGCPKms(d.Get("google_cloud_kms").(map[string]interface{}))
	}

	if d.HasChange("aws_kms_config") {
		encrypt.AwsKms = expandAwsKmsConfig(d.Get("aws_kms_config").([]interface{}))
	}

	if d.HasChange("azure_key_vault_config") {
		encrypt.AzureKeyVault = expandAzureKeyVaultConfig(d.Get("azure_key_vault_config").([]interface{}))
	}

	if d.HasChange("google_cloud_kms_config") {
		encrypt.GoogleCloudKms = expandGCPKmsConfig(d.Get("google_cloud_kms_config").([]interface{}))
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

func flattenAWSKMS(m *matlas.AwsKms) map[string]interface{} {
	return map[string]interface{}{
		"enabled":                cast.ToString(m.Enabled),
		"access_key_id":          m.AccessKeyID,
		"customer_master_key_id": m.CustomerMasterKeyID,
		"region":                 m.Region,
		"role_id":                m.RoleID,
	}
}

func flattenAzureVault(m *matlas.AzureKeyVault) map[string]interface{} {
	return map[string]interface{}{
		"enabled":             cast.ToString(m.Enabled),
		"client_id":           m.ClientID,
		"azure_environment":   m.AzureEnvironment,
		"subscription_id":     m.SubscriptionID,
		"resource_group_name": m.ResourceGroupName,
		"key_vault_name":      m.KeyVaultName,
		"key_identifier":      m.KeyIdentifier,
		"secret":              m.Secret,
		"tenant_id":           m.TenantID,
	}
}

func flattenGCPKms(m *matlas.GoogleCloudKms) map[string]interface{} {
	return map[string]interface{}{
		"enabled":                 cast.ToString(m.Enabled),
		"service_account_key":     m.ServiceAccountKey,
		"key_version_resource_id": m.KeyVersionResourceID,
	}
}

func expandAwsKmsConfig(awsKms []interface{}) matlas.AwsKms {
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
		return fmt.Errorf("empty aws_kms_config")
	}

	v := awsKms[0].(map[string]interface{})

	ak, akOk := v["access_key_id"]
	sa, saOk := v["secret_access_key"]
	r, rOk := v["role_id"]

	if ((akOk && ak != "") && (saOk && sa != "") && (rOk && r != "")) || ((akOk && ak != "") && (rOk && r != "")) || ((saOk && sa != "") && (rOk && r != "")) {
		return fmt.Errorf("for credentials: `access_key_id` and `secret_access_key` are allowed but not `role_id`." +
			" For roles: `access_key_id` and `secret_access_key` are not allowed but `role_id` is allowed")
	}

	return nil
}

func expandAzureKeyVaultConfig(azure []interface{}) matlas.AzureKeyVault {
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

func expandGCPKmsConfig(gcpKms []interface{}) matlas.GoogleCloudKms {
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

func flattenAWSKMSConfig(m *matlas.AwsKms) []interface{} {
	objMap := make(map[string]interface{}, 1)

	if cast.ToBool(m.Enabled) {
		objMap["enabled"] = m.Enabled
	}

	objMap["access_key_id"] = m.AccessKeyID
	objMap["customer_master_key_id"] = m.CustomerMasterKeyID
	objMap["region"] = m.Region
	objMap["role_id"] = m.RoleID

	return []interface{}{objMap}
}

func flattenAzureVaultConfig(m *matlas.AzureKeyVault) []interface{} {
	objMap := make(map[string]interface{}, 1)

	if cast.ToBool(m.Enabled) {
		objMap["enabled"] = m.Enabled
	}

	objMap["client_id"] = m.ClientID
	objMap["azure_environment"] = m.AzureEnvironment
	objMap["subscription_id"] = m.SubscriptionID
	objMap["resource_group_name"] = m.ResourceGroupName
	objMap["key_vault_name"] = m.KeyVaultName
	objMap["key_identifier"] = m.KeyIdentifier
	objMap["secret"] = m.Secret
	objMap["tenant_id"] = m.TenantID

	return []interface{}{objMap}
}

func flattenGCPKmsConfig(m *matlas.GoogleCloudKms) []interface{} {
	objMap := make(map[string]interface{}, 1)

	if cast.ToBool(m.Enabled) {
		objMap["enabled"] = m.Enabled
	}
	objMap["service_account_key"] = m.ServiceAccountKey
	objMap["key_version_resource_id"] = m.KeyVersionResourceID

	return []interface{}{objMap}
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
