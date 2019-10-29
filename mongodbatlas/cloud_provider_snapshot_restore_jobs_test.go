package mongodbatlas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

func TestCloudProviderSnapshotRestoreJobs_List(t *testing.T) {
	setup()
	defer teardown()

	requestParameters := &SnapshotReqPathParameters{
		GroupID:     "5b6212af90dc76637950a2c6",
		ClusterName: "MyCluster",
	}

	path := fmt.Sprintf("/groups/%s/clusters/%s/backup/restoreJobs", requestParameters.GroupID, requestParameters.ClusterName)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"links": [
				{
					"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0/backup/restoreJobs?pageNum=1&itemsPerPage=100",
					"rel": "self"
				}
			],
			"results": [
				{
					"cancelled": false,
					"deliveryType": "automated",
					"expired": false,
					"expiresAt": "2018-08-02T02:08:48Z",
					"id": "5b622f7087d9d6039fafe03f",
					"links": [
						{
							"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0/backup/restoreJobs/5b622f7087d9d6039fafe03f",
							"rel": "self"
						},
						{
							"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0/backup/snapshots/5b6211ff87d9d663c59d3dee",
							"rel": "http://mms.mongodb.com/snapshot"
						},
						{
							"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0",
							"rel": "http://mms.mongodb.com/cluster"
						},
						{
							"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6",
							"rel": "http://mms.mongodb.com/group"
						}
					],
					"snapshotId": "5b6211ff87d9d663c59d3feg",
					"targetClusterName": "MyOtherCluster",
					"targetGroupId": "5b6212af90dc76637950a2c6",
					"timestamp": "2018-08-01T20:02:07Z"
				},
				{
					"cancelled": false,
					"createdAt": "2018-08-01T22:05:41Z",
					"deliveryType": "download",
					"deliveryUrl": ["https://restore.mongodb.net:27017/restore-5b622e3587d9d6039faf713f.tar.gz"],
					"expired": false,
					"expiresAt": "2018-08-02T02:03:33Z",
					"id": "5b622e3587d9d6039faf713f",
					"links": [
						{
							"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0/backup/restoreJobs/5b622e3587d9d6039faf713f",
							"rel": "self"
						},
						{
							"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0/backup/snapshots/5b6211ff87d9d663c59d3dee",
							"rel": "http://mms.mongodb.com/snapshot"
						},
						{
							"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0",
							"rel": "http://mms.mongodb.com/cluster"
						},
						{
							"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6",
							"rel": "http://mms.mongodb.com/group"
						}
					],
					"snapshotId": "5b6211ff87d9d663c59d3feg",
					"timestamp": "2018-08-01T20:02:07Z"
				}
			],
			"totalCount": 2
		}`)
	})

	cloudProviderSnapshots, _, err := client.CloudProviderSnapshotRestoreJobs.List(ctx, requestParameters)
	if err != nil {
		t.Errorf("CloudProviderSnapshotRestoreJobs.List returned error: %v", err)
	}

	expected := &CloudProviderSnapshotRestoreJobs{
		Links: []*Link{
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0/backup/restoreJobs?pageNum=1&itemsPerPage=100",
				Rel:  "self",
			},
		},
		Results: []*CloudProviderSnapshotRestoreJob{
			{
				Cancelled:    false,
				DeliveryType: "automated",
				Expired:      false,
				ExpiresAt:    "2018-08-02T02:08:48Z",
				ID:           "5b622f7087d9d6039fafe03f",
				Links: []*Link{
					{
						Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0/backup/restoreJobs/5b622f7087d9d6039fafe03f",
						Rel:  "self",
					},
					{
						Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0/backup/snapshots/5b6211ff87d9d663c59d3dee",
						Rel:  "http://mms.mongodb.com/snapshot",
					},
					{
						Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0",
						Rel:  "http://mms.mongodb.com/cluster",
					},
					{
						Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6",
						Rel:  "http://mms.mongodb.com/group",
					},
				},
				SnapshotID:        "5b6211ff87d9d663c59d3feg",
				TargetClusterName: "MyOtherCluster",
				TargetGroupID:     "5b6212af90dc76637950a2c6",
				Timestamp:         "2018-08-01T20:02:07Z",
			},
			{
				Cancelled:    false,
				CreatedAt:    "2018-08-01T22:05:41Z",
				DeliveryType: "download",
				DeliveryURL:  []string{"https://restore.mongodb.net:27017/restore-5b622e3587d9d6039faf713f.tar.gz"},
				Expired:      false,
				ExpiresAt:    "2018-08-02T02:03:33Z",
				ID:           "5b622e3587d9d6039faf713f",
				Links: []*Link{
					{
						Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0/backup/restoreJobs/5b622e3587d9d6039faf713f",
						Rel:  "self",
					},
					{
						Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0/backup/snapshots/5b6211ff87d9d663c59d3dee",
						Rel:  "http://mms.mongodb.com/snapshot",
					},
					{
						Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0",
						Rel:  "http://mms.mongodb.com/cluster",
					},
					{
						Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6",
						Rel:  "http://mms.mongodb.com/group",
					},
				},
				SnapshotID: "5b6211ff87d9d663c59d3feg",
				Timestamp:  "2018-08-01T20:02:07Z",
			},
		},
		TotalCount: 2,
	}

	if diff := deep.Equal(cloudProviderSnapshots, expected); diff != nil {
		t.Error(diff)
	}
	if !reflect.DeepEqual(cloudProviderSnapshots, expected) {
		t.Errorf("CloudProviderSnapshotRestoreJobs.List\n got=%#v\nwant=%#v", cloudProviderSnapshots, expected)
	}
}

func TestCloudProviderSnapshotRestoreJobs_Get(t *testing.T) {
	setup()
	defer teardown()

	requestParameters := &SnapshotReqPathParameters{
		GroupID:     "5b6212af90dc76637950a2c6",
		ClusterName: "MyCluster",
		JobID:       "5b622f7087d9d6039fafe03f",
	}

	path := fmt.Sprintf("/groups/%s/clusters/%s/backup/restoreJobs/%s", requestParameters.GroupID, requestParameters.ClusterName, requestParameters.JobID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"cancelled": false,
			"deliveryType": "automated",
			"expired": false,
			"expiresAt": "2018-08-02T02:08:48Z",
			"id": "5b622f7087d9d6039fafe03f",
			"links": [
				{
					"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0/backup/restoreJobs/5b622f7087d9d6039fafe03f",
					"rel": "self"
				},
				{
					"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0/backup/snapshots/5b6211ff87d9d663c59d3dee",
					"rel": "http://mms.mongodb.com/snapshot"
				},
				{
					"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0",
					"rel": "http://mms.mongodb.com/cluster"
				},
				{
					"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6",
					"rel": "http://mms.mongodb.com/group"
				}
			],
			"snapshotId": "5b6211ff87d9d663c59d3feg",
			"targetClusterName": "MyOtherCluster",
			"targetGroupId": "5b6212af90dc76637950a2c6",
			"timestamp": "2018-08-01T20:02:07Z"
		}`)
	})

	cloudProviderSnapshot, _, err := client.CloudProviderSnapshotRestoreJobs.Get(ctx, requestParameters)
	if err != nil {
		t.Errorf("CloudProviderSnapshotRestoreJobs.Get returned error: %v", err)
	}

	expected := &CloudProviderSnapshotRestoreJob{
		Cancelled:    false,
		DeliveryType: "automated",
		Expired:      false,
		ExpiresAt:    "2018-08-02T02:08:48Z",
		ID:           "5b622f7087d9d6039fafe03f",
		Links: []*Link{
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0/backup/restoreJobs/5b622f7087d9d6039fafe03f",
				Rel:  "self",
			},
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0/backup/snapshots/5b6211ff87d9d663c59d3dee",
				Rel:  "http://mms.mongodb.com/snapshot",
			},
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0",
				Rel:  "http://mms.mongodb.com/cluster",
			},
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6",
				Rel:  "http://mms.mongodb.com/group",
			},
		},
		SnapshotID:        "5b6211ff87d9d663c59d3feg",
		TargetClusterName: "MyOtherCluster",
		TargetGroupID:     "5b6212af90dc76637950a2c6",
		Timestamp:         "2018-08-01T20:02:07Z",
	}

	if diff := deep.Equal(cloudProviderSnapshot, expected); diff != nil {
		t.Error(diff)
	}
	if !reflect.DeepEqual(cloudProviderSnapshot, expected) {
		t.Errorf("CloudProviderSnapshotRestoreJobs.Get\n got=%#v\nwant=%#v", cloudProviderSnapshot, expected)
	}
}

