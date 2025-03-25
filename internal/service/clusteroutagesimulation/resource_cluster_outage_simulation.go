package clusteroutagesimulation

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312001/admin"
)

const (
	errorClusterOutageSimulationCreate  = "error starting MongoDB Atlas Cluster Outage Simulation for Project (%s), Cluster (%s): %s"
	errorClusterOutageSimulationRead    = "error getting MongoDB Atlas Cluster Outage Simulation for Project (%s), Cluster (%s): %s"
	errorClusterOutageSimulationDelete  = "error ending MongoDB Atlas Cluster Outage Simulation for Project (%s), Cluster (%s): %s"
	errorClusterOutageSimulationSetting = "error setting `%s` for MongoDB Atlas Cluster Outage Simulation: %s"
	defaultOutageFilterType             = "REGION"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
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

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	requestBody := admin.ClusterOutageSimulation{
		OutageFilters: newOutageFilters(d),
	}

	_, _, err := connV2.ClusterOutageSimulationApi.StartOutageSimulation(ctx, projectID, clusterName, &requestBody).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterOutageSimulationCreate, projectID, clusterName, err))
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"START_REQUESTED", "STARTING"},
		Target:     []string{"SIMULATING"},
		Refresh:    resourceRefreshFunc(ctx, clusterName, projectID, connV2),
		Timeout:    timeout,
		MinTimeout: 1 * time.Minute,
		Delay:      3 * time.Minute,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterOutageSimulationCreate, projectID, clusterName, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return resourceRead(ctx, d, meta)
}

func newOutageFilters(d *schema.ResourceData) *[]admin.AtlasClusterOutageSimulationOutageFilter {
	outageFilters := make([]admin.AtlasClusterOutageSimulationOutageFilter, len(d.Get("outage_filters").([]any)))

	for k, v := range d.Get("outage_filters").([]any) {
		a := v.(map[string]any)
		outageFilters[k] = admin.AtlasClusterOutageSimulationOutageFilter{
			CloudProvider: conversion.StringPtr(a["cloud_provider"].(string)),
			RegionName:    conversion.StringPtr(a["region_name"].(string)),
			Type:          conversion.StringPtr(defaultOutageFilterType),
		}
	}

	return &outageFilters
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	outageSimulation, resp, err := connV2.ClusterOutageSimulationApi.GetOutageSimulation(ctx, projectID, clusterName).Execute()

	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf(errorClusterOutageSimulationRead, projectID, clusterName, err))
	}

	if err = convertOutageSimulationToSchema(outageSimulation, d); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return nil
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	_, _, err := connV2.ClusterOutageSimulationApi.EndOutageSimulation(ctx, projectID, clusterName).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterOutageSimulationDelete, projectID, clusterName, err))
	}

	log.Println("[INFO] Waiting for MongoDB Cluster Outage Simulation to end")

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"RECOVERY_REQUESTED", "RECOVERING", "COMPLETE"},
		Target:     []string{"DELETED"},
		Refresh:    resourceRefreshFunc(ctx, clusterName, projectID, connV2),
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

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return diag.FromErr(fmt.Errorf("updating a Cluster Outage Simulation is not supported"))
}

func resourceRefreshFunc(ctx context.Context, clusterName, projectID string, client *admin.APIClient) retry.StateRefreshFunc {
	return func() (any, string, error) {
		outageSimulation, resp, err := client.ClusterOutageSimulationApi.GetOutageSimulation(ctx, projectID, clusterName).Execute()

		if err != nil {
			if validate.StatusNotFound(resp) {
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

func convertOutageSimulationToSchema(outageSimulation *admin.ClusterOutageSimulation, d *schema.ResourceData) error {
	if err := d.Set("state", outageSimulation.GetState()); err != nil {
		return fmt.Errorf(errorClusterOutageSimulationSetting, "state", err)
	}
	if err := d.Set("start_request_date", conversion.TimePtrToStringPtr(outageSimulation.StartRequestDate)); err != nil {
		return fmt.Errorf(errorClusterOutageSimulationSetting, "start_request_date", err)
	}
	if err := d.Set("simulation_id", outageSimulation.GetId()); err != nil {
		return fmt.Errorf(errorClusterOutageSimulationSetting, "simulation_id", err)
	}
	if err := d.Set("project_id", outageSimulation.GetGroupId()); err != nil {
		return fmt.Errorf(errorClusterOutageSimulationSetting, "project_id", err)
	}
	if err := d.Set("cluster_name", outageSimulation.GetClusterName()); err != nil {
		return fmt.Errorf(errorClusterOutageSimulationSetting, "cluster_name", err)
	}

	if outageFilters := convertOutageFiltersToSchema(outageSimulation.GetOutageFilters(), d); outageFilters != nil {
		if err := d.Set("outage_filters", outageFilters); err != nil {
			return fmt.Errorf(errorClusterOutageSimulationSetting, "outage_filters", err)
		}
	}
	return nil
}

func convertOutageFiltersToSchema(outageFilters []admin.AtlasClusterOutageSimulationOutageFilter, d *schema.ResourceData) []map[string]any {
	outageFilterList := make([]map[string]any, 0)
	for _, v := range outageFilters {
		outageFilterList = append(outageFilterList, map[string]any{
			"cloud_provider": v.GetCloudProvider(),
			"region_name":    v.GetRegionName(),
			"type":           v.GetType(),
		})
	}
	return outageFilterList
}
