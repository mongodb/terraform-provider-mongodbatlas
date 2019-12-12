package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mwielbut/pointy"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	errorAuditingCreate = "error create MongoDB Auditing (%s): %s"
	errorAuditingUpdate = "error update MongoDB Auditing (%s): %s"
	errorAuditingRead   = "error reading MongoDB Auditing (%s): %s"
	// errorAuditingDelete = "error delete MongoDB Auditing (%s): %s"
)

func resourceMongoDBAtlasAuditing() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasAuditingCreate,
		Read:   resourceMongoDBAtlasAuditingRead,
		Update: resourceMongoDBAtlasAuditingUpdate,
		Delete: resourceMongoDBAtlasAuditingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"audit_authorization_success": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"audit_filter": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"configuration_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasAuditingCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)
	auditingReq := &matlas.Auditing{}

	if auditAuth, ok := d.GetOk("audit_authorization_success"); ok {
		auditingReq.AuditAuthorizationSuccess = pointy.Bool(auditAuth.(bool))
	}
	if auditFilter, ok := d.GetOk("audit_filter"); ok {
		auditingReq.AuditFilter = auditFilter.(string)
	}
	if enabled, ok := d.GetOk("enabled"); ok {
		auditingReq.Enabled = pointy.Bool(enabled.(bool))
	}

	_, _, err := conn.Auditing.Configure(context.Background(), projectID, auditingReq)
	if err != nil {
		return fmt.Errorf(errorAuditingCreate, projectID, err)
	}

	d.SetId(projectID)
	return resourceMongoDBAtlasAuditingRead(d, meta)
}

func resourceMongoDBAtlasAuditingRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	auditing, _, err := conn.Auditing.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf(errorAuditingRead, d.Id(), err)
	}

	if err := d.Set("audit_authorization_success", auditing.AuditAuthorizationSuccess); err != nil {
		return fmt.Errorf(errorAuditingRead, d.Id(), err)
	}
	if err := d.Set("audit_filter", auditing.AuditFilter); err != nil {
		return fmt.Errorf(errorAuditingRead, d.Id(), err)
	}
	if err := d.Set("enabled", auditing.Enabled); err != nil {
		return fmt.Errorf(errorAuditingRead, d.Id(), err)
	}
	if err := d.Set("configuration_type", auditing.ConfigurationType); err != nil {
		return fmt.Errorf(errorAuditingRead, d.Id(), err)
	}
	return nil
}

func resourceMongoDBAtlasAuditingUpdate(d *schema.ResourceData, meta interface{}) error {
	//Get the client connection.
	conn := meta.(*matlas.Client)

	auditingReq := &matlas.Auditing{}

	if d.HasChange("audit_authorization_success") {
		auditingReq.AuditAuthorizationSuccess = pointy.Bool((d.Get("audit_authorization_success").(bool)))
	}
	if d.HasChange("audit_filter") {
		auditingReq.AuditFilter = d.Get("audit_filter").(string)
	}
	if d.HasChange("enabled") {
		auditingReq.Enabled = pointy.Bool(d.Get("enabled").(bool))
	}

	_, _, err := conn.Auditing.Configure(context.Background(), d.Id(), auditingReq)
	if err != nil {
		return fmt.Errorf(errorAuditingUpdate, d.Id(), err)
	}

	return nil
}

func resourceMongoDBAtlasAuditingDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