func TestCloudProviderSnapshotRestoreJobs_Create(t *testing.T) {
	setup()
	defer teardown()

	requestParameters := &SnapshotReqPathParameters{
		GroupID:     "5b6212af90dc76637950a2c6",
		ClusterName: "MyClusterName",
	}

	createRequest := &CloudProviderSnapshotRestoreJob{
		SnapshotID:        "5b6211ff87d9d663c59d3feg",
		DeliveryType:      "automated",
		TargetClusterName: "MyOtherCluster",
		TargetGroupID:     "5b6212af90dc76637950a2c6",
	}

	path := fmt.Sprintf("/groups/%s/clusters/%s/backup/restoreJobs", requestParameters.GroupID, requestParameters.ClusterName)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"snapshotId":        "5b6211ff87d9d663c59d3feg",
			"deliveryType":      "automated",
			"targetClusterName": "MyOtherCluster",
			"targetGroupId":     "5b6212af90dc76637950a2c6",
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Decode json: %v", err)
		}

		if diff := deep.Equal(v, expected); diff != nil {
			t.Error(diff)
		}
		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprint(w, `{
			"cancelled": false,
			"deliveryType": "automated",
			"expired": false,
			"expiresAt": "2018-08-02T02:08:48Z",
			"id": "5b622f7087d9d6039fafe03f",
			"links": [
				{
					"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0/backup/restoreJobs/5b622f7087d9d6039fafe03f",
					"rel": "self"
				},
				{
					"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0/backup/snapshots/5b6211ff87d9d663c59d3dee",
					"rel": "http://mms.mongodb.com/snapshot"
				},
				{
					"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0",
					"rel": "http://mms.mongodb.com/cluster"
				},
				{
					"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6",
					"rel": "http://mms.mongodb.com/group"
				}
			],
			"snapshotId": "5b6211ff87d9d663c59d3feg",
			"targetClusterName": "MyOtherCluster",
			"targetGroupId": "5b6212af90dc76637950a2c6",
			"timestamp": "2018-08-01T20:02:07Z"
		}`)
	})

	cloudProviderSnapshot, _, err := client.CloudProviderSnapshotRestoreJobs.Create(ctx, requestParameters, createRequest)
	if err != nil {
		t.Errorf("CloudProviderSnapshotRestoreJobs.Create returned error: %v", err)
	}

	expected := &CloudProviderSnapshotRestoreJob{
		Cancelled:    false,
		DeliveryType: "automated",
		Expired:      false,
		ExpiresAt:    "2018-08-02T02:08:48Z",
		ID:           "5b622f7087d9d6039fafe03f",
		Links: []*Link{
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0/backup/restoreJobs/5b622f7087d9d6039fafe03f",
				Rel:  "self",
			},
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0/backup/snapshots/5b6211ff87d9d663c59d3dee",
				Rel:  "http://mms.mongodb.com/snapshot",
			},
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6/clusters/Cluster0",
				Rel:  "http://mms.mongodb.com/cluster",
			},
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5b6212af90dc76637950a2c6",
				Rel:  "http://mms.mongodb.com/group",
			},
		},
		SnapshotID:        "5b6211ff87d9d663c59d3feg",
		TargetClusterName: "MyOtherCluster",
		TargetGroupID:     "5b6212af90dc76637950a2c6",
		Timestamp:         "2018-08-01T20:02:07Z",
	}

	if diff := deep.Equal(cloudProviderSnapshot, expected); diff != nil {
		t.Error(diff)
	}
	if !reflect.DeepEqual(cloudProviderSnapshot, expected) {
		t.Errorf("CloudProviderSnapshotRestoreJobs.Create\n got=%#v\nwant=%#v", cloudProviderSnapshot, expected)
	}
}

func TestCloudProviderSnapshotRestoreJobs_Delete(t *testing.T) {
	setup()
	defer teardown()

	requestParameters := &SnapshotReqPathParameters{
		GroupID:     "5b6212af90dc76637950a2c6",
		ClusterName: "MyCluster",
		JobID:       "5b622f7087d9d6039fafe03f",
	}

	path := fmt.Sprintf("/groups/%s/clusters/%s/backup/restoreJobs/%s", requestParameters.GroupID, requestParameters.ClusterName, requestParameters.JobID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.CloudProviderSnapshotRestoreJobs.Delete(ctx, requestParameters)
	if err != nil {
		t.Errorf("CloudProviderSnapshotRestoreJobs.Delete returned error: %v", err)
	}
}
