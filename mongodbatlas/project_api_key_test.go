package mongodbatlas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

func TestProjectAPIKeys_ListAPIKeys(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/groups/1/apiKeys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"results": [
				{
					"desc": "test-apikey",
					"id": "5c47503320eef5699e1cce8d",
					"privateKey": "********-****-****-db2c132ca78d",
					"publicKey": "ewmaqvdo",
					"roles": [
						{
							"groupId": "1",
							"roleName": "GROUP_OWNER"
						},
						{
							"orgId": "1",
							"roleName": "ORG_MEMBER"
						}
					]
				},
				{
					"desc": "test-apikey-2",
					"id": "5c47503320eef5699e1cce8f",
					"privateKey": "********-****-****-db2c132ca78f",
					"publicKey": "ewmaqvde",
					"roles": [
						{
							"groupId": "1",
							"roleName": "GROUP_OWNER"
						},
						{
							"orgId": "1",
							"roleName": "ORG_MEMBER"
						}
					]
				}
			],
			"totalCount": 2
		}`)
	})

	apiKeys, _, err := client.ProjectAPIKeys.List(ctx, "1", nil)

	if err != nil {
		t.Errorf("APIKeys.List returned error: %v", err)
	}

	expected := []APIKey{
		{
			ID:         "5c47503320eef5699e1cce8d",
			Desc:       "test-apikey",
			PrivateKey: "********-****-****-db2c132ca78d",
			PublicKey:  "ewmaqvdo",
			Roles: []AtlasRole{
				{
					GroupID:  "1",
					RoleName: "GROUP_OWNER",
				},
				{
					OrgID:    "1",
					RoleName: "ORG_MEMBER",
				},
			},
		},
		{
			ID:         "5c47503320eef5699e1cce8f",
			Desc:       "test-apikey-2",
			PrivateKey: "********-****-****-db2c132ca78f",
			PublicKey:  "ewmaqvde",
			Roles: []AtlasRole{
				{
					GroupID:  "1",
					RoleName: "GROUP_OWNER",
				},
				{
					OrgID:    "1",
					RoleName: "ORG_MEMBER",
				},
			},
		},
	}
	if !reflect.DeepEqual(apiKeys, expected) {
		t.Errorf("APIKeys.List\n got=%#v\nwant=%#v", apiKeys, expected)
	}
}

func TestProjectAPIKeys_Assign(t *testing.T) {
	setup()
	defer teardown()

	groupID := "5953c5f380eef53887615f9a"
	keyID := "5d1d12c087d9d63e6d682438"

	mux.HandleFunc(fmt.Sprintf("/groups/%s/apiKeys/%s", groupID, keyID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
	})

	_, err := client.ProjectAPIKeys.Assign(ctx, groupID, keyID, &AssignAPIKey{})
	if err != nil {
		t.Errorf("ProjectAPIKeys.Assign returned error: %v", err)
	}
}

func TestProjectAPIKeys_Unassign(t *testing.T) {
	setup()
	defer teardown()

	groupID := "5953c5f380eef53887615f9a"
	keyID := "5d1d12c087d9d63e6d682438"

	mux.HandleFunc(fmt.Sprintf("/groups/%s/apiKeys/%s", groupID, keyID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.ProjectAPIKeys.Unassign(ctx, groupID, keyID)
	if err != nil {
		t.Errorf("ProjectAPIKeys.Unassign returned error: %v", err)
	}
}

func TestProjectAPIKeys_Create(t *testing.T) {
	setup()
	defer teardown()

	orgID := "1"

	createRequest := &APIKeyInput{
		Desc:  "test-apiKey",
		Roles: []string{"GROUP_OWNER"},
	}

	mux.HandleFunc(fmt.Sprintf("/groups/%s/apiKeys", orgID), func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"desc":  "test-apiKey",
			"roles": []interface{}{"GROUP_OWNER"},
		}

		jsonBlob := `
		{
			"desc": "test-apikey",
			"id": "5c47503320eef5699e1cce8d",
			"privateKey": "********-****-****-db2c132ca78d",
			"publicKey": "ewmaqvdo",
			"roles": [
				{
					"groupId": "1",
					"roleName": "GROUP_OWNER"
				},
				{
					"orgId": "1",
					"roleName": "ORG_MEMBER"
				}
			]
		}
		`

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if diff := deep.Equal(v, expected); diff != nil {
			t.Errorf("Clusters.Create Request Body = %v", diff)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprint(w, jsonBlob)
	})

	apiKey, _, err := client.ProjectAPIKeys.Create(ctx, orgID, createRequest)
	if err != nil {
		t.Errorf("ProjectAPIKeys.Create returned error: %v", err)
	}

	if desc := apiKey.Desc; desc != "test-apikey" {
		t.Errorf("expected username '%s', received '%s'", "test-apikeye", desc)
	}

	if pk := apiKey.PublicKey; pk != "ewmaqvdo" {
		t.Errorf("expected publicKey '%s', received '%s'", orgID, pk)
	}

}
