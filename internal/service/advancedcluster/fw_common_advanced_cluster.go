package advancedcluster

import (
	"context"
	"fmt"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
)

var defaultLabel = matlas.Label{Key: "Infrastructure Tool", Value: "MongoDB Atlas Terraform Provider"}

func ClusterDSAdvancedConfigurationListAttr() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"default_read_concern": schema.StringAttribute{
					Computed: true,
				},
				"default_write_concern": schema.StringAttribute{
					Computed: true,
				},
				"fail_index_key_too_long": schema.BoolAttribute{
					Computed: true,
				},
				"javascript_enabled": schema.BoolAttribute{
					Computed: true,
				},
				"minimum_enabled_tls_protocol": schema.StringAttribute{
					Computed: true,
				},
				"no_table_scan": schema.BoolAttribute{
					Computed: true,
				},
				"oplog_size_mb": schema.Int64Attribute{
					Computed: true,
				},
				"sample_size_bi_connector": schema.Int64Attribute{
					Computed: true,
				},
				"sample_refresh_interval_bi_connector": schema.Int64Attribute{
					Computed: true,
				},
				"oplog_min_retention_hours": schema.Int64Attribute{
					Computed: true,
				},
				"transaction_lifetime_limit_seconds": schema.Int64Attribute{
					Computed: true,
				},
			},
		},
	}
}

func ClusterDSBiConnectorConfigListAttr() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"enabled": schema.BoolAttribute{
					Computed: true,
				},
				"read_preference": schema.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func ClusterDSLabelsSetAttr() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Computed:           true,
		DeprecationMessage: fmt.Sprintf(constant.DeprecationParamByDateWithReplacement, "September 2024", "tags"),
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"key": schema.StringAttribute{
					Computed: true,
				},
				"value": schema.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func ClusterDSTagsSetAttr() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"key": schema.StringAttribute{
					Computed: true,
				},
				"value": schema.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func RemoveDefaultLabel(labels []TfLabelModel) []TfLabelModel {
	result := make([]TfLabelModel, 0)

	for _, item := range labels {
		if item.Key.ValueString() == defaultLabel.Key && item.Value.ValueString() == defaultLabel.Value {
			continue
		}
		result = append(result, item)
	}

	return result
}

func NewTfAdvancedConfigurationModelDSFromAtlas(ctx context.Context, conn *matlas.Client, projectID, clusterName string) ([]*TfAdvancedConfigurationModel, error) {
	processArgs, _, err := conn.Clusters.GetProcessArgs(ctx, projectID, clusterName)
	if err != nil {
		return nil, err
	}

	advConfigModel := NewTfAdvancedConfigurationModel(processArgs)
	return advConfigModel, err
}
