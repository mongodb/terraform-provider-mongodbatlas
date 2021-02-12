package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasCustomDNSConfigurationAWS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasCustomDNSConfigurationAWSRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasCustomDNSConfigurationAWSRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)

	customDNSSetting, _, err := conn.CustomAWSDNS.Get(context.Background(), projectID)
	if err != nil {
		return fmt.Errorf(errorCustomDNSConfigurationRead, err)
	}

	if err := d.Set("enabled", customDNSSetting.Enabled); err != nil {
		return fmt.Errorf(errorCustomDNSConfigurationSetting, "enabled", projectID, err)
	}

	d.SetId(projectID)

	return nil
}
