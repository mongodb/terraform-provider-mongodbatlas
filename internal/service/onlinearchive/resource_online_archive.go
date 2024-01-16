package onlinearchive

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	admin20231001002 "go.mongodb.org/atlas-sdk/v20231001002/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mwielbut/pointy"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	errorOnlineArchivesCreate = "error creating MongoDB Atlas Online Archive:: %s"
	errorOnlineArchivesDelete = "error deleting MongoDB Atlas Online Archive: %s archive_id (%s)"
	scheduleTypeDefault       = "DEFAULT"
)

func Resource() *schema.Resource {
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
		"data_expiration_rule": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"expire_after_days": {
						Type:     schema.TypeInt,
						Required: true,
					},
				},
			},
		},
		"data_process_region": {
			Type:       schema.TypeList,
			MinItems:   1,
			MaxItems:   1,
			ConfigMode: schema.SchemaConfigModeAttr,
			Optional:   true,
			Computed:   true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"region": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"cloud_provider": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
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
						ValidateFunc: validation.StringInSlice([]string{"DAILY", "MONTHLY", "WEEKLY"}, false),
					},
					"end_hour": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"end_minute": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"start_hour": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"start_minute": {
						Type:     schema.TypeInt,
						Optional: true,
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

func resourceMongoDBAtlasOnlineArchiveCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn20231001002 := meta.(*config.MongoDBClient).Atlas20231001002
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	inputRequest := mapToArchivePayload(d)
	outputRequest, _, err := conn20231001002.OnlineArchiveApi.CreateOnlineArchive(ctx, projectID, clusterName, &inputRequest).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorOnlineArchivesCreate, err))
	}

	archiveID := *outputRequest.Id

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
		"archive_id":   archiveID,
	}))

	if d.Get("sync_creation").(bool) {
		stateConf := &retry.StateChangeConf{
			Pending:    []string{"PENDING", "ARCHIVING", "PAUSING", "PAUSED", "ORPHANED", "REPEATING"},
			Target:     []string{"IDLE", "ACTIVE"},
			Refresh:    resourceOnlineRefreshFunc(ctx, projectID, clusterName, archiveID, conn20231001002),
			Timeout:    3 * time.Hour,
			MinTimeout: 1 * time.Minute,
			Delay:      3 * time.Minute,
		}

		// Wait, catching any errors
		_, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating the online archive status %s for cluster %s", clusterName, archiveID))
		}
	}

	return resourceMongoDBAtlasOnlineArchiveRead(ctx, d, meta)
}

func resourceOnlineRefreshFunc(ctx context.Context, projectID, clusterName, archiveID string, client *admin20231001002.APIClient) retry.StateRefreshFunc {
	return func() (any, string, error) {
		c, resp, err := client.OnlineArchiveApi.GetOnlineArchive(ctx, projectID, archiveID, clusterName).Execute()

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

		if conversion.IsStringPresent(c.State) {
			log.Printf("[DEBUG] status for MongoDB archive_id: %s: %s", archiveID, *c.State)
		}

		return c, conversion.SafeString(c.State), nil
	}
}

func resourceMongoDBAtlasOnlineArchiveRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn20231001002 := meta.(*config.MongoDBClient).Atlas20231001002
	ids := conversion.DecodeStateID(d.Id())

	archiveID := ids["archive_id"]
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	onlineArchive, resp, err := conn20231001002.OnlineArchiveApi.GetOnlineArchive(context.Background(), projectID, archiveID, clusterName).Execute()
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error MongoDB Atlas Online Archive with id %s, read error %s", archiveID, err.Error()))
	}

	mapValues := fromOnlineArchiveToMap(onlineArchive)

	for key, val := range mapValues {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(fmt.Errorf("error MongoDB Atlas Online Archive with id %s, read error %s", archiveID, err.Error()))
		}
	}
	return nil
}

func resourceMongoDBAtlasOnlineArchiveDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas
	ids := conversion.DecodeStateID(d.Id())
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

