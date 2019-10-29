package mongodbatlas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/mwielbut/pointy"
)

func TestContainers_ListContainers(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/groups/1/containers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"results": [
				{
					"atlasCidrBlock": "10.8.0.0/18",
					"id": "1112269b3bf99403840e8934",
					"gcpProjectId": "my-sample-project-191923",
					"networkName": "test1",
					"providerName": "GCP",
					"provisioned": true
				},
				{
					"atlasCidrBlock" : "10.8.0.0/21",
					"id" : "1112269b3bf99403840e8934",
					"providerName" : "AWS",
					"provisioned" : true,
					"regionName" : "US_EAST_1",
					"vpcId" : "vpc-zz0zzzzz"
				}
			],
			"totalCount": 2
		}`)
	})

	containers, _, err := client.Containers.List(ctx, "1", nil)
	if err != nil {
		t.Errorf("Containers.List returned error: %v", err)
	}

	GCPContainer := Container{
		AtlasCIDRBlock: "10.8.0.0/18",
		ID:             "1112269b3bf99403840e8934",
		GCPProjectID:   "my-sample-project-191923",
		NetworkName:    "test1",
		ProviderName:   "GCP",
		Provisioned:    pointy.Bool(true),
	}

	AWSContainer := Container{
		AtlasCIDRBlock: "10.8.0.0/21",
		ID:             "1112269b3bf99403840e8934",
		ProviderName:   "AWS",
		Provisioned:    pointy.Bool(true),
		RegionName:     "US_EAST_1",
		VPCID:          "vpc-zz0zzzzz",
	}

	expected := []Container{GCPContainer, AWSContainer}
	if !reflect.DeepEqual(containers, expected) {
		t.Errorf("Containers.List\n got=%#v\nwant=%#v", containers, expected)
	}
}

func TestContainers_ListContainersMultiplePages(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/groups/1/containers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		dr := containersResponse{
			Results: []Container{
				{
					AtlasCIDRBlock: "10.8.0.0/18",
					ID:             "1112269b3bf99403840e8934",
					GCPProjectID:   "my-sample-project-191923",
					NetworkName:    "test1",
					ProviderName:   "GCP",
					Provisioned:    pointy.Bool(true),
				},
				{
					AtlasCIDRBlock: "10.8.0.0/21",
					ID:             "1112269b3bf99403840e8934",
					ProviderName:   "AWS",
					Provisioned:    pointy.Bool(true),
					RegionName:     "US_EAST_1",
					VPCID:          "vpc-zz0zzzzz",
				},
			},
			Links: []*Link{
				{Href: "http://example.com/api/atlas/v1.0/groups/1/containers?pageNum=2&itemsPerPage=2", Rel: "self"},
				{Href: "http://example.com/api/atlas/v1.0/groups/1/containers?pageNum=2&itemsPerPage=2", Rel: "previous"},
			},
		}

		b, err := json.Marshal(dr)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprint(w, string(b))
	})

	_, resp, err := client.Containers.List(ctx, "1", nil)
	if err != nil {
		t.Fatal(err)
	}

	checkCurrentPage(t, resp, 2)
}

func TestContainers_RetrievePageByNumber(t *testing.T) {
	setup()
	defer teardown()

	jBlob := `
	{
		"links": [
			{
				"href": "http://example.com/api/atlas/v1.0/groups/1/containers?pageNum=1&itemsPerPage=1",
				"rel": "previous"
			},
			{
				"href": "http://example.com/api/atlas/v1.0/groups/1/containers?pageNum=2&itemsPerPage=1",
				"rel": "self"
			},
			{
				"href": "http://example.com/api/atlas/v1.0/groups/1/containers?itemsPerPage=3&pageNum=2",
				"rel": "next"
			}
		],
		"results": [
			{
				"atlasCidrBlock": "10.8.0.0/18",
				"id": "1112269b3bf99403840e8934",
				"gcpProjectId": "my-sample-project-191923",
				"networkName": "test1",
				"providerName": "GCP",
				"provisioned": true
			}
		],
		"totalCount": 3
	}`

	mux.HandleFunc("/groups/1/containers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, jBlob)
	})

	opt := &ListOptions{PageNum: 2}
	_, resp, err := client.Containers.List(ctx, "1", opt)

	if err != nil {
		t.Fatal(err)
	}

	checkCurrentPage(t, resp, 2)
}

func TestContainers_Create(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1"

	createRequest := &Container{
		AtlasCIDRBlock: "10.8.0.0/18",
		ProviderName:   "GCP",
	}

	mux.HandleFunc(fmt.Sprintf("/groups/%s/containers", groupID), func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"atlasCidrBlock": "10.8.0.0/18",
			"providerName":   "GCP",
		}

		jsonBlob := `
		{
			"atlasCidrBlock": "10.8.0.0/18",
			"id": "1112269b3bf99403840e8934",
			"gcpProjectId": "my-sample-project-191923",
			"networkName": "test1",
			"providerName": "GCP",
			"provisioned": true
		}
		`

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}
		if diff := deep.Equal(v, expected); diff != nil {
			t.Errorf("Containers.Create Request Body = %v", diff)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprint(w, jsonBlob)
	})

	container, _, err := client.Containers.Create(ctx, groupID, createRequest)
	if err != nil {
		t.Errorf("Containers.Create returned error: %v", err)
	}

	expectedCID := "1112269b3bf99403840e8934"
	expectedGCPID := "my-sample-project-191923"

	if cID := container.ID; cID != expectedCID {
		t.Errorf("expected containerId '%s', received '%s'", expectedCID, cID)
	}

	if id := container.GCPProjectID; id != expectedGCPID {
		t.Errorf("expected gcpProjectId '%s', received '%s'", expectedGCPID, id)
	}

	if isProvisioned := container.Provisioned; !*isProvisioned {
		t.Errorf("expected provisioned '%t', received '%t'", true, *isProvisioned)
	}

}

func TestContainers_Update(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1"
	containerID := "1112269b3bf99403840e8934"

	updateRequest := &Container{
		AtlasCIDRBlock: "10.8.0.0/18",
		GCPProjectID:   "my-sample-project-191923",
		ProviderName:   "GCP",
	}

	mux.HandleFunc(fmt.Sprintf("/groups/%s/containers/%s", groupID, containerID), func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"atlasCidrBlock": "10.8.0.0/18",
			"gcpProjectId":   "my-sample-project-191923",
			"providerName":   "GCP",
		}

		jsonBlob := `
		{
			"atlasCidrBlock": "10.8.0.0/18",
			"id": "1112269b3bf99403840e8934",
			"gcpProjectId": "my-sample-project-191923",
			"networkName": "test1",
			"providerName": "GCP",
			"provisioned": true
		}
		`

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if diff := deep.Equal(v, expected); diff != nil {
			t.Errorf("Containers.Update Request Body = %v", diff)
		}
		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprint(w, jsonBlob)
	})

	container, _, err := client.Containers.Update(ctx, groupID, containerID, updateRequest)
	if err != nil {
		t.Errorf("Containers.Update returned error: %v", err)
	}

	if cID := container.ID; cID != containerID {
		t.Errorf("expected containerID '%s', received '%s'", containerID, cID)
	}

	expectedGCPID := "my-sample-project-191923"

	if id := container.GCPProjectID; id != expectedGCPID {
		t.Errorf("expected gcpProjectId '%s', received '%s'", expectedGCPID, id)
	}

	if isProvisioned := container.Provisioned; !*isProvisioned {
		t.Errorf("expected provisioned '%t', received '%t'", true, *isProvisioned)
	}

}

func TestContainers_Delete(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1"
	id := "1112222b3bf99403840e8934"

	mux.HandleFunc(fmt.Sprintf("/groups/%s/containers/%s", groupID, id), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Containers.Delete(ctx, groupID, id)
	if err != nil {
		t.Errorf("Container.Delete returned error: %v", err)
	}
}
