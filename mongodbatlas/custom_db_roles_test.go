package mongodbatlas

import (
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
