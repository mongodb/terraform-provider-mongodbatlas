package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorDataLakePipelineCreate      = "error creating MongoDB Atlas DataLake Pipeline: %s"
	errorDataLakePipelineRead        = "error reading MongoDB Atlas DataLake Pipeline (%s): %s"
	errorDataLakePipelineImport      = "error importing MongoDB Atlas DataLake Pipeline (%s): %s"
	errorDataLakePipelineImportField = "error importing field (%s) for MongoDB Atlas DataLake Pipeline (%s): %s"
	errorDataLakePipelineDelete      = "error deleting MongoDB Atlas DataLake Pipeline (%s): %s"
	errorDataLakePipelineUpdate      = "error updating MongoDB Atlas DataLake Pipeline: %s"
	errorDataLakePipelineSetting     = "error setting `%s` for MongoDB Atlas DataLake Pipeline (%s): %s"
)

func resourceMongoDBAtlasDataLakePipeline() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasDataLakePipelineCreate,
		ReadContext:   resourceMongoDBAtlasDataLakePipelineRead,
		UpdateContext: resourceMongoDBAtlasDataLakePipelineUpdate,
		DeleteContext: resourceMongoDBAtlasDataLakePipelineDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasDataLakePipelineImportState,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sink": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"provider": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"partition_fields": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"order": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"source": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"collection_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"database_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"policy_item_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"transformations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"snapshots":           schemaDataLakePipelineSnapshots(),
			"ingestion_schedules": schemaDataLakePipelineIngestionSchedules(),
		},
	}
}

func schemaDataLakePipelineIngestionSchedules() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"frequency_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"retention_unit": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"retention_value": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"frequency_interval": {
					Type:     schema.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func schemaDataLakePipelineSnapshots() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"provider": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"created_at": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"expires_at": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"frequency_yype": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"master_key": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"mongod_version": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"replica_set_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"snapshot_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"status": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"size": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"copy_region": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"policies": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func resourceMongoDBAtlasDataLakePipelineCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	dataLakePipelineReqBody := &matlas.DataLakePipeline{
		GroupID:         projectID,
		Name:            name,
		Sink:            newDataLakePipelineSink(d),
		Source:          newDataLakePipelineSource(d),
		Transformations: newDataLakePipelineTransformation(d),
	}

	dataLakePipeline, _, err := conn.DataLakePipeline.Create(ctx, projectID, dataLakePipelineReqBody)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineCreate, err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"name":       dataLakePipeline.Name,
	}))

	return resourceMongoDBAtlasDataLakePipelineRead(ctx, d, meta)
}

func resourceMongoDBAtlasDataLakePipelineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	name := ids["name"]

	dataLakePipeline, resp, err := conn.DataLakePipeline.Get(ctx, projectID, name)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineRead, name, err))
	}

	if err := d.Set("id", dataLakePipeline.ID); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "id", name, err))
	}

	if err := d.Set("state", dataLakePipeline.State); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "state", name, err))
	}

	if err := d.Set("created_date", dataLakePipeline.CreatedDate); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "created_date", name, err))
	}

	if err := d.Set("last_updated_date", dataLakePipeline.LastUpdatedDate); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "last_updated_date", name, err))
	}

	if err := d.Set("sink", flattenDataLakePipelineSink(dataLakePipeline.Sink)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "sink", name, err))
	}

	if err := d.Set("source", flattenDataLakePipelineSource(dataLakePipeline.Source)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "source", name, err))
	}

	if err := d.Set("transformations", flattenDataLakePipelineTransformations(dataLakePipeline.Transformations)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "transformations", name, err))
	}

	snapshots, _, err := conn.DataLakePipeline.ListSnapshots(ctx, projectID, name, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineRead, name, err))
	}

	if err := d.Set("snapshots", flattenDataLakePipelineSnapshots(snapshots.Results)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "snapshots", name, err))
	}

	ingestionSchedules, _, err := conn.DataLakePipeline.ListIngestionSchedules(ctx, projectID, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineRead, name, err))
	}

	if err := d.Set("ingestion_schedules", flattenDataLakePipelineIngestionSchedules(ingestionSchedules)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "ingestion_schedules", name, err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"name":       dataLakePipeline.Name,
	}))

	return nil
}

func resourceMongoDBAtlasDataLakePipelineUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	dataLakePipelineReqBody := &matlas.DataLakePipeline{
		GroupID:         projectID,
		Name:            name,
		Sink:            newDataLakePipelineSink(d),
		Source:          newDataLakePipelineSource(d),
		Transformations: newDataLakePipelineTransformation(d),
	}

	_, _, err := conn.DataLakePipeline.Update(ctx, projectID, name, dataLakePipelineReqBody)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineUpdate, err))
	}

	return resourceMongoDBAtlasDataLakePipelineRead(ctx, d, meta)
}

