package cloudbackupsnapshot

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312001/admin"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluralRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"page_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"items_per_page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expires_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_key_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mongod_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_size_bytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_provider": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"members": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cloud_provider": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"replica_set_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"replica_set_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourcePluralRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	params := &admin.ListReplicaSetBackupsApiParams{
		GroupId:      d.Get("project_id").(string),
		ClusterName:  d.Get("cluster_name").(string),
		PageNum:      conversion.Pointer(d.Get("page_num").(int)),
		ItemsPerPage: conversion.Pointer(d.Get("items_per_page").(int)),
	}

	snapshots, _, err := connV2.CloudBackupsApi.ListReplicaSetBackupsWithParams(ctx, params).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting cloudProviderSnapshots information: %s", err))
	}
	shards, _, _ := connV2.CloudBackupsApi.ListShardedClusterBackups(ctx, params.GroupId, params.ClusterName).Execute()

	if err := d.Set("results", flattenCloudBackupSnapshots(snapshots.GetResults(), shards)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `results`: %s", err))
	}

	if err := d.Set("total_count", snapshots.GetTotalCount()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `total_count`: %s", err))
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenCloudBackupSnapshots(snapshots []admin.DiskBackupReplicaSet, shards *admin.PaginatedCloudBackupShardedClusterSnapshot) []map[string]any {
	if len(snapshots) == 0 {
		return nil
	}
	results := make([]map[string]any, len(snapshots))
	for i := range snapshots {
		snapshot := &snapshots[i]
		results[i] = map[string]any{
			"id":                 snapshot.GetId(),
			"created_at":         conversion.TimePtrToStringPtr(snapshot.CreatedAt),
			"expires_at":         conversion.TimePtrToStringPtr(snapshot.ExpiresAt),
			"description":        snapshot.GetDescription(),
			"master_key_uuid":    snapshot.GetMasterKeyUUID(),
			"mongod_version":     snapshot.GetMongodVersion(),
			"snapshot_type":      snapshot.GetSnapshotType(),
			"status":             snapshot.GetStatus(),
			"storage_size_bytes": snapshot.GetStorageSizeBytes(),
			"type":               snapshot.GetType(),
			"cloud_provider":     snapshot.GetCloudProvider(),
			"replica_set_name":   snapshot.GetReplicaSetName(),
		}
		if shards != nil {
			shardResults := shards.GetResults()
			for j := range shardResults {
				if shardResults[j].GetId() == snapshot.GetId() {
					results[i]["members"] = flattenCloudMembers(shardResults[j].GetMembers())
					results[i]["snapshot_ids"] = shardResults[j].GetSnapshotIds()
				}
			}
		}
	}
	return results
}
