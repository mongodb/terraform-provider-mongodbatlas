package mongodbatlas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestProjectIPWhitelist_ListProjectIPWhitelist(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/groups/1/whitelist", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"results": [{"groupId":"1", "cidrBlock":"0.0.0.0/12", "ipAddress":"0.0.0.0"},{"groupId":"1", "cidrBlock":"0.0.0.1/12", "ipAddress":"0.0.0.1"}], "totalCount":2}`)
	})

	projectIPWhitelists, _, err := client.ProjectIPWhitelist.List(ctx, "1", nil)
	if err != nil {
		t.Errorf("ProjectIPWhitelist.List returned error: %v", err)
	}

	expected := []ProjectIPWhitelist{{GroupID: "1", CIDRBlock: "0.0.0.0/12", IPAddress: "0.0.0.0"}, {GroupID: "1", CIDRBlock: "0.0.0.1/12", IPAddress: "0.0.0.1"}}
	if !reflect.DeepEqual(projectIPWhitelists, expected) {
		t.Errorf("ProjectIPWhitelist.List\n got=%#v\nwant=%#v", projectIPWhitelists, expected)
	}
}

func TestProjectIPWhitelist_ListProjectIPWhitelistMultiplePages(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/groups/1/whitelist", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		dr := projectIPWhitelistsResponse{
			Results: []ProjectIPWhitelist{
				{GroupID: "1", CIDRBlock: "0.0.0.1/12"},
				{GroupID: "1", CIDRBlock: "0.0.0.0/12"},
			},
			Links: []*Link{
				{Href: "http://example.com/api/atlas/v1.0/groups/1/whitelist?pageNum=2&itemsPerPage=2", Rel: "self"},
				{Href: "http://example.com/api/atlas/v1.0/groups/1/whitelist?pageNum=1&itemsPerPage=2", Rel: "previous"},
			},
		}

		b, err := json.Marshal(dr)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprint(w, string(b))
	})

	_, resp, err := client.ProjectIPWhitelist.List(ctx, "1", nil)
	if err != nil {
		t.Fatal(err)
	}

	checkCurrentPage(t, resp, 2)
}

func TestProjectIPWhitelist_RetrievePageByNumber(t *testing.T) {
	setup()
	defer teardown()

	jBlob := `
	{
		"links": [
			{
				"href": "http://example.com/api/atlas/v1.0/groups/1/whitelist?pageNum=1&itemsPerPage=1",
				"rel": "previous"
			},
			{
				"href": "http://example.com/api/atlas/v1.0/groups/1/whitelist?pageNum=2&itemsPerPage=1",
				"rel": "self"
			},
			{
				"href": "http://example.com/api/atlas/v1.0/groups/1/whitelist?itemsPerPage=3&pageNum=2",
				"rel": "next"
			}
		],
		"results": [
			{
				"groupId":"1",
				"cidrBlock":"0.0.0.0/12",
				"ipAddress":"0.0.0.0",
				"links": [
					{
						"href": "http://example.com/api/atlas/v1.0/groups/1/whitelist/0.0.0.0",
						"rel": "self"
					}
				]
			}
		],
		"totalCount": 3
	}`

	mux.HandleFunc("/groups/1/whitelist", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, jBlob)
	})

	opt := &ListOptions{PageNum: 2}
	_, resp, err := client.ProjectIPWhitelist.List(ctx, "1", opt)
	if err != nil {
		t.Fatal(err)
	}

	checkCurrentPage(t, resp, 2)
}

func TestProjectIPWhitelist_Create(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1"

	createRequest := []*ProjectIPWhitelist{
		{
			GroupID:   groupID,
			CIDRBlock: "0.0.0.1/12",
		},
	}

	mux.HandleFunc(fmt.Sprintf("/groups/%s/whitelist", groupID), func(w http.ResponseWriter, r *http.Request) {
		expected := []map[string]interface{}{
			{
				"cidrBlock": "0.0.0.1/12",
				"groupId":   groupID,
			},
		}

		jsonBlob := `
		{
			"results": [{
				"groupId":"1",
				"cidrBlock":"0.0.0.1/12"
			}]
		}
		`

		var v []map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, jsonBlob)
	})

	projectIPWhitelist, _, err := client.ProjectIPWhitelist.Create(ctx, groupID, createRequest)
	if err != nil {
		t.Errorf("ProjectIPWhitelist.Create returned error: %v", err)
		return
	}

	if len(projectIPWhitelist) > 1 {
		t.Error("expected projectIPWhitelist response a list of one element")
	}

	if cidrBlock := projectIPWhitelist[0].CIDRBlock; cidrBlock != "0.0.0.1/12" {
		t.Errorf("expected cidrBlock '%s', received '%s'", "0.0.0.1/12", cidrBlock)
	}

	if id := projectIPWhitelist[0].GroupID; id != groupID {
		t.Errorf("expected groupId '%s', received '%s'", groupID, id)
	}

}

