package mongodbatlas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

func TestPeers_ListPeers(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/groups/1/peers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"results": [
				{
					"containerId": "507f1f77bcf86cd799439011",
					"gcpProjectId": "my-sample-project-191923",
					"id": "1112222b3bf99403840e8934",
					"networkName": "test1",
					"status": "ADDING_PEER"
				},
				{
					"containerId": "507f1f77bcf86cd799439011",
					"gcpProjectId": "my-sample-project-191923",
					"id": "1112222b3bf99403840e8934",
					"networkName": "test1",
					"status": "ADDING_PEER"
				}
			],
			"totalCount": 2
		}`)
	})

	peers, _, err := client.Peers.List(ctx, "1", nil)
	if err != nil {
		t.Errorf("Peers.List returned error: %v", err)
	}

	peer := Peer{
		ContainerID:  "507f1f77bcf86cd799439011",
		GCPProjectID: "my-sample-project-191923",
		ID:           "1112222b3bf99403840e8934",
		NetworkName:  "test1",
		Status:       "ADDING_PEER",
	}

	expected := []Peer{peer, peer}
	if !reflect.DeepEqual(peers, expected) {
		t.Errorf("Peers.List\n got=%#v\nwant=%#v", peers, expected)
	}
}

func TestPeers_ListPeersMultiplePages(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/groups/1/peers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		dr := peersResponse{
			Results: []Peer{
				{
					ContainerID:  "507f1f77bcf86cd799439011",
					GCPProjectID: "my-sample-project-191923",
					ID:           "1112222b3bf99403840e8934",
					NetworkName:  "test1",
					Status:       "ADDING_PEER",
				},
				{
					ContainerID:  "507f1f77bcf86cd799439011",
					GCPProjectID: "my-sample-project-191923",
					ID:           "1112222b3bf99403840e8935",
					NetworkName:  "test1",
					Status:       "ADDING_PEER",
				},
			},
			Links: []*Link{
				{Href: "http://example.com/api/atlas/v1.0/groups/1/peers?pageNum=2&itemsPerPage=2", Rel: "self"},
				{Href: "http://example.com/api/atlas/v1.0/groups/1/peers?pageNum=2&itemsPerPage=2", Rel: "previous"},
			},
		}

		b, err := json.Marshal(dr)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprint(w, string(b))
	})

	_, resp, err := client.Peers.List(ctx, "1", nil)
	if err != nil {
		t.Fatal(err)
	}

	checkCurrentPage(t, resp, 2)
}

func TestPeers_RetrievePageByNumber(t *testing.T) {
	setup()
	defer teardown()

	jBlob := `
	{
		"links": [
			{
				"href": "http://example.com/api/atlas/v1.0/groups/1/peers?pageNum=1&itemsPerPage=1",
				"rel": "previous"
			},
			{
				"href": "http://example.com/api/atlas/v1.0/groups/1/peers?pageNum=2&itemsPerPage=1",
				"rel": "self"
			},
			{
				"href": "http://example.com/api/atlas/v1.0/groups/1/peers?itemsPerPage=3&pageNum=2",
				"rel": "next"
			}
		],
		"results": [
			{
				"containerId": "507f1f77bcf86cd799439011",
				"gcpProjectId": "my-sample-project-191923",
				"id": "1112222b3bf99403840e8934",
				"networkName": "test1",
				"status": "ADDING_PEER"
			}
		],
		"totalCount": 3
	}`

	mux.HandleFunc("/groups/1/peers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, jBlob)
	})

	opt := &ListOptions{PageNum: 2}
	_, resp, err := client.Peers.List(ctx, "1", opt)

	if err != nil {
		t.Fatal(err)
	}

	checkCurrentPage(t, resp, 2)
}

func TestPeers_Create(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1"

	createRequest := &Peer{
		ContainerID:  "507f1f77bcf86cd799439011",
		GCPProjectID: "my-sample-project-191923",
		ID:           "1112222b3bf99403840e8934",
		NetworkName:  "test1",
		Status:       "ADDING_PEER",
	}

	mux.HandleFunc(fmt.Sprintf("/groups/%s/peers", groupID), func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"containerId":  "507f1f77bcf86cd799439011",
			"gcpProjectId": "my-sample-project-191923",
			"id":           "1112222b3bf99403840e8934",
			"networkName":  "test1",
			"status":       "ADDING_PEER",
		}

		jsonBlob := `
		{
			"containerId":  "507f1f77bcf86cd799439011",
			"gcpProjectId": "my-sample-project-191923",
			"id":           "1112222b3bf99403840e8934",
			"networkName":  "test1",
			"status":       "ADDING_PEER"
		}
		`

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}
		if diff := deep.Equal(v, expected); diff != nil {
			t.Errorf("Peers.Create Request Body = %v", diff)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprint(w, jsonBlob)
	})

	peer, _, err := client.Peers.Create(ctx, groupID, createRequest)
	if err != nil {
		t.Errorf("Peers.Create returned error: %v", err)
	}

	expectedCID := "507f1f77bcf86cd799439011"
	expectedGCPID := "my-sample-project-191923"

	if cID := peer.ContainerID; cID != expectedCID {
		t.Errorf("expected containerId '%s', received '%s'", expectedCID, cID)
	}

	if id := peer.GCPProjectID; id != expectedGCPID {
		t.Errorf("expected gcpProjectId '%s', received '%s'", expectedGCPID, id)
	}

}

func TestPeers_Update(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1"
	peerID := "1112222b3bf99403840e8934"

	updateRequest := &Peer{
		ContainerID:  "507f1f77bcf86cd799439011",
		GCPProjectID: "my-sample-project-191923",
		ID:           "1112222b3bf99403840e8934",
		NetworkName:  "test1",
		Status:       "ADDING_PEER",
	}

	mux.HandleFunc(fmt.Sprintf("/groups/%s/peers/%s", groupID, peerID), func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"containerId":  "507f1f77bcf86cd799439011",
			"gcpProjectId": "my-sample-project-191923",
			"id":           "1112222b3bf99403840e8934",
			"networkName":  "test1",
			"status":       "ADDING_PEER",
		}

		jsonBlob := `
		{
			"containerId":  "507f1f77bcf86cd799439011",
			"gcpProjectId": "my-sample-project-191923",
			"id":           "1112222b3bf99403840e8934",
			"networkName":  "test1",
			"status":       "ADDING_PEER"
		}
		`

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if diff := deep.Equal(v, expected); diff != nil {
			t.Errorf("Peers.Update Request Body = %v", diff)
		}
		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprint(w, jsonBlob)
	})

	perr, _, err := client.Peers.Update(ctx, groupID, peerID, updateRequest)
	if err != nil {
		t.Errorf("Peers.Update returned error: %v", err)
	}

	if pID := perr.ID; pID != peerID {
		t.Errorf("expected peerID '%s', received '%s'", peerID, pID)
	}

	expectedGCPID := "my-sample-project-191923"

	if id := perr.GCPProjectID; id != expectedGCPID {
		t.Errorf("expected groupId '%s', received '%s'", expectedGCPID, id)
	}

}

func TestPeers_Delete(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1"
	id := "1112222b3bf99403840e8934"

	mux.HandleFunc(fmt.Sprintf("/groups/%s/peers/%s", groupID, id), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Peers.Delete(ctx, groupID, id)
	if err != nil {
		t.Errorf("Peer.Delete returned error: %v", err)
	}
}
