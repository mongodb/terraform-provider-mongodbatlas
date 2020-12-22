package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorCustomDNSConfigurationCreate  = "error creating custom dns configuration cluster aws information: %s"
	errorCustomDNSConfigurationRead    = "error getting custom dns configuration cluster aws information: %s"
	errorCustomDNSConfigurationUpdate  = "error updating custom dns configuration cluster aws information: %s"
	errorCustomDNSConfigurationDelete  = "error deleting custom dns configuration cluster aws (%s): %s"
	errorCustomDNSConfigurationSetting = "error setting `%s` for custom dns configuration cluster aws (%s): %s"
)

func resourceMongoDBAtlasCustomDNSConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasCustomDNSConfigurationCreate,
		Read:   resourceMongoDBAtlasCustomDNSConfigurationRead,
		Update: resourceMongoDBAtlasCustomDNSConfigurationUpdate,
		Delete: resourceMongoDBAtlasCustomDNSConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceMongoDBAtlasCustomDNSConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	orgID := d.Get("project_id").(string)

	// Creating(Updating) the Custom DNS Configuration for Atlas Clusters on AWS
	_, _, err := conn.CustomAWSDNS.Update(context.Background(), orgID,
		&matlas.AWSCustomDNSSetting{
			Enabled: d.Get("enabled").(bool),
		})
	if err != nil {
		return fmt.Errorf(errorCustomDNSConfigurationCreate, err)
	}

	d.SetId(orgID)

	return resourceMongoDBAtlasCustomDNSConfigurationRead(d, meta)
}

func resourceMongoDBAtlasCustomDNSConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	dnsResp, _, err := conn.CustomAWSDNS.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf(errorCustomDNSConfigurationRead, err)
	}

	if err = d.Set("enabled", dnsResp.Enabled); err != nil {
		return fmt.Errorf(errorCustomDNSConfigurationSetting, "name", d.Id(), err)
	}

	return nil
}

func resourceMongoDBAtlasCustomDNSConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	if d.HasChange("enabled") {
		_, _, err := conn.CustomAWSDNS.Update(context.Background(), d.Id(), &matlas.AWSCustomDNSSetting{
			Enabled: d.Get("enabled").(bool),
		})
		if err != nil {
			return fmt.Errorf(errorCustomDNSConfigurationUpdate, err)
		}
	}

	return resourceMongoDBAtlasCustomDNSConfigurationRead(d, meta)
}

func resourceMongoDBAtlasCustomDNSConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	_, _, err := conn.CustomAWSDNS.Update(context.Background(), d.Id(), &matlas.AWSCustomDNSSetting{
		Enabled: false,
	})
	if err != nil {
		return fmt.Errorf(errorCustomDNSConfigurationDelete, d.Id(), err)
	}

	return nil
}