func resourceMongoDBAtlasOnlineArchiveImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	conn20231001002 := meta.(*config.MongoDBClient).Atlas20231001002
	parts := strings.Split(d.Id(), "-")

	var projectID, clusterName, archiveID string

	if len(parts) != 3 {
		if len(parts) < 3 {
			return nil, errors.New("import format error to import a MongoDB Atlas Online Archive, use the format {project_id}-{cluster_name}-{archive_id}")
		}

		projectID = parts[0]
		clusterName = strings.Join(parts[1:len(parts)-2], "")
		archiveID = parts[len(parts)-1]
	} else {
		projectID, clusterName, archiveID = parts[0], parts[1], parts[2]
	}

	outOnlineArchive, _, err := conn20231001002.OnlineArchiveApi.GetOnlineArchive(ctx, projectID, archiveID, clusterName).Execute()

	if err != nil {
		return nil, fmt.Errorf("could not import Online Archive %s in project %s, error %s", archiveID, projectID, err.Error())
	}

	// soft error, because after the import will be a read execution
	if err := d.Set("project_id", projectID); err != nil {
		log.Printf("error setting project id %s for Online Archive id: %s", err, archiveID)
	}

	mapValues := fromOnlineArchiveToMap(outOnlineArchive)
	for key, val := range mapValues {
		if err := d.Set(key, val); err != nil {
			return nil, fmt.Errorf("error MongoDB Atlas Online Archive with id %s, read error: %w", archiveID, err)
		}
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"archive_id":   conversion.SafeString(outOnlineArchive.Id),
		"cluster_name": conversion.SafeString(outOnlineArchive.ClusterName),
		"project_id":   projectID,
	}))

	return []*schema.ResourceData{d}, nil
}

func mapToArchivePayload(d *schema.ResourceData) admin20231001002.BackupOnlineArchiveCreate {
	// shared input
	requestInput := admin20231001002.BackupOnlineArchiveCreate{
		DbName:   d.Get("db_name").(string),
		CollName: d.Get("coll_name").(string),
	}
	if collType := d.Get("collection_type").(string); collType != "" {
		requestInput.CollectionType = &collType
	}

	requestInput.Criteria = mapCriteria(d)
	requestInput.DataExpirationRule = mapDataExpirationRule(d)
	requestInput.DataProcessRegion = mapDataProcessRegion(d)
	requestInput.Schedule = mapSchedule(d)

	if partitions, ok := d.GetOk("partition_fields"); ok {
		list := partitions.([]any)

		if len(list) > 0 {
			partitionList := make([]admin20231001002.PartitionField, 0, len(list))
			for _, partition := range list {
				item := partition.(map[string]any)
				query := admin20231001002.PartitionField{
					FieldName: item["field_name"].(string),
					Order:     item["order"].(int),
				}

				if dbType, ok := item["field_type"]; ok && dbType != nil {
					if dbType.(string) != "" {
						query.FieldType = admin20231001002.PtrString(dbType.(string))
					}
				}

				partitionList = append(partitionList, query)
			}

			requestInput.PartitionFields = partitionList
		}
	}

	return requestInput
}

func resourceMongoDBAtlasOnlineArchiveUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn20231001002 := meta.(*config.MongoDBClient).Atlas20231001002

	ids := conversion.DecodeStateID(d.Id())

	atlasID := ids["archive_id"]
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	// if the criteria or the pausedHasChange is enable then perform an update
	pausedHasChange := d.HasChange("paused")
	criteriaHasChange := d.HasChange("criteria")
	dataExpirationRuleHasChange := d.HasChange("data_expiration_rule")
	dataProcessRegionHasChange := d.HasChange("data_process_region")
	scheduleHasChange := d.HasChange("schedule")

	collectionTypeHasChange := d.HasChange("collection_type")

	// nothing to do, let's go
	if !pausedHasChange && !criteriaHasChange && !collectionTypeHasChange && !scheduleHasChange && !dataExpirationRuleHasChange && !dataProcessRegionHasChange {
		return nil
	}

	request := admin20231001002.BackupOnlineArchive{}

	// reading current value
	if pausedHasChange {
		request.Paused = pointy.Bool(d.Get("paused").(bool))
	}

	if criteriaHasChange {
		newCriteria := mapCriteria(d)
		request.Criteria = &newCriteria
	}

	if dataExpirationRuleHasChange {
		newExpirationRule := mapDataExpirationRule(d)
		if newExpirationRule == nil {
			// expiration rule has been removed from tf config, empty dataExpirationRule object needs to be sent in patch request
			request.DataExpirationRule = &admin20231001002.DataExpirationRule{}
		} else {
			request.DataExpirationRule = newExpirationRule
		}
	}

	if dataProcessRegionHasChange {
		newDataProcessRegion := mapDataProcessRegion(d)
		if newDataProcessRegion == nil {
			request.DataProcessRegion = &admin20231001002.DataProcessRegion{}
		} else {
			request.DataProcessRegion = newDataProcessRegion
		}
	}

	if scheduleHasChange {
		request.Schedule = mapSchedule(d)
	}

	if collType := d.Get("collection_type").(string); collectionTypeHasChange && collType != "" {
		request.CollectionType = admin20231001002.PtrString(collType)
	}

	_, _, err := conn20231001002.OnlineArchiveApi.UpdateOnlineArchive(ctx, projectID, atlasID, clusterName, &request).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating Mongo Online Archive id: %s %s", atlasID, err.Error()))
	}

	return resourceMongoDBAtlasOnlineArchiveRead(ctx, d, meta)
}

