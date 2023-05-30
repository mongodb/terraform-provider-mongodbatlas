package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorClusterOutageSimulationCreate  = "error starting MongoDB Atlas Cluster Outage Simulation for Project (%s), Cluster (%s): %s"
	errorClusterOutageSimulationRead    = "error getting MongoDB Atlas Cluster Outage Simulation for Project (%s), Cluster (%s): %s"
	errorClusterOutageSimulationDelete  = "error ending MongoDB Atlas Cluster Outage Simulation for Project (%s), Cluster (%s): %s"
	errorClusterOutageSimulationSetting = "error setting `%s` for MongoDB Atlas Cluster Outage Simulation: %s"
	defaultOutageFilterType             = "REGION"
)

func resourceMongoDBAtlasClusterOutageSimulation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBClusterOutageSimulationCreate,
		ReadContext:   resourceMongoDBAClusterOutageSimulationRead,
		UpdateContext: resourceMongoDBClusterOutageSimulationUpdate,
		DeleteContext: resourceMongoDBAtlasClusterOutageSimulationDelete,
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(25 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"outage_filters": {
				Type:     schema.TypeList,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_provider": {
							Type:     schema.TypeString,
							Required: true,
						},
						"region_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"start_request_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"simulation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBClusterOutageSimulationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	requestBody := &matlas.ClusterOutageSimulationRequest{
		OutageFilters: newOutageFilters(d),
	}

	_, _, err := conn.ClusterOutageSimulation.StartOutageSimulation(ctx, projectID, clusterName, requestBody)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterOutageSimulationCreate, projectID, clusterName, err))
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"START_REQUESTED", "STARTING"},
		Target:     []string{"SIMULATING"},
		Refresh:    resourceClusterOutageSimulationRefreshFunc(ctx, clusterName, projectID, conn),
		Timeout:    timeout,
		MinTimeout: 1 * time.Minute,
		Delay:      3 * time.Minute,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterOutageSimulationCreate, projectID, clusterName, err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return resourceMongoDBAClusterOutageSimulationRead(ctx, d, meta)
}

func newOutageFilters(d *schema.ResourceData) []matlas.ClusterOutageSimulationOutageFilter {
	outageFilters := make([]matlas.ClusterOutageSimulationOutageFilter, len(d.Get("outage_filters").([]interface{})))

	for k, v := range d.Get("outage_filters").([]interface{}) {
		a := v.(map[string]interface{})
		outageFilters[k] = matlas.ClusterOutageSimulationOutageFilter{
			CloudProvider: pointy.String(a["cloud_provider"].(string)),
			RegionName:    pointy.String(a["region_name"].(string)),
			Type:          pointy.String(defaultOutageFilterType),
		}
	}

	return outageFilters
}

func resourceMongoDBAClusterOutageSimulationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	outageSimulation, resp, err := conn.ClusterOutageSimulation.GetOutageSimulation(ctx, projectID, clusterName)

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf(errorClusterOutageSimulationRead, projectID, clusterName, err))
	}

	if err = convertOutageSimulationToSchema(outageSimulation, d); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return nil
}

func resourceMongoDBAtlasClusterOutageSimulationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	_, _, err := conn.ClusterOutageSimulation.EndOutageSimulation(ctx, projectID, clusterName)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterOutageSimulationDelete, projectID, clusterName, err))
	}

	log.Println("[INFO] Waiting for MongoDB Cluster Outage Simulation to end")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"RECOVERY_REQUESTED", "RECOVERING", "COMPLETE"},
		Target:     []string{"DELETED"},
		Refresh:    resourceClusterOutageSimulationRefreshFunc(ctx, clusterName, projectID, conn),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterOutageSimulationDelete, projectID, clusterName, err))
	}

	return nil
}

func resourceMongoDBClusterOutageSimulationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.FromErr(fmt.Errorf("updating a Cluster Outage Simulation is not supported"))
}

func resourceClusterOutageSimulationRefreshFunc(ctx context.Context, clusterName, projectID string, client *matlas.Client) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		outageSimulation, resp, err := client.ClusterOutageSimulation.GetOutageSimulation(ctx, projectID, clusterName)

		if err != nil {
			if resp.StatusCode == 404 {
				return "", "DELETED", nil
			}
			return nil, "", err
		}

		if *outageSimulation.State != "" {
			log.Printf("[DEBUG] status for MongoDB cluster outage simulation: %s: %s", clusterName, *outageSimulation.State)
		}

		return outageSimulation, *outageSimulation.State, nil
	}
}

func convertOutageSimulationToSchema(outageSimulation *matlas.ClusterOutageSimulation, d *schema.ResourceData) error {
	if err := d.Set("state", outageSimulation.State); err != nil {
		return fmt.Errorf(errorClusterOutageSimulationSetting, "state", err)
	}
	if err := d.Set("start_request_date", outageSimulation.StartRequestDate); err != nil {
		return fmt.Errorf(errorClusterOutageSimulationSetting, "start_request_date", err)
	}
	if err := d.Set("simulation_id", outageSimulation.ID); err != nil {
		return fmt.Errorf(errorClusterOutageSimulationSetting, "simulation_id", err)
	}
	if err := d.Set("project_id", outageSimulation.GroupID); err != nil {
		return fmt.Errorf(errorClusterOutageSimulationSetting, "project_id", err)
	}
	if err := d.Set("cluster_name", outageSimulation.ClusterName); err != nil {
		return fmt.Errorf(errorClusterOutageSimulationSetting, "cluster_name", err)
	}

	if outageFilters := convertOutageFiltersToSchema(outageSimulation.OutageFilters, d); outageFilters != nil {
		if err := d.Set("outage_filters", outageFilters); err != nil {
			return fmt.Errorf(errorClusterOutageSimulationSetting, "outage_filters", err)
		}
	}
	return nil
}

func convertOutageFiltersToSchema(outageFilters []matlas.ClusterOutageSimulationOutageFilter, d *schema.ResourceData) []map[string]interface{} {
	outageFilterList := make([]map[string]interface{}, 0)
	for _, v := range outageFilters {
		outageFilterList = append(outageFilterList, map[string]interface{}{
			"cloud_provider": v.CloudProvider,
			"region_name":    v.RegionName,
			"type":           v.Type,
		})
	}
	return outageFilterList
}