func TestProjectIPWhitelist_GetProjectIPWhitelist(t *testing.T) {
	setup()
	defer teardown()

	ipAddress := "0.0.0.0"

	mux.HandleFunc(fmt.Sprintf("/groups/1/whitelist/%s", ipAddress), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, fmt.Sprintf(`{"ipAddress":"%s"}`, ipAddress))
	})

	projectIPWhitelists, _, err := client.ProjectIPWhitelist.Get(ctx, "1", ipAddress)
	if err != nil {
		t.Errorf("ProjectIPWhitelist.Get returned error: %v", err)
	}

	expected := &ProjectIPWhitelist{IPAddress: ipAddress}
	if !reflect.DeepEqual(projectIPWhitelists, expected) {
		t.Errorf("ProjectIPWhitelist.Get\n got=%#v\nwant=%#v", projectIPWhitelists, expected)
	}
}

func TestProjectIPWhitelist_GetProjectIPWhitelistByCIDR(t *testing.T) {
	setup()
	defer teardown()

	cidr := "0.0.0.0/32"

	mux.HandleFunc(fmt.Sprintf("/groups/1/whitelist/%s", cidr), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, fmt.Sprintf(`{"cidrBlock":"%s"}`, cidr))
	})

	projectIPWhitelists, _, err := client.ProjectIPWhitelist.Get(ctx, "1", cidr)
	if err != nil {
		t.Errorf("ProjectIPWhitelist.Get returned error: %v", err)
	}

	expected := &ProjectIPWhitelist{CIDRBlock: cidr}
	if !reflect.DeepEqual(projectIPWhitelists, expected) {
		t.Errorf("ProjectIPWhitelist.Get\n got=%#v\nwant=%#v", projectIPWhitelists, expected)
	}
}

func TestProjectIPWhitelist_Update(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1"
	ipAddress := "0.0.0.1"

	createRequest := []*ProjectIPWhitelist{{
		GroupID:   groupID,
		IPAddress: ipAddress,
	}}

	mux.HandleFunc(fmt.Sprintf("/groups/%s/whitelist/%s", groupID, ipAddress), func(w http.ResponseWriter, r *http.Request) {
		expected := []map[string]interface{}{{
			"ipAddress": ipAddress,
			"groupId":   groupID,
		}}

		jsonBlob := `
		{
			"results": [
				{
					"ipAddress": "0.0.0.1",
					"groupId": "1"
				}
			]
		}
		`

		var v []map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, jsonBlob)
	})

	projectIPWhitelist, _, err := client.ProjectIPWhitelist.Update(ctx, groupID, ipAddress, createRequest)
	if err != nil {
		t.Errorf("ProjectIPWhitelist.Update returned error: %v", err)
		return
	}

	if ip := projectIPWhitelist[0].IPAddress; ip != ipAddress {
		t.Errorf("expected ipAddress '%s', received '%s'", ipAddress, ipAddress)
	}

	if id := projectIPWhitelist[0].GroupID; id != groupID {
		t.Errorf("expected groupId '%s', received '%s'", groupID, id)
	}

}

func TestProjectIPWhitelist_Delete(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1"
	ipAddress := "0.0.0.1"

	mux.HandleFunc(fmt.Sprintf("/groups/%s/whitelist/%s", groupID, ipAddress), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.ProjectIPWhitelist.Delete(ctx, groupID, ipAddress)
	if err != nil {
		t.Errorf("ProjectIPWhitelist.Delete returned error: %v", err)
	}
}
