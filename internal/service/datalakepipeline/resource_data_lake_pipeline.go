package datalakepipeline

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312001/admin"
)

const (
	errorDataLakePipelineCreate      = "error creating MongoDB Atlas DataLake Pipeline: %s"
	errorDataLakePipelineRead        = "error reading MongoDB Atlas DataLake Pipeline (%s): %s"
	errorDataLakePipelineImport      = "error importing MongoDB Atlas DataLake Pipeline (%s): %s"
	errorDataLakePipelineImportField = "error importing field (%s) for MongoDB Atlas DataLake Pipeline (%s): %s"
	errorDataLakePipelineDelete      = "error deleting MongoDB Atlas DataLake Pipeline (%s): %s"
	errorDataLakePipelineUpdate      = "error updating MongoDB Atlas DataLake Pipeline: %s"
	errorDataLakePipelineSetting     = "error setting `%s` for MongoDB Atlas DataLake Pipeline (%s): %s"
	ErrorDataLakeSetting             = "error setting `%s` for MongoDB Atlas DataLake (%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Data Lake is deprecated. As of September 2024, Data Lake is deprecated and will reach end-of-life. To learn more, see https://dochub.mongodb.org/core/data-lake-deprecation",
		CreateContext:      resourceCreate,
		ReadContext:        resourceRead,
		UpdateContext:      resourceUpdate,
		DeleteContext:      resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImport,
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
			"snapshots":           schemaSnapshots(),
			"ingestion_schedules": schemaSchedules(),
		},
	}
}

