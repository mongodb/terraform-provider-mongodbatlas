package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorOnlineArchivesCreate = "error creating MongoDB Atlas Online Archive:: %s"
	errorOnlineArchivesDelete = "error deleting MongoDB Atlas Online Archive: %s archive_id (%s)"
)

func resourceMongoDBAtlasOnlineArchive() *schema.Resource {
	return &schema.Resource{
		Schema:        getMongoDBAtlasOnlineArchiveSchema(),
		CreateContext: resourceMongoDBAtlasOnlineArchiveCreate,
		ReadContext:   resourceMongoDBAtlasOnlineArchiveRead,
		DeleteContext: resourceMongoDBAtlasOnlineArchiveDelete,
		UpdateContext: resourceMongoDBAtlasOnlineArchiveUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasOnlineArchiveImportState,
		},
	}
}

// https://docs.atlas.mongodb.com/reference/api/online-archive-create-one
func getMongoDBAtlasOnlineArchiveSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// argument values
		"project_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"cluster_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"coll_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"collection_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice([]string{"STANDARD", "TIMESERIES"}, false),
		},
		"db_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"criteria": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"DATE", "CUSTOM"}, false),
					},
					"date_field": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"date_format": {
						Type:         schema.TypeString,
						Optional:     true,
						Computed:     true, // api will set the default
						ValidateFunc: validation.StringInSlice([]string{"ISODATE", "EPOCH_SECONDS", "EPOCH_MILLIS", "EPOCH_NANOSECONDS"}, false),
					},
					"expire_after_days": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"query": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"schedule": {
			Type:     schema.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"DAILY", "MONTHLY", "WEEKLY", "DEFAULT"}, false),
					},
					"end_hour": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"end_minute": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"start_hour": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"start_minute": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"day_of_month": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"day_of_week": {
						Type:     schema.TypeInt,
						Optional: true,
					},
				},
			},
		},
		"partition_fields": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"field_name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"order": {
						Type:         schema.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"field_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"archive_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"paused": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"sync_creation": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func resourceMongoDBAtlasOnlineArchiveCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	inputRequest := mapToArchivePayload(d)
	outputRequest, _, err := conn.OnlineArchives.Create(ctx, projectID, clusterName, &inputRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorOnlineArchivesCreate, err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": outputRequest.ClusterName,
		"archive_id":   outputRequest.ID,
	}))

	if d.Get("sync_creation").(bool) {
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"PENDING", "ARCHIVING", "PAUSING", "PAUSED", "ORPHANED", "REPEATING"},
			Target:     []string{"IDLE", "ACTIVE"},
			Refresh:    resourceOnlineRefreshFunc(ctx, projectID, outputRequest.ClusterName, outputRequest.ID, conn),
			Timeout:    3 * time.Hour,
			MinTimeout: 1 * time.Minute,
			Delay:      3 * time.Minute,
		}

		// Wait, catching any errors
		_, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating the online archive status %s for cluster %s", outputRequest.ClusterName, outputRequest.ID))
		}
	}

	return resourceMongoDBAtlasOnlineArchiveRead(ctx, d, meta)
}

func resourceOnlineRefreshFunc(ctx context.Context, projectID, clusterName, archiveID string, client *matlas.Client) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, resp, err := client.OnlineArchives.Get(ctx, projectID, clusterName, archiveID)

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

		if c.State != "" {
			log.Printf("[DEBUG] status for MongoDB archive_id: %s: %s", archiveID, c.State)
		}

		return c, c.State, nil
	}
}

func resourceMongoDBAtlasOnlineArchiveRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())

	atlasID := ids["archive_id"]
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	onlineArchive, resp, err := conn.OnlineArchives.Get(context.Background(), projectID, clusterName, atlasID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error MongoDB Atlas Online Archive with id %s, read error %s", atlasID, err.Error()))
	}

	mapValues := fromOnlineArchiveToMapInCreate(onlineArchive)

	for key, val := range mapValues {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(fmt.Errorf("error MongoDB Atlas Online Archive with id %s, read error %s", atlasID, err.Error()))
		}
	}
	return nil
}

func resourceMongoDBAtlasOnlineArchiveDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	atlasID := ids["archive_id"]
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	_, err := conn.OnlineArchives.Delete(ctx, projectID, clusterName, atlasID)

	if err != nil {
		alreadyDeleted := strings.Contains(err.Error(), "404") && !d.IsNewResource()
		if alreadyDeleted {
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorOnlineArchivesDelete, err, atlasID))
	}
	return nil
}

