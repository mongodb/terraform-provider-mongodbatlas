package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorOnlineArchivesCreate  = "error creating MongoDB Online Archive: %s"
	errorOnlineMissingComputed = "error MongoDB Online Archive missing: %s"
)

func resourceMongoDBAtlasOnlineArchive() *schema.Resource {
	return &schema.Resource{
		Schema: getMongoDBAtlasOnlineArchiveSchema(),
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

	if err = syncSchema(d, &inputRequest, outputRequest); err != nil {
		return fmt.Errorf(errorOnlineArchivesCreate, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": inputRequest.ClusterName,
		"atlas_id":     outputRequest.ID,
	}))

	return nil
}

func mapToArchivePayload(d *schema.ResourceData) matlas.OnlineArchive {
	// shared input
	requestInput := matlas.OnlineArchive{
		ClusterName: d.Get("cluster_name").(string),
		DBName:      d.Get("db_name").(string),
		CollName:    d.Get("coll_name").(string),
	}

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
	if criteriaInput.Type == "CUSTOM" {
	}

	requestInput.Criteria = criteriaInput

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

func syncSchema(d *schema.ResourceData, in, out *matlas.OnlineArchive) error {
	// computed attribute
	if err := d.Set("atlas_id", out.ID); err != nil {
		return fmt.Errorf(errorOnlineMissingComputed, err)
	}

	if err := d.Set("paused", out.Paused); err != nil {
		return fmt.Errorf(errorOnlineMissingComputed, err)
	}

	if err := d.Set("state", out.State); err != nil {
		return fmt.Errorf(errorOnlineMissingComputed, err)
	}

	return nil
}
