package mongodbatlas

import (
	"fmt"
	"net/http"
	"testing"
)

func TestProjectAPIKeys_Assign(t *testing.T) {
	setup()
	defer teardown()

	groupID := "5953c5f380eef53887615f9a"
	projectID := "5d1d12c087d9d63e6d682438"

	mux.HandleFunc(fmt.Sprintf("/groups/%s/apiKeys/%s", groupID, projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	_, err := client.ProjectAPIKeys.Assign(ctx, groupID, projectID)
	if err != nil {
		t.Errorf("ProjectIPWhitelist.Assign returned error: %v", err)
	}
}

func TestProjectAPIKeys_Unassign(t *testing.T) {
	setup()
	defer teardown()

	groupID := "5953c5f380eef53887615f9a"
	projectID := "5d1d12c087d9d63e6d682438"

	mux.HandleFunc(fmt.Sprintf("/groups/%s/apiKeys/%s", groupID, projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.ProjectAPIKeys.Unassign(ctx, groupID, projectID)
	if err != nil {
		t.Errorf("ProjectIPWhitelist.Assign returned error: %v", err)
	}
}
