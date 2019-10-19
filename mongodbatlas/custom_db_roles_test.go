package mongodbatlas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestCustomDBRoles_ListCustomDBRoles(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/groups/1/customDBRoles/roles", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[{"actions":[{"action":"CREATE_INDEX","resources":[{"collection":"test-collection","db":"test-db"}]}],"inheritedRoles":[{"db":"test-db","role":"read"}],"roleName":"test-role-name"}]`)
	})

	customDBRoles, _, err := client.CustomDBRoles.List(ctx, "1", nil)
	if err != nil {
		t.Errorf("CustomDBRoles.List returned error: %v", err)
	}

	expected := &[]CustomDbRole{{
		Actions: []Action{{
			Action: "CREATE_INDEX",
			Resources: []Resource{{
				Collection: "test-collection",
				Db:         "test-db",
				Cluster:    false,
			}},
		}},
		InheritedRoles: []InheritedRole{{
			Db:   "test-db",
			Role: "read",
		}},
		RoleName: "test-role-name",
	}}
	if !reflect.DeepEqual(customDBRoles, expected) {
		t.Errorf("CustomDBRoles.List\n got=%#v\nwant=%#v", customDBRoles, expected)
	}
}

func TestCustomDBRoles_GetCustomDBRole(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/groups/1/customDBRoles/roles/test-role-name", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"actions":[{"action":"CREATE_INDEX","resources":[{"collection":"test-collection","db":"test-db"}]}],"inheritedRoles":[{"db":"test-db","role":"read"}],"roleName":"test-role-name"}`)
	})

	customDBRole, _, err := client.CustomDBRoles.Get(ctx, "1", "test-role-name")
	if err != nil {
		t.Errorf("CustomDBRoles.Get returned error: %v", err)
	}

	expected := &CustomDbRole{
		Actions: []Action{{
			Action: "CREATE_INDEX",
			Resources: []Resource{{
				Collection: "test-collection",
				Db:         "test-db",
				Cluster:    false,
			}},
		}},
		InheritedRoles: []InheritedRole{{
			Db:   "test-db",
			Role: "read",
		}},
		RoleName: "test-role-name",
	}
	if !reflect.DeepEqual(customDBRole, expected) {
		t.Errorf("CustomDBRoles.List\n got=%#v\nwant=%#v", customDBRole, expected)
	}
}

func TestCustomDBRoles_CreateCustomDBRole(t *testing.T) {
	setup()
	defer teardown()

	createRequest := &CustomDbRole{
		Actions: []Action{{
			Action: "CREATE_INDEX",
			Resources: []Resource{{
				Collection: "test-collection",
				Db:         "test-db",
				Cluster:    false,
			}},
		}},
		InheritedRoles: []InheritedRole{{
			Db:   "test-db",
			Role: "read",
		}},
		RoleName: "test-role-name",
	}

	mux.HandleFunc("/groups/1/customDBRoles/roles", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"actions": []interface{}{map[string]interface{}{
				"action": "CREATE_INDEX",
				"resources": []interface{}{map[string]interface{}{
					"collection": "test-collection",
					"db":         "test-db",
				}},
			}},
			"inheritedRoles": []interface{}{map[string]interface{}{
				"db":   "test-db",
				"role": "read",
			}},
			"roleName": "test-role-name",
		}

		jsonBlob := `
		{
			"actions": [
				{
					"action": "CREATE_INDEX",
					"resources": [
						{
							"collection": "test-collection",
							"db": "test-db"
						}
					]
				}
			],
			"inheritedRoles": [
				{
					"db": "test-db",
					"role": "read"
				}
			],
			"roleName":"test-role-name"
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

		fmt.Fprintf(w, jsonBlob)
	})

	customDBRole, _, err := client.CustomDBRoles.Create(ctx, "1", createRequest)
	if err != nil {
		t.Errorf("CustomDBRoles.Create returned error: %v", err)
	}

	if roleName := customDBRole.RoleName; roleName != "test-role-name" {
		t.Errorf("expected roleName '%s', received '%s'", "test-role-name", roleName)
	}
}

func TestCustomDBRoles_UpdateCustomDBRole(t *testing.T) {
	setup()
	defer teardown()

	updateRequest := &CustomDbRole{
		Actions: []Action{{
			Action: "CREATE_INDEX",
			Resources: []Resource{{
				Collection: "test-collection",
				Db:         "test-db",
				Cluster:    false,
			}},
		}},
		InheritedRoles: []InheritedRole{{
			Db:   "test-db",
			Role: "read",
		}},
		RoleName: "test-role-name",
	}

	mux.HandleFunc("/groups/1/customDBRoles/roles/test-role-name", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"actions": []interface{}{map[string]interface{}{
				"action": "CREATE_INDEX",
				"resources": []interface{}{map[string]interface{}{
					"collection": "test-collection",
					"db":         "test-db",
				}},
			}},
			"inheritedRoles": []interface{}{map[string]interface{}{
				"db":   "test-db",
				"role": "read",
			}},
			"roleName": "test-role-name",
		}

		jsonBlob := `
		{
			"actions": [
				{
					"action": "CREATE_INDEX",
					"resources": [
						{
							"collection": "test-collection",
							"db": "test-db"
						}
					]
				}
			],
			"inheritedRoles": [
				{
					"db": "test-db",
					"role": "read"
				}
			],
			"roleName":"test-role-name"
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

		fmt.Fprintf(w, jsonBlob)
	})

	customDBRole, _, err := client.CustomDBRoles.Update(ctx, "1", "test-role-name", updateRequest)
	if err != nil {
		t.Errorf("CustomDBRoles.Update returned error: %v", err)
	}

	if roleName := customDBRole.RoleName; roleName != "test-role-name" {
		t.Errorf("expected roleName '%s', received '%s'", "test-role-name", roleName)
	}
}

func TestDatabaseUsers_DeleteCustomDBRole(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1"
	roleName := "test-role-name"

	mux.HandleFunc(fmt.Sprintf("/groups/%s/customDBRoles/roles/%s", groupID, roleName), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.CustomDBRoles.Delete(ctx, groupID, roleName)
	if err != nil {
		t.Errorf("CustomDBRole.Delete returned error: %v", err)
	}
}