func resourceMongoDBAtlasDataLakePipelineDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	name := ids["name"]

	_, err := conn.DataLakePipeline.Delete(ctx, projectID, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineDelete, name, err))
	}

	return nil
}

func resourceMongoDBAtlasDataLakePipelineImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	projectID, name, err := splitDataLakePipelineImportID(d.Id())
	if err != nil {
		return nil, err
	}

	dataLakePipeline, _, err := conn.DataLakePipeline.Get(ctx, projectID, name)
	if err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImport, name, err)
	}

	if err := d.Set("name", name); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "name", name, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "project_id", name, err)
	}

	if err := d.Set("id", dataLakePipeline.ID); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "id", name, err)
	}

	if err := d.Set("state", dataLakePipeline.State); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "state", name, err)
	}

	if err := d.Set("created_date", dataLakePipeline.CreatedDate); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "created_date", name, err)
	}

	if err := d.Set("last_updated_date", dataLakePipeline.LastUpdatedDate); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "last_updated_date", name, err)
	}

	if err := d.Set("sink", flattenDataLakePipelineSink(dataLakePipeline.Sink)); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "sink", name, err)
	}

	if err := d.Set("source", flattenDataLakePipelineSource(dataLakePipeline.Source)); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "source", name, err)
	}

	if err := d.Set("transformations", flattenDataLakePipelineTransformations(dataLakePipeline.Transformations)); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "transformations", name, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"name":       dataLakePipeline.Name,
	}))

	snapshots, _, err := conn.DataLakePipeline.ListSnapshots(ctx, projectID, name, nil)
	if err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImport, name, err)
	}

	if err := d.Set("snapshots", flattenDataLakePipelineSnapshots(snapshots.Results)); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "snapshots", name, err)
	}

	ingestionSchedules, _, err := conn.DataLakePipeline.ListIngestionSchedules(ctx, projectID, name)
	if err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImport, name, err)
	}

	if err := d.Set("ingestion_schedules", flattenDataLakePipelineIngestionSchedules(ingestionSchedules)); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "ingestion_schedules", name, err)
	}

	return []*schema.ResourceData{d}, nil
}

func splitDataLakePipelineImportID(id string) (projectID, name string, err error) {
	var parts = strings.Split(id, "--")

	if len(parts) != 2 {
		err = errors.New("import format error: to import a Data Lake, use the format {project_id}--{name}")
		return
	}

	projectID = parts[0]
	name = parts[1]

	return
}

func newDataLakePipelineSink(d *schema.ResourceData) *matlas.DataLakePipelineSink {
	if sink, ok := d.Get("sink").([]interface{}); ok && len(sink) == 1 {
		sinkMap := sink[0].(map[string]interface{})
		dataLakePipelineSink := &matlas.DataLakePipelineSink{}

		if sinkType, ok := sinkMap["type"].(string); ok {
			dataLakePipelineSink.Type = sinkType
		}

		if provider, ok := sinkMap["provider"].(string); ok {
			dataLakePipelineSink.MetadataProvider = provider
		}

		if region, ok := sinkMap["region"].(string); ok {
			dataLakePipelineSink.MetadataRegion = region
		}

		dataLakePipelineSink.PartitionFields = newDataLakePipelinePartitionField(sinkMap)
		return dataLakePipelineSink
	}

	return nil
}

func newDataLakePipelinePartitionField(sinkMap map[string]interface{}) []*matlas.DataLakePipelinePartitionField {
	partitionFields, ok := sinkMap["partition_fields"].([]interface{})
	if !ok || len(partitionFields) == 0 {
		return nil
	}

	fields := make([]*matlas.DataLakePipelinePartitionField, len(partitionFields))
	for i, partitionField := range partitionFields {
		fieldMap := partitionField.(map[string]interface{})
		fields[i] = &matlas.DataLakePipelinePartitionField{
			FieldName: fieldMap["field_name"].(string),
			Order:     int32(fieldMap["order"].(int)),
		}
	}

	return fields
}

func newDataLakePipelineSource(d *schema.ResourceData) *matlas.DataLakePipelineSource {
	source, ok := d.Get("source").([]interface{})
	if !ok || len(source) == 0 {
		return nil
	}

	sourceMap := source[0].(map[string]interface{})
	dataLakePipelineSource := &matlas.DataLakePipelineSource{}

	if sourceType, ok := sourceMap["type"].(string); ok {
		dataLakePipelineSource.Type = sourceType
	}

	if clusterName, ok := sourceMap["cluster_name"].(string); ok {
		dataLakePipelineSource.ClusterName = clusterName
	}

	if collectionName, ok := sourceMap["collection_name"].(string); ok {
		dataLakePipelineSource.CollectionName = collectionName
	}

	if databaseName, ok := sourceMap["database_name"].(string); ok {
		dataLakePipelineSource.DatabaseName = databaseName
	}

	if policyID, ok := sourceMap["policy_item_id"].(string); ok {
		dataLakePipelineSource.PolicyItemID = policyID
	}

	return dataLakePipelineSource
}

