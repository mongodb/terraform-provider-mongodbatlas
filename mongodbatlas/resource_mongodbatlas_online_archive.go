package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorOnlineArchivesCreate = "error creating MongoDB Online Archive: %s"
	errorOnlineArchivesDelete = "error deleting MongoDB Online Archive: %s atlas_id (%s)"
)

func resourceMongoDBAtlasOnlineArchive() *schema.Resource {
	return &schema.Resource{
		Schema: getMongoDBAtlasOnlineArchiveSchema(),
		Create: resourceMongoDBAtlasOnlineArchiveCreate,
		Read:   resourceMongoDBAtlasOnlineArchiveRead,
		Delete: resourceMongoDBAtlasOnlineArchiveDelete,
		Update: resourceMongoDBAtlasOnlineArchiveUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasOnlineArchiveImportState,
		},
	}
}

// https://docs.atlas.mongodb.com/reference/api/online-archive-create-one
func getMongoDBAtlasOnlineArchiveSchema() map[string]*schema.Schema {
	criteriaValidator := func(val interface{}, key string) (warns []string, errs []error) {
		in := val.(map[string]interface{})

		_type, ok := in["type"]
		_, dateFieldOk := in["date_field"]
		_, expiredOk := in["expire_after_days"]
		_, queryOk := in["query"]

		if !ok {
			return
		}

		if _type == "DATE" {
			if !dateFieldOk {
				errs = append(errs, fmt.Errorf("error: criteria.date_field is required for DATE type"))
			}

			if !expiredOk {
				errs = append(errs, fmt.Errorf("error: criteria.expire_after_days is required for DATE type"))
			}
		}

		if _type == "CUSTOM" {
			if !queryOk {
				errs = append(errs, fmt.Errorf("error: criteria.query is required for CUSTOM type"))
			}
		}

		return
	}

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
		"db_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"criteria": {
			Type:     schema.TypeMap,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"DATE, CUSTOM"}, false),
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
						Type:     schema.TypeFloat,
						Optional: true,
					},
					"query": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
			ValidateFunc: criteriaValidator,
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
						Type:         schema.TypeFloat,
						Required:     true,
						ValidateFunc: validation.FloatAtLeast(0.0),
					},
					"field_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		// mongodb_atlas id
		"atlas_id": {
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
	}
}

func resourceMongoDBAtlasOnlineArchiveCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)

	inputRequest := mapToArchivePayload(d)
	outputRequest, _, err := conn.OnlineArchives.Create(context.Background(), projectID, inputRequest.ClusterName, &inputRequest)

	if err != nil {
		return fmt.Errorf(errorOnlineArchivesCreate, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": inputRequest.ClusterName,
		"atlas_id":     outputRequest.ID,
	}))

	return resourceMongoDBAtlasOnlineArchiveRead(d, meta)
}

func resourceMongoDBAtlasOnlineArchiveRead(d *schema.ResourceData, meta interface{}) error {
	// getting the atlas id
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())

	atlasID := ids["atlas_id"]
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	outOnlineArchive, _, err := conn.OnlineArchives.Get(context.Background(), projectID, clusterName, atlasID)

	if err != nil {
		reset := strings.Contains(err.Error(), "404") && !d.IsNewResource()
		if reset {
			d.SetId("")
			return nil
		}
	}

	return syncSchema(d, outOnlineArchive)
}

func resourceMongoDBAtlasOnlineArchiveDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	atlasID := ids["atlas_id"]
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	_, err := conn.OnlineArchives.Delete(context.Background(), projectID, clusterName, atlasID)

	if err != nil {
		alreadyDeleted := strings.Contains(err.Error(), "404") && !d.IsNewResource()
		if alreadyDeleted {
			return nil
		}

		return fmt.Errorf(errorOnlineArchivesDelete, err, atlasID)
	}
	return nil
}

func resourceMongoDBAtlasOnlineArchiveImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)
	parts := strings.Split(d.Id(), "-")

	if len(parts) != 3 {
		return nil, errors.New("import format error to import a Mongo Online Archive, use the format {project_id}-{cluste_rname}-{atlas_id} ")
	}

	projectID, clusterName, atlasID := parts[0], parts[1], parts[2]

	outOnlineArchive, _, err := conn.OnlineArchives.Get(context.Background(), projectID, clusterName, atlasID)

	if err != nil {
		return nil, fmt.Errorf("could not import Online Archive %s in project %s, error %s", atlasID, projectID, err.Error())
	}

	// soft error, because after the import will be a read execution
	if err := d.Set("project_id", projectID); err != nil {
		log.Printf("error setting project id %s for Online Archive id: %s", err, atlasID)
	}

	d.SetId(encodeStateID(map[string]string{
		"atlas_id":     outOnlineArchive.ID,
		"cluster_name": outOnlineArchive.ClusterName,
		"project_id":   projectID,
	}))

	return []*schema.ResourceData{d}, nil
}

