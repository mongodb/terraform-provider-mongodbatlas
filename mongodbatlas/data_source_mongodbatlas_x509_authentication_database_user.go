package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasX509AuthDBUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasX509AuthDBUserRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"customer_x509_cas": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"not_after": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subject": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasX509AuthDBUserRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	username := d.Get("username").(string)

	if username != "" {
		certificates, _, err := conn.X509AuthDBUsers.GetUserCertificates(context.Background(), projectID, username)
		if err != nil {
			return fmt.Errorf(errorX509AuthDBUsersRead, username, projectID, err)
		}

		if err := d.Set("certificates", flattenCertificates(certificates)); err != nil {
			return fmt.Errorf(errorX509AuthDBUsersSetting, "certificates", username, err)
		}
	}

	customerX509, _, err := conn.X509AuthDBUsers.GetCurrentX509Conf(context.Background(), projectID)
	if err != nil {
		return fmt.Errorf(errorCustomerX509AuthDBUsersRead, projectID, err)
	}

	if err := d.Set("customer_x509_cas", customerX509.Cas); err != nil {
		return fmt.Errorf(errorX509AuthDBUsersSetting, "certificates", username, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":          projectID,
		"username":            username,
		"current_certificate": "",
	}))

	return nil
}