func resourceMongoDBAtlasOnlineArchiveImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas
	parts := strings.Split(d.Id(), "-")

	var projectID, clusterName, atlasID string

	if len(parts) != 3 {
		if len(parts) < 3 {
			return nil, errors.New("import format error to import a MongoDB Atlas Online Archive, use the format {project_id}-{cluster_name}-{archive_id}")
		}

		projectID = parts[0]
		clusterName = strings.Join(parts[1:len(parts)-2], "")
		atlasID = parts[len(parts)-1]
	} else {
		projectID, clusterName, atlasID = parts[0], parts[1], parts[2]
	}

	outOnlineArchive, _, err := conn.OnlineArchives.Get(ctx, projectID, clusterName, atlasID)

	if err != nil {
		return nil, fmt.Errorf("could not import Online Archive %s in project %s, error %s", atlasID, projectID, err.Error())
	}

	// soft error, because after the import will be a read execution
	if err := d.Set("project_id", projectID); err != nil {
		log.Printf("error setting project id %s for Online Archive id: %s", err, atlasID)
	}

	mapValues := fromOnlineArchiveToMapInCreate(outOnlineArchive)
	for key, val := range mapValues {
		if err := d.Set(key, val); err != nil {
			return nil, fmt.Errorf("error MongoDB Atlas Online Archive with id %s, read error: %w", atlasID, err)
		}
	}

	d.SetId(encodeStateID(map[string]string{
		"archive_id":   outOnlineArchive.ID,
		"cluster_name": outOnlineArchive.ClusterName,
		"project_id":   projectID,
	}))

	return []*schema.ResourceData{d}, nil
}

func mapToArchivePayload(d *schema.ResourceData) matlas.OnlineArchive {
	// shared input
	requestInput := matlas.OnlineArchive{
		DBName:         d.Get("db_name").(string),
		CollName:       d.Get("coll_name").(string),
		CollectionType: d.Get("collection_type").(string),
	}

	requestInput.Criteria = mapCriteria(d)
	requestInput.Schedule = mapSchedule(d)

	if partitions, ok := d.GetOk("partition_fields"); ok {
		list := partitions.([]interface{})

		if len(list) > 0 {
			partitionList := make([]*matlas.PartitionFields, 0, len(list))
			for _, partition := range list {
				item := partition.(map[string]interface{})
				localOrder := item["order"].(int)
				localOrderFloat := float64(localOrder)

				query := &matlas.PartitionFields{
					FieldName: item["field_name"].(string),
					Order:     pointy.Float64(localOrderFloat),
				}

				if dbType, ok := item["field_type"]; ok && dbType != nil {
					if dbType.(string) != "" {
						query.FieldType = dbType.(string)
					}
				}

				partitionList = append(partitionList, query)
			}

			requestInput.PartitionFields = partitionList
		}
	}

	return requestInput
}

func resourceMongoDBAtlasOnlineArchiveUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	ids := decodeStateID(d.Id())

	atlasID := ids["archive_id"]
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	// if the criteria or the paused is enable then perform an update
	paused := d.HasChange("paused")
	criteria := d.HasChange("criteria")
	schedule := d.HasChange("schedule")

	collectionType := d.HasChange("collection_type")

	// nothing to do, let's go
	if !paused && !criteria && !collectionType && !schedule {
		return nil
	}

	request := matlas.OnlineArchive{}

	// reading current value
	if paused {
		request.Paused = pointy.Bool(d.Get("paused").(bool))
	}

	if criteria {
		request.Criteria = mapCriteria(d)
	}

	if schedule {
		request.Schedule = mapSchedule(d)
	}

	if collectionType {
		request.CollectionType = d.Get("collection_type").(string)
	}

	_, _, err := conn.OnlineArchives.Update(ctx, projectID, clusterName, atlasID, &request)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating Mongo Online Archive id: %s %s", atlasID, err.Error()))
	}

	return resourceMongoDBAtlasOnlineArchiveRead(ctx, d, meta)
}