func schemaSchedules() *schema.Schema {
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

func schemaSnapshots() *schema.Schema {
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

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	params := &admin.DataLakeIngestionPipeline{
		GroupId:         conversion.StringPtr(projectID),
		Name:            conversion.StringPtr(name),
		Sink:            newSink(d),
		Source:          newSource(d),
		Transformations: newTransformation(d),
	}

	pipeline, _, err := connV2.DataLakePipelinesApi.CreatePipeline(ctx, projectID, params).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineCreate, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"name":       pipeline.GetName(),
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	name := ids["name"]

	pipeline, resp, err := connV2.DataLakePipelinesApi.GetPipeline(ctx, projectID, name).Execute()
	if validate.StatusNotFound(resp) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineRead, name, err))
	}

	if err := d.Set("id", pipeline.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorDataLakeSetting, "id", name, err))
	}

	if err := d.Set("state", pipeline.GetState()); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "state", name, err))
	}

	if err := d.Set("created_date", conversion.TimePtrToStringPtr(pipeline.CreatedDate)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "created_date", name, err))
	}

	if err := d.Set("last_updated_date", conversion.TimePtrToStringPtr(pipeline.LastUpdatedDate)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "last_updated_date", name, err))
	}

	if err := d.Set("sink", flattenSink(pipeline.Sink)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "sink", name, err))
	}

	if err := d.Set("source", flattenSource(pipeline.Source)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "source", name, err))
	}

	if err := d.Set("transformations", flattenTransformations(pipeline.GetTransformations())); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "transformations", name, err))
	}

	snapshots, _, err := connV2.DataLakePipelinesApi.ListPipelineSnapshots(ctx, projectID, name).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineRead, name, err))
	}

	if err := d.Set("snapshots", flattenSnapshots(snapshots.GetResults())); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "snapshots", name, err))
	}

	ingestionSchedules, _, err := connV2.DataLakePipelinesApi.ListPipelineSchedules(ctx, projectID, name).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineRead, name, err))
	}

	if err := d.Set("ingestion_schedules", flattenIngestionSchedules(ingestionSchedules)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineSetting, "ingestion_schedules", name, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"name":       pipeline.GetName(),
	}))

	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	params := &admin.DataLakeIngestionPipeline{
		GroupId:         conversion.StringPtr(projectID),
		Name:            conversion.StringPtr(name),
		Sink:            newSink(d),
		Source:          newSource(d),
		Transformations: newTransformation(d),
	}

	_, _, err := connV2.DataLakePipelinesApi.UpdatePipeline(ctx, projectID, name, params).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineUpdate, err))
	}
	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	name := ids["name"]

	_, _, err := connV2.DataLakePipelinesApi.DeletePipeline(ctx, projectID, name).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakePipelineDelete, name, err))
	}

	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID, name, err := splitDataLakePipelineImportID(d.Id())
	if err != nil {
		return nil, err
	}

	pipeline, _, err := connV2.DataLakePipelinesApi.GetPipeline(ctx, projectID, name).Execute()
	if err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImport, name, err)
	}

	if err := d.Set("name", name); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "name", name, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "project_id", name, err)
	}

	if err := d.Set("id", pipeline.GetId()); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "id", name, err)
	}

	if err := d.Set("state", pipeline.GetState()); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "state", name, err)
	}

	if err := d.Set("created_date", conversion.TimePtrToStringPtr(pipeline.CreatedDate)); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "created_date", name, err)
	}

	if err := d.Set("last_updated_date", conversion.TimePtrToStringPtr(pipeline.LastUpdatedDate)); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "last_updated_date", name, err)
	}

	if err := d.Set("sink", flattenSink(pipeline.Sink)); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "sink", name, err)
	}

	if err := d.Set("source", flattenSource(pipeline.Source)); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "source", name, err)
	}

	if err := d.Set("transformations", flattenTransformations(pipeline.GetTransformations())); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "transformations", name, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"name":       pipeline.GetName(),
	}))

	snapshots, _, err := connV2.DataLakePipelinesApi.ListPipelineSnapshots(ctx, projectID, name).Execute()
	if err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImport, name, err)
	}

	if err := d.Set("snapshots", flattenSnapshots(snapshots.GetResults())); err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImportField, "snapshots", name, err)
	}

	ingestionSchedules, _, err := connV2.DataLakePipelinesApi.ListPipelineSchedules(ctx, projectID, name).Execute()
	if err != nil {
		return nil, fmt.Errorf(errorDataLakePipelineImport, name, err)
	}

	if err := d.Set("ingestion_schedules", flattenIngestionSchedules(ingestionSchedules)); err != nil {
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

func newSink(d *schema.ResourceData) *admin.IngestionSink {
	if sink, ok := d.Get("sink").([]any); ok && len(sink) == 1 {
		sinkMap := sink[0].(map[string]any)
		dataLakePipelineSink := &admin.IngestionSink{}

		if sinkType, ok := sinkMap["type"].(string); ok {
			dataLakePipelineSink.Type = conversion.StringPtr(sinkType)
		}
		if provider, ok := sinkMap["provider"].(string); ok {
			dataLakePipelineSink.MetadataProvider = conversion.StringPtr(provider)
		}
		if region, ok := sinkMap["region"].(string); ok {
			dataLakePipelineSink.MetadataRegion = conversion.StringPtr(region)
		}
		dataLakePipelineSink.PartitionFields = newPartitionField(sinkMap)
		return dataLakePipelineSink
	}
	return nil
}

func newPartitionField(sinkMap map[string]any) *[]admin.DataLakePipelinesPartitionField {
	partitionFields, ok := sinkMap["partition_fields"].([]any)
	if !ok || len(partitionFields) == 0 {
		return nil
	}
	fields := make([]admin.DataLakePipelinesPartitionField, len(partitionFields))
	for i, partitionField := range partitionFields {
		fieldMap := partitionField.(map[string]any)
		fields[i] = admin.DataLakePipelinesPartitionField{
			FieldName: fieldMap["field_name"].(string),
			Order:     fieldMap["order"].(int),
		}
	}
	return &fields
}

func newSource(d *schema.ResourceData) *admin.IngestionSource {
	source, ok := d.Get("source").([]any)
	if !ok || len(source) == 0 {
		return nil
	}

	sourceMap := source[0].(map[string]any)
	dataLakePipelineSource := new(admin.IngestionSource)

	if sourceType, ok := sourceMap["type"].(string); ok {
		dataLakePipelineSource.Type = conversion.StringPtr(sourceType)
	}

	if clusterName, ok := sourceMap["cluster_name"].(string); ok {
		dataLakePipelineSource.ClusterName = conversion.StringPtr(clusterName)
	}

	if collectionName, ok := sourceMap["collection_name"].(string); ok {
		dataLakePipelineSource.CollectionName = conversion.StringPtr(collectionName)
	}

	if databaseName, ok := sourceMap["database_name"].(string); ok {
		dataLakePipelineSource.DatabaseName = conversion.StringPtr(databaseName)
	}

	if policyID, ok := sourceMap["policy_item_id"].(string); ok {
		dataLakePipelineSource.PolicyItemId = conversion.StringPtr(policyID)
	}

	return dataLakePipelineSource
}

func newTransformation(d *schema.ResourceData) *[]admin.FieldTransformation {
	trasformations, ok := d.Get("transformations").([]any)
	if !ok || len(trasformations) == 0 {
		return nil
	}

	dataLakePipelineTransformations := make([]admin.FieldTransformation, 0)
	for _, trasformation := range trasformations {
		trasformationMap := trasformation.(map[string]any)
		dataLakeTransformation := admin.FieldTransformation{}

		if transformationType, ok := trasformationMap["type"].(string); ok {
			dataLakeTransformation.Type = conversion.StringPtr(transformationType)
		}

		if transformationField, ok := trasformationMap["field"].(string); ok {
			dataLakeTransformation.Field = conversion.StringPtr(transformationField)
		}

		if conversion.SafeString(dataLakeTransformation.Field) != "" || conversion.SafeString(dataLakeTransformation.Type) != "" {
			dataLakePipelineTransformations = append(dataLakePipelineTransformations, dataLakeTransformation)
		}
	}
	return &dataLakePipelineTransformations
}

func flattenSource(source *admin.IngestionSource) []map[string]any {
	if source == nil {
		return nil
	}
	return []map[string]any{
		{
			"type":            source.GetType(),
			"cluster_name":    source.GetClusterName(),
			"collection_name": source.GetCollectionName(),
			"database_name":   source.GetDatabaseName(),
			"project_id":      source.GetGroupId(),
		},
	}
}

func flattenSink(sink *admin.IngestionSink) []map[string]any {
	if sink == nil {
		return nil
	}
	return []map[string]any{
		{
			"type":             sink.GetType(),
			"provider":         sink.GetMetadataProvider(),
			"region":           sink.GetMetadataRegion(),
			"partition_fields": flattenPartitionFields(sink.GetPartitionFields()),
		},
	}
}

func flattenIngestionSchedules(schedules []admin.DiskBackupApiPolicyItem) []map[string]any {
	if len(schedules) == 0 {
		return nil
	}
	out := make([]map[string]any, len(schedules))
	for i, schedule := range schedules {
		out[i] = map[string]any{
			"id":                 schedule.GetId(),
			"frequency_type":     schedule.GetFrequencyType(),
			"frequency_interval": schedule.GetFrequencyInterval(),
			"retention_unit":     schedule.GetRetentionUnit(),
			"retention_value":    schedule.GetRetentionValue(),
		}
	}
	return out
}

func flattenSnapshots(snapshots []admin.DiskBackupSnapshot) []map[string]any {
	if len(snapshots) == 0 {
		return nil
	}
	out := make([]map[string]any, len(snapshots))
	for i := range snapshots {
		snapshot := &snapshots[i]
		out[i] = map[string]any{
			"id":               snapshot.GetId(),
			"provider":         snapshot.GetCloudProvider(),
			"created_at":       conversion.TimePtrToStringPtr(snapshot.CreatedAt),
			"expires_at":       conversion.TimePtrToStringPtr(snapshot.ExpiresAt),
			"frequency_yype":   snapshot.GetFrequencyType(),
			"master_key":       snapshot.GetMasterKeyUUID(),
			"mongod_version":   snapshot.GetMongodVersion(),
			"replica_set_name": snapshot.GetReplicaSetName(),
			"type":             snapshot.GetType(),
			"snapshot_type":    snapshot.GetSnapshotType(),
			"status":           snapshot.GetStatus(),
			"size":             snapshot.GetStorageSizeBytes(),
			"policies":         snapshot.GetPolicyItems(),
		}
	}
	return out
}

func flattenTransformations(transformations []admin.FieldTransformation) []map[string]any {
	if len(transformations) == 0 {
		return nil
	}
	out := make([]map[string]any, len(transformations))
	for i, transformation := range transformations {
		out[i] = map[string]any{
			"type":  transformation.GetType(),
			"field": transformation.GetField(),
		}
	}
	return out
}

func flattenPartitionFields(fields []admin.DataLakePipelinesPartitionField) []map[string]any {
	if len(fields) == 0 {
		return nil
	}
	out := make([]map[string]any, len(fields))
	for i, field := range fields {
		out[i] = map[string]any{
			"field_name": field.GetFieldName(),
			"order":      field.GetOrder(),
		}
	}
	return out
}
