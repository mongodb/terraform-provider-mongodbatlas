package mongodbatlas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDatabaseUsers_ListDatabaseUsers(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/groups/1/databaseUsers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"results": [{"groupId":"1", "username":"test-username"},{"groupId":"1", "username":"test-username-1"}], "totalCount":2}`)
	})

	dbUsers, _, err := client.DatabaseUsers.List(ctx, "1", nil)
	if err != nil {
		t.Errorf("DatabaseUsers.List returned error: %v", err)
	}

	expected := []DatabaseUser{{GroupID: "1", Username: "test-username"}, {GroupID: "1", Username: "test-username-1"}}
	if !reflect.DeepEqual(dbUsers, expected) {
		t.Errorf("DatabaseUsers.List\n got=%#v\nwant=%#v", dbUsers, expected)
	}
}

func TestDatabaseUsers_ListDatabaseUsersMultiplePages(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/groups/1/databaseUsers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		dr := databaseUsers{
			Results: []DatabaseUser{
				{GroupID: "1", Username: "test-one"},
				{GroupID: "1", Username: "test-two"},
			},
			Links: []*Link{
				{Href: "http://example.com/api/atlas/v1.0/groups/1/databaseUsers?pageNum=2&itemsPerPage=2", Rel: "self"},
				{Href: "http://example.com/api/atlas/v1.0/groups/1/databaseUsers?pageNum=2&itemsPerPage=2", Rel: "previous"},
			},
		}

		b, err := json.Marshal(dr)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprint(w, string(b))
	})

	_, resp, err := client.DatabaseUsers.List(ctx, "1", nil)
	if err != nil {
		t.Fatal(err)
	}

	checkCurrentPage(t, resp, 2)
}

func TestDatabaseUsers_RetrievePageByNumber(t *testing.T) {
	setup()
	defer teardown()

	jBlob := `
	{
		"links": [
			{
				"href": "http://example.com/api/atlas/v1.0/groups/1/databaseUsers?pageNum=1&itemsPerPage=1",
				"rel": "previous"
			},
			{
				"href": "http://example.com/api/atlas/v1.0/groups/1/databaseUsers?pageNum=2&itemsPerPage=1",
				"rel": "self"
			},
			{
				"href": "http://example.com/api/atlas/v1.0/groups/1/databaseUsers?itemsPerPage=3&pageNum=2",
				"rel": "next"
			}
		],
		"results": [
			{
				"databaseName": "admin",
				"groupId": "1",
				"ldapAuthType": "NONE",
				"links": [
					{
						"href": "http://example.com/api/atlas/v1.0/groups/1/databaseUsers/admin/test-test",
						"rel": "self"
					}
				],
				"roles": [
					{
						"databaseName": "admin",
						"roleName": "atlasAdmin"
					}
				],
				"username": "test-test"
			}
		],
		"totalCount": 3
	}`

	mux.HandleFunc("/groups/1/databaseUsers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, jBlob)
	})

	opt := &ListOptions{PageNum: 2}
	_, resp, err := client.DatabaseUsers.List(ctx, "1", opt)
	if err != nil {
		t.Fatal(err)
	}

	checkCurrentPage(t, resp, 2)
}

func TestDatabaseUsers_Create(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1"

	createRequest := &DatabaseUser{
		GroupID:      groupID,
		Username:     "test-username",
		Password:     "test-password",
		DatabaseName: "test-databasename",
		Roles: []Role{{
			DatabaseName:   "test-databasename",
			CollectionName: "test-collection-name",
			RoleName:       "test-role-name",
		}},
	}

	mux.HandleFunc(fmt.Sprintf("/groups/%s/databaseUsers", groupID), func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"username":     "test-username",
			"password":     "test-password",
			"databaseName": "test-databasename",
			"groupId":      groupID,
			"roles": []interface{}{map[string]interface{}{
				"databaseName":   "test-databasename",
				"collectionName": "test-collection-name",
				"roleName":       "test-role-name",
			}},
		}

		jsonBlob := `
		{
			"username": "test-username",
			"password": "test-password",
			"databaseName": "test-databasename",
			"groupId": "1",
			"roles": [
				{
					"databaseName": "test-databasename",
					"collectionName": "test-collection-name",
					"roleName": "test-role-name"
				}
			]
		}
		`

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprint(w, jsonBlob)
	})

	dbUser, _, err := client.DatabaseUsers.Create(ctx, groupID, createRequest)
	if err != nil {
		t.Errorf("DatabaseUsers.Create returned error: %v", err)
	}

	if username := dbUser.Username; username != "test-username" {
		t.Errorf("expected username '%s', received '%s'", "test-username", username)
	}

	if id := dbUser.GroupID; id != groupID {
		t.Errorf("expected groupId '%s', received '%s'", groupID, id)
	}

}

func TestDatabaseUsers_GetDatabaseUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/groups/1/databaseUsers/admin/test-username", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"username":"test-username"}`)
	})

	dbUsers, _, err := client.DatabaseUsers.Get(ctx, "1", "test-username")
	if err != nil {
		t.Errorf("DatabaseUser.Get returned error: %v", err)
	}

	expected := &DatabaseUser{Username: "test-username"}
	if !reflect.DeepEqual(dbUsers, expected) {
		t.Errorf("DatabaseUsers.Get\n got=%#v\nwant=%#v", dbUsers, expected)
	}
}

func TestDatabaseUsers_Update(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1"

	createRequest := &DatabaseUser{
		GroupID:      groupID,
		Username:     "test-username",
		Password:     "test-password",
		DatabaseName: "test-databasename",
		Roles: []Role{{
			DatabaseName:   "test-databasename",
			CollectionName: "test-collection-name",
			RoleName:       "test-role-name",
		}},
	}

	mux.HandleFunc(fmt.Sprintf("/groups/%s/databaseUsers/admin/%s", groupID, "test-username"), func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"username":     "test-username",
			"password":     "test-password",
			"databaseName": "test-databasename",
			"groupId":      groupID,
			"roles": []interface{}{map[string]interface{}{
				"databaseName":   "test-databasename",
				"collectionName": "test-collection-name",
				"roleName":       "test-role-name",
			}},
		}

		jsonBlob := `
		{
			"username": "test-username",
			"password": "test-password",
			"databaseName": "test-databasename",
			"groupId": "1",
			"roles": [
				{
					"databaseName": "test-databasename",
					"collectionName": "test-collection-name",
					"roleName": "test-role-name"
				}
			]
		}
		`

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprint(w, jsonBlob)
	})

	dbUser, _, err := client.DatabaseUsers.Update(ctx, groupID, "test-username", createRequest)
	if err != nil {
		t.Errorf("DatabaseUsers.Update returned error: %v", err)
	}

	if username := dbUser.Username; username != "test-username" {
		t.Errorf("expected username '%s', received '%s'", "test-username", username)
	}

	if id := dbUser.GroupID; id != groupID {
		t.Errorf("expected groupId '%s', received '%s'", groupID, id)
	}

}

func TestDatabaseUsers_Delete(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1"
	username := "test-username"

	mux.HandleFunc(fmt.Sprintf("/groups/%s/databaseUsers/admin/%s", groupID, username), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.DatabaseUsers.Delete(ctx, groupID, username)
	if err != nil {
		t.Errorf("DatabaseUser.Delete returned error: %v", err)
	}
}
