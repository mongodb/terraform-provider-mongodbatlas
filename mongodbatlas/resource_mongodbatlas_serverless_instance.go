package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func resourceMongoDBAtlasServerlessInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasServerlessInstanceCreate,
		ReadContext:   resourceMongoDBAtlasServerlessInstanceRead,
		UpdateContext: resourceMongoDBAtlasServerlessInstanceUpdate,
		DeleteContext: resourceMongoDBAtlasServerlessInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasServerlessInstanceImportState,
		},
		Schema: returnServerlessInstanceSchema(),
	}
}

func resourceMongoDBAtlasServerlessInstanceUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Serverless Instance Update not Implemented on MONGODB ATLAS API")
	return nil
}

func returnServerlessInstanceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"project_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"provider_settings_backing_provider_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"provider_settings_provider_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"provider_settings_region_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"connection_strings_standard_srv": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"create_date": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"mongo_db_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"links": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"href": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"rel": {
						Type:     schema.TypeString,
						Computed: true,
					},
				}},
		},
		"state_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}

func resourceMongoDBAtlasServerlessInstanceImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	projectID, name, err := splitServerlessInstanceImportID(d.Id())
	if err != nil {
		return nil, err
	}

	u, _, err := conn.ServerlessInstances.Get(ctx, *projectID, *name)
	if err != nil {
		return nil, fmt.Errorf("couldn't import cluster %s in project %s, error: %s", *name, *projectID, err)
	}

	if err := d.Set("project_id", u.GroupID); err != nil {
		log.Printf(errorClusterSetting, "project_id", u.ID, err)
	}

	if err := d.Set("name", u.Name); err != nil {
		log.Printf(errorClusterSetting, "name", u.ID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": *projectID,
		"name":       u.Name,
	}))

	return []*schema.ResourceData{d}, nil
}

func resourceMongoDBAtlasServerlessInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	serverlessName := ids["name"]

	_, err := conn.ServerlessInstances.Delete(ctx, projectID, serverlessName)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting MongoDB Serverless Instance (%s): %s", serverlessName, err))
	}

	log.Println("[INFO] Waiting for MongoDB Serverless Instance to be destroyed")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"IDLE", "CREATING", "UPDATING", "REPAIRING", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    resourceServerlessInstanceRefreshFunc(ctx, serverlessName, projectID, conn),
		Timeout:    3 * time.Hour,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute, // Wait 30 secs before starting
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting MongoDB Serverless Instance (%s): %s", serverlessName, err))
	}

	return nil
}

func resourceMongoDBAtlasServerlessInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	instanceName := ids["name"]

	serverlessInstance, _, err := conn.ServerlessInstances.Get(ctx, projectID, instanceName)
	if err != nil {
		// case 404
		// deleted in the backend case
		reset := strings.Contains(err.Error(), "404") && !d.IsNewResource()

		if reset {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting serverless instance information: %s", err)
	}

	if err := d.Set("id", serverlessInstance.ID); err != nil {
		return diag.Errorf("error setting `is` for serverless instance (%s): %s", d.Id(), err)
	}

	if err := d.Set("provider_settings_backing_provider_name", serverlessInstance.ProviderSettings.BackingProviderName); err != nil {
		return diag.Errorf("error setting `provider_settings_backing_provider_name` for serverless instance (%s): %s", d.Id(), err)
	}

	if err := d.Set("provider_settings_provider_name", serverlessInstance.ProviderSettings.ProviderName); err != nil {
		return diag.Errorf("error setting `provider_settings_provider_name` for serverless instance (%s): %s", d.Id(), err)
	}

	if err := d.Set("provider_settings_region_name", serverlessInstance.ProviderSettings.RegionName); err != nil {
		return diag.Errorf("error setting `provider_settings_region_name` for serverless instance (%s): %s", d.Id(), err)
	}

	if err := d.Set("connection_strings_standard_srv", serverlessInstance.ConnectionStrings.StandardSrv); err != nil {
		return diag.Errorf("error setting `connection_strings_standard_srv` for serverless instance (%s): %s", d.Id(), err)
	}

	if err := d.Set("create_date", serverlessInstance.CreateDate); err != nil {
		return diag.Errorf("error setting `create_date` for serverless instance (%s): %s", d.Id(), err)
	}

	if err := d.Set("mongo_db_version", serverlessInstance.MongoDBVersion); err != nil {
		return diag.Errorf("error setting `mongo_db_version` for serverless instance (%s): %s", d.Id(), err)
	}

	if err := d.Set("links", flattenServerlessInstanceLinks(serverlessInstance.Links)); err != nil {
		return diag.Errorf("error setting `links` for serverless instance (%s): %s", d.Id(), err)
	}

	if err := d.Set("state_name", serverlessInstance.StateName); err != nil {
		return diag.Errorf("error setting `state_name` for serverless instance (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceMongoDBAtlasServerlessInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	name := d.Get("name").(string)

	serverlessProviderSettings := &matlas.ServerlessProviderSettings{
		BackingProviderName: d.Get("provider_settings_backing_provider_name").(string),
		ProviderName:        d.Get("provider_settings_provider_name").(string),
		RegionName:          d.Get("provider_settings_region_name").(string),
	}

	serverlessInstanceRequest := &matlas.ServerlessCreateRequestParams{
		Name:             name,
		ProviderSettings: serverlessProviderSettings,
	}

	_, _, err := conn.ServerlessInstances.Create(ctx, projectID, serverlessInstanceRequest)
	if err != nil {
		return diag.Errorf("error creating serverless instance: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING", "REPEATING", "PENDING"},
		Target:     []string{"IDLE"},
		Refresh:    resourceServerlessInstanceRefreshFunc(ctx, d.Get("name").(string), projectID, conn),
		Timeout:    3 * time.Hour,
		MinTimeout: 1 * time.Minute,
		Delay:      3 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error creating MongoDB Serverless Instance: %s", err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"name":       name,
	}))

	return resourceMongoDBAtlasServerlessInstanceRead(ctx, d, meta)
}

func resourceServerlessInstanceRefreshFunc(ctx context.Context, name, projectID string, client *matlas.Client) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, resp, err := client.ServerlessInstances.Get(ctx, projectID, name)

		if err != nil && strings.Contains(err.Error(), "reset by peer") {
			return nil, "REPEATING", nil
		}

		if err != nil && c == nil && resp == nil {
			return nil, "", err
		} else if err != nil {
			if resp.StatusCode == 404 {
				return "", "DELETED", nil
			}
			if resp.StatusCode == 503 {
				return "", "PENDING", nil
			}
			return nil, "", err
		}

		if c.StateName != "" {
			log.Printf("[DEBUG] status for MongoDB Serverless Instance: %s: %s", name, c.StateName)
		}

		return c, c.StateName, nil
	}
}

func flattenServerlessInstanceLinks(links []*matlas.Link) []map[string]interface{} {
	linksList := make([]map[string]interface{}, 0)

	for _, link := range links {
		mLink := map[string]interface{}{
			"href": link.Href,
			"rel":  link.Rel,
		}
		linksList = append(linksList, mLink)
	}

	return linksList
}

func splitServerlessInstanceImportID(id string) (projectID, instanceName *string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("import format error: to import a serverless instance, use the format {project_id}-{name}")
		return
	}

	projectID = &parts[1]
	instanceName = &parts[2]

	return
}