func newDataLakePipelineTransformation(d *schema.ResourceData) []*matlas.DataLakePipelineTransformation {
	trasformations, ok := d.Get("transformations").([]interface{})
	if !ok || len(trasformations) == 0 {
		return nil
	}

	dataLakePipelineTransformations := make([]*matlas.DataLakePipelineTransformation, len(trasformations))
	for i, trasformation := range trasformations {
		trasformationMap := trasformation.(map[string]interface{})
		dataLakeTransformation := &matlas.DataLakePipelineTransformation{}

		if transformationType, ok := trasformationMap["type"].(string); ok {
			dataLakeTransformation.Type = transformationType
		}

		if transformationField, ok := trasformationMap["field"].(string); ok {
			dataLakeTransformation.Field = transformationField
		}

		if dataLakeTransformation.Field != "" || dataLakeTransformation.Type != "" {
			dataLakePipelineTransformations[i] = dataLakeTransformation
		}
	}

	return dataLakePipelineTransformations
}

func flattenDataLakePipelineSource(atlasPipelineSource *matlas.DataLakePipelineSource) []map[string]interface{} {
	if atlasPipelineSource == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":            atlasPipelineSource.Type,
			"cluster_name":    atlasPipelineSource.ClusterName,
			"collection_name": atlasPipelineSource.CollectionName,
			"database_name":   atlasPipelineSource.DatabaseName,
			"project_id":      atlasPipelineSource.GroupID,
		},
	}
}

func flattenDataLakePipelineSink(atlasPipelineSink *matlas.DataLakePipelineSink) []map[string]interface{} {
	if atlasPipelineSink == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":             atlasPipelineSink.Type,
			"provider":         atlasPipelineSink.MetadataProvider,
			"region":           atlasPipelineSink.MetadataRegion,
			"partition_fields": flattenDataLakePipelinePartitionFields(atlasPipelineSink.PartitionFields),
		},
	}
}

func flattenDataLakePipelineIngestionSchedules(atlasPipelineIngestionSchedules []*matlas.DataLakePipelineIngestionSchedule) []map[string]interface{} {
	if len(atlasPipelineIngestionSchedules) == 0 {
		return nil
	}

	out := make([]map[string]interface{}, len(atlasPipelineIngestionSchedules))
	for i, schedule := range atlasPipelineIngestionSchedules {
		out[i] = map[string]interface{}{
			"id":                 schedule.ID,
			"frequency_type":     schedule.FrequencyType,
			"frequency_interval": schedule.FrequencyInterval,
			"retention_unit":     schedule.RetentionUnit,
			"retention_value":    schedule.RetentionValue,
		}
	}

	return out
}

func flattenDataLakePipelineSnapshots(snapshots []*matlas.DataLakePipelineSnapshot) []map[string]interface{} {
	if len(snapshots) == 0 {
		return nil
	}

	out := make([]map[string]interface{}, len(snapshots))
	for i, snapshot := range snapshots {
		out[i] = map[string]interface{}{
			"id":               snapshot.ID,
			"provider":         snapshot.CloudProvider,
			"created_at":       snapshot.CreatedAt,
			"expires_at":       snapshot.ExpiresAt,
			"frequency_yype":   snapshot.FrequencyType,
			"master_key":       snapshot.MasterKeyUUID,
			"mongod_version":   snapshot.MongodVersion,
			"replica_set_name": snapshot.ReplicaSetName,
			"type":             snapshot.Type,
			"snapshot_type":    snapshot.SnapshotType,
			"status":           snapshot.Status,
			"size":             snapshot.StorageSizeBytes,
			"policies":         snapshot.PolicyItems,
		}
	}
	return out
}

func flattenDataLakePipelineTransformations(atlasPipelineTransformation []*matlas.DataLakePipelineTransformation) []map[string]interface{} {
	if len(atlasPipelineTransformation) == 0 {
		return nil
	}

	out := make([]map[string]interface{}, len(atlasPipelineTransformation))
	for i, atlasPipelineTransformation := range atlasPipelineTransformation {
		out[i] = map[string]interface{}{
			"type":  atlasPipelineTransformation.Type,
			"field": atlasPipelineTransformation.Field,
		}
	}
	return out
}

func flattenDataLakePipelinePartitionFields(atlasDataLakePipelinePartitionFields []*matlas.DataLakePipelinePartitionField) []map[string]interface{} {
	if len(atlasDataLakePipelinePartitionFields) == 0 {
		return nil
	}

	out := make([]map[string]interface{}, len(atlasDataLakePipelinePartitionFields))
	for i, atlasDataLakePipelinePartitionField := range atlasDataLakePipelinePartitionFields {
		out[i] = map[string]interface{}{
			"field_name": atlasDataLakePipelinePartitionField.FieldName,
			"order":      atlasDataLakePipelinePartitionField.Order,
		}
	}
	return out
}