func fromOnlineArchiveToMap(in *admin20231001002.BackupOnlineArchive) map[string]any {
	// computed attribute
	schemaVals := map[string]any{
		"cluster_name":    in.ClusterName,
		"archive_id":      in.Id,
		"paused":          in.Paused,
		"state":           in.State,
		"coll_name":       in.CollName,
		"collection_type": in.CollectionType,
		"db_name":         in.DbName,
	}

	criteria := map[string]any{
		"type":        in.Criteria.Type,
		"date_field":  in.Criteria.DateField,
		"date_format": in.Criteria.DateFormat,
		"query":       in.Criteria.Query,
	}

	var schedule map[string]any
	// When schedule is not provided in CREATE/UPDATE the GET returns Schedule.Type = DEFAULT
	// In this case, we don't want to update the schema as there is no SCHEDULE
	if in.Schedule != nil && in.Schedule.Type != scheduleTypeDefault {
		schedule = map[string]any{
			"type":         in.Schedule.Type,
			"day_of_month": in.Schedule.DayOfMonth,
			"day_of_week":  in.Schedule.DayOfWeek,
			"end_hour":     in.Schedule.EndHour,
			"end_minute":   in.Schedule.EndMinute,
			"start_hour":   in.Schedule.StartHour,
			"start_minute": in.Schedule.StartMinute,
		}
	}

	// note: criteria is a conditional field, not required when type is equal to CUSTOM
	if in.Criteria.ExpireAfterDays != nil {
		criteria["expire_after_days"] = *in.Criteria.ExpireAfterDays
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

	schemaVals["criteria"] = []any{criteria}

	if schedule != nil {
		schemaVals["schedule"] = []any{schedule}
	}

	var dataExpirationRule map[string]any
	if in.DataExpirationRule != nil && in.DataExpirationRule.ExpireAfterDays != nil {
		dataExpirationRule = map[string]any{
			"expire_after_days": in.DataExpirationRule.ExpireAfterDays,
		}
		schemaVals["data_expiration_rule"] = []any{dataExpirationRule}
	}

	var dataProcessRegion map[string]any
	if in.DataProcessRegion != nil && (in.DataProcessRegion.CloudProvider != nil || in.DataProcessRegion.Region != nil) {
		dataProcessRegion = map[string]any{
			"cloud_provider": in.DataProcessRegion.CloudProvider,
			"region":         in.DataProcessRegion.Region,
		}
		schemaVals["data_process_region"] = []any{dataProcessRegion}
	}

	partitionFields := in.PartitionFields
	if len(partitionFields) == 0 {
		return schemaVals
	}
	partitionFieldsMap := make([]map[string]any, 0, len(partitionFields))
	for _, field := range partitionFields {
		fieldMap := map[string]any{
			"field_name": field.FieldName,
			"field_type": field.FieldType,
			"order":      field.Order,
		}

		partitionFieldsMap = append(partitionFieldsMap, fieldMap)
	}
	schemaVals["partition_fields"] = partitionFieldsMap

	return schemaVals
}

func mapDataExpirationRule(d *schema.ResourceData) *admin20231001002.DataExpirationRule {
	if dataExpireRules, ok := d.GetOk("data_expiration_rule"); ok && len(dataExpireRules.([]any)) > 0 {
		dataExpireRule := dataExpireRules.([]any)[0].(map[string]any)
		result := admin20231001002.DataExpirationRule{}
		if expireAfterDays, ok := dataExpireRule["expire_after_days"]; ok {
			result.ExpireAfterDays = pointy.Int(expireAfterDays.(int))
		}
		return &result
	}
	return nil
}

func mapDataProcessRegion(d *schema.ResourceData) *admin20231001002.DataProcessRegion {
	if dataProcessRegions, ok := d.GetOk("data_process_region"); ok && len(dataProcessRegions.([]any)) > 0 {
		dataProcessRegion := dataProcessRegions.([]any)[0].(map[string]any)
		result := admin20231001002.DataProcessRegion{}
		if cloudProvider, ok := dataProcessRegion["cloud_provider"]; ok {
			result.CloudProvider = pointy.String(cloudProvider.(string))
		}
		if region, ok := dataProcessRegion["region"]; ok {
			result.Region = pointy.String(region.(string))
		}
		return &result
	}
	return nil
}

func mapCriteria(d *schema.ResourceData) admin20231001002.Criteria {
	criteriaList := d.Get("criteria").([]any)

	criteria := criteriaList[0].(map[string]any)

	criteriaInput := admin20231001002.Criteria{
		Type: admin20231001002.PtrString(criteria["type"].(string)),
	}

	if criteriaInput.Type != nil && *criteriaInput.Type == "DATE" {
		if dateField := criteria["date_field"].(string); dateField != "" {
			criteriaInput.DateField = admin20231001002.PtrString(dateField)
		}

		criteriaInput.ExpireAfterDays = pointy.Int(criteria["expire_after_days"].(int))

		// optional
		if dformat, ok := criteria["date_format"]; ok && dformat.(string) != "" {
			criteriaInput.DateFormat = admin20231001002.PtrString(dformat.(string))
		}
	}

	if criteriaInput.Type != nil && *criteriaInput.Type == "CUSTOM" {
		if query := criteria["query"].(string); query != "" {
			criteriaInput.Query = admin20231001002.PtrString(query)
		}
	}

	// Pending update client missing QUERY field
	return criteriaInput
}

func mapSchedule(d *schema.ResourceData) *admin20231001002.OnlineArchiveSchedule {
	// scheduleInput := &matlas.OnlineArchiveSchedule{

	// We have to provide schedule.type="DEFAULT" when the schedule block is not provided or removed
	scheduleInput := &admin20231001002.OnlineArchiveSchedule{
		Type: scheduleTypeDefault,
	}

	scheduleTFConfigInterface := d.Get("schedule")
	if scheduleTFConfigInterface == nil {
		return scheduleInput
	}

	scheduleTFConfigList, ok := scheduleTFConfigInterface.([]any)
	if !ok {
		return scheduleInput
	}

	if len(scheduleTFConfigList) == 0 {
		return scheduleInput
	}

	scheduleTFConfig := scheduleTFConfigList[0].(map[string]any)
	scheduleInput = &admin20231001002.OnlineArchiveSchedule{
		Type: scheduleTFConfig["type"].(string),
	}

	if endHour, ok := scheduleTFConfig["end_hour"].(int); ok {
		scheduleInput.EndHour = pointy.Int(endHour)
	}

	if endMinute, ok := scheduleTFConfig["end_minute"].(int); ok {
		scheduleInput.EndMinute = pointy.Int(endMinute)
	}

	if startHour, ok := scheduleTFConfig["start_hour"].(int); ok {
		scheduleInput.StartHour = pointy.Int(startHour)
	}

	if startMinute, ok := scheduleTFConfig["start_minute"].(int); ok {
		scheduleInput.StartMinute = pointy.Int(startMinute)
	}

	if dayOfWeek, ok := scheduleTFConfig["day_of_week"].(int); ok && dayOfWeek != 0 { // needed to verify attribute is actually defined
		scheduleInput.DayOfWeek = pointy.Int(dayOfWeek)
	}

	if dayOfMonth, ok := scheduleTFConfig["day_of_month"].(int); ok && dayOfMonth != 0 {
		scheduleInput.DayOfMonth = pointy.Int(dayOfMonth)
	}

	return scheduleInput
}

func isEmpty(val any) bool {
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