func mapToArchivePayload(d *schema.ResourceData) matlas.OnlineArchive {
	// shared input
	requestInput := matlas.OnlineArchive{
		ClusterName: d.Get("cluster_name").(string),
		DBName:      d.Get("db_name").(string),
		CollName:    d.Get("coll_name").(string),
	}

	requestInput.Criteria = mapCriteria(d)

	if partitions, ok := d.GetOk("partition_fields"); ok {
		list := partitions.([]interface{})

		if len(list) > 0 {
			partitionList := make([]*matlas.PartitionFields, 0, len(list))
			for _, partition := range list {
				item := partition.(map[string]interface{})

				partitionList = append(partitionList,
					&matlas.PartitionFields{
						FieldName: item["field_name"].(string),
						Order:     pointy.Float64(item["order"].(float64)),
					},
				)
			}

			requestInput.PartitionFields = partitionList
		}
	}

	return requestInput
}

func resourceMongoDBAtlasOnlineArchiveUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	ids := decodeStateID(d.Id())

	atlasID := ids["atlas_id"]
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	// if the criteria or the paused is enable then perform an update
	paused := d.HasChange("paused")
	criteria := d.HasChange("criteria")

	// nothing to do, let's go
	if !paused && !criteria {
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

	_, _, err := conn.OnlineArchives.Update(context.Background(), projectID, clusterName, atlasID, &request)

	if err != nil {
		return fmt.Errorf("error updating Mongo Online Archive id: %s %s", atlasID, err.Error())
	}

	return resourceMongoDBAtlasOnlineArchiveRead(d, meta)
}

func syncSchema(d *schema.ResourceData, in *matlas.OnlineArchive) error {
	// computed attribute
	schemaVals := map[string]interface{}{
		"cluster_name": in.ClusterName,
		"atlas_id":     in.ID,
		"paused":       in.Paused,
		"state":        in.State,
		"coll_name":    in.CollName,
	}

	// criteria
	criteria := map[string]interface{}{
		"type":              in.Criteria.Type,
		"date_field":        in.Criteria.DateField,
		"date_format":       in.Criteria.DateFormat,
		"expire_after_days": in.Criteria.ExpireAfterDays,
		// missing query check in client
	}

	// clean up criteria for empty values
	for key, val := range criteria {
		if isEmpty(val) {
			delete(criteria, key)
		}
	}

	schemaVals["criteria"] = criteria

	// partitions fields
	if len(in.PartitionFields) > 0 {
		expected := make([]map[string]interface{}, 0, len(in.PartitionFields))
		for _, partition := range in.PartitionFields {
			if partition == nil {
				continue
			}

			partition := map[string]interface{}{
				"field_name": partition.FieldName,
				"field_type": partition.FieldType,
				"order":      partition.Order,
			}

			expected = append(expected, partition)
		}
		schemaVals["partition_fields"] = expected
	}

	return nil
}

func mapCriteria(d *schema.ResourceData) *matlas.OnlineArchiveCriteria {
	criteria := d.Get("criteria").(map[string]interface{})

	criteriaInput := &matlas.OnlineArchiveCriteria{
		Type: criteria["type"].(string),
	}

	if criteriaInput.Type == "DATE" {
		criteriaInput.DateField = criteria["date_field"].(string)
		criteriaInput.ExpireAfterDays = criteria["expire_after_days"].(float64)
		// optional
		if dformat, ok := criteria["date_format"]; ok {
			if len(dformat.(string)) > 0 {
				criteriaInput.DateFormat = dformat.(string)
			}
		}
	}

	// Pending update client missing QUERY field
	// if criteriaInput.Type == "CUSTOM" {}
	return criteriaInput
}

func isEmpty(val interface{}) bool {
	if val == nil {
		return true
	}

	switch v := val.(type) {
	case *bool:
		if v == nil {
			return true
		}
	case *float64:
		if v == nil {
			return true
		}
	case *int64:
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