func fromOnlineArchiveToMap(in *matlas.OnlineArchive) map[string]interface{} {
	// computed attribute
	schemaVals := map[string]interface{}{
		"cluster_name":    in.ClusterName,
		"archive_id":      in.ID,
		"paused":          in.Paused,
		"state":           in.State,
		"coll_name":       in.CollName,
		"collection_type": in.CollectionType,
	}

	criteria := map[string]interface{}{
		"type":        in.Criteria.Type,
		"date_field":  in.Criteria.DateField,
		"date_format": in.Criteria.DateFormat,
		"query":       in.Criteria.Query,
	}

	schedule := map[string]interface{}{
		"type":         in.Schedule.Type,
		"day_of_month": in.Schedule.DayOfMonth,
		"day_of_week":  in.Schedule.DayOfWeek,
		"end_hour":     in.Schedule.EndHour,
		"end_minute":   in.Schedule.EndMinute,
		"start_hour":   in.Schedule.StartHour,
		"start_minute": in.Schedule.StartMinute,
	}

	// note: criteria is a conditional field, not required when type is equal to CUSTOM
	if in.Criteria.ExpireAfterDays != nil {
		criteria["expire_after_days"] = int(*in.Criteria.ExpireAfterDays)
	}

	// clean up criteria for empty values
	for key, val := range criteria {
		if isEmpty(val) {
			delete(criteria, key)
		}
	}

	// clean up schedule for empty values
	for key, val := range schedule {
		if isEmpty(val) {
			delete(schedule, key)
		}
	}

	schemaVals["criteria"] = []interface{}{criteria}
	schemaVals["schedule"] = []interface{}{schedule}

	// partitions fields
	if len(in.PartitionFields) == 0 {
		return schemaVals
	}

	partitionFieldsMap := make([]map[string]interface{}, 0, len(in.PartitionFields))
	for _, field := range in.PartitionFields {
		if field == nil {
			continue
		}

		fieldMap := map[string]interface{}{
			"field_name": field.FieldName,
			"field_type": field.FieldType,
			"order":      field.Order,
		}

		partitionFieldsMap = append(partitionFieldsMap, fieldMap)
	}
	schemaVals["partition_fields"] = partitionFieldsMap

	return schemaVals
}

func fromOnlineArchiveToMapInCreate(in *matlas.OnlineArchive) map[string]interface{} {
	localSchema := fromOnlineArchiveToMap(in)
	delete(localSchema, "partition_fields")
	return localSchema
}

func mapCriteria(d *schema.ResourceData) *matlas.OnlineArchiveCriteria {
	criteriaList := d.Get("criteria").([]interface{})

	criteria := criteriaList[0].(map[string]interface{})

	criteriaInput := &matlas.OnlineArchiveCriteria{
		Type: criteria["type"].(string),
	}

	if criteriaInput.Type == "DATE" {
		criteriaInput.DateField = criteria["date_field"].(string)

		conversion := criteria["expire_after_days"].(int)

		criteriaInput.ExpireAfterDays = pointy.Float64(float64(conversion))
		// optional
		if dformat, ok := criteria["date_format"]; ok {
			if len(dformat.(string)) > 0 {
				criteriaInput.DateFormat = dformat.(string)
			}
		}
	}

	if criteriaInput.Type == "CUSTOM" {
		criteriaInput.Query = criteria["query"].(string)
	}

	// Pending update client missing QUERY field
	return criteriaInput
}

func mapSchedule(d *schema.ResourceData) *matlas.OnlineArchiveSchedule {
	scheduleTFConfigList := d.Get("schedule").([]interface{})
	scheduleTFConfig := scheduleTFConfigList[0].(map[string]interface{})

	scheduleInput := &matlas.OnlineArchiveSchedule{
		Type: scheduleTFConfig["type"].(string),
	}

	if endHour, ok := scheduleTFConfig["end_hour"].(int); ok {
		scheduleInput.EndHour = pointy.Int32(int32(endHour))
	}

	if endMinute, ok := scheduleTFConfig["end_minute"].(int); ok {
		scheduleInput.EndMinute = pointy.Int32(int32(endMinute))
	}

	if startHour, ok := scheduleTFConfig["start_hour"].(int); ok {
		scheduleInput.StartHour = pointy.Int32(int32(startHour))
	}

	if startMinute, ok := scheduleTFConfig["start_minute"].(int); ok {
		scheduleInput.StartMinute = pointy.Int32(int32(startMinute))
	}

	if dayOfWeek, ok := scheduleTFConfig["day_of_week"].(int); ok {
		scheduleInput.DayOfWeek = int32(dayOfWeek)
	}

	if dayOfMonth, ok := scheduleTFConfig["day_of_month"].(int); ok {
		scheduleInput.DayOfMonth = int32(dayOfMonth)
	}

	return scheduleInput
}

func isEmpty(val interface{}) bool {
	if val == nil {
		return true
	}

	switch v := val.(type) {
	case *bool, *float64, *int64:
		if v == nil {
			return true
		}
	case string:
		return v == ""
	case *string:
		if v == nil {
			return true
		}
		return *v == ""
	}

	return false
}
