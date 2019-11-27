package mongodbatlas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

func TestAPIKeys_ListAPIKeys(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/1/apiKeys", func(w http.ResponseWriter, r *http.Request) {
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

	apiKeys, _, err := client.APIKeys.List(ctx, "1", nil)

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

func TestAPIKeys_ListAPIKeysMultiplePages(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/1/apiKeys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		dr := apiKeysResponse{
			Results: []APIKey{
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
			},
			Links: []*Link{
				{Href: "http://example.com/api/atlas/v1.0/orgs/1/apiKeys?pageNum=2&itemsPerPage=2", Rel: "self"},
				{Href: "http://example.com/api/atlas/v1.0/orgs/1/apiKeys?pageNum=2&itemsPerPage=2", Rel: "previous"},
			},
		}

		b, err := json.Marshal(dr)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprint(w, string(b))
	})

	_, resp, err := client.APIKeys.List(ctx, "1", nil)
	if err != nil {
		t.Fatal(err)
	}

	checkCurrentPage(t, resp, 2)
}

func TestAPIKeys_RetrievePageByNumber(t *testing.T) {
	setup()
	defer teardown()

	jBlob := `
	{
		"links": [
			{
				"href": "http://example.com/api/atlas/v1.0/orgs/1/apikeys?pageNum=1&itemsPerPage=1",
				"rel": "previous"
			},
			{
				"href": "http://example.com/api/atlas/v1.0/orgs/1/apikeys?pageNum=2&itemsPerPage=1",
				"rel": "self"
			},
			{
				"href": "http://example.com/api/atlas/v1.0/orgs/1/apikeys?itemsPerPage=3&pageNum=2",
				"rel": "next"
			}
		],
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
			}
		],
		"totalCount": 3
	}`

	mux.HandleFunc("/orgs/1/apiKeys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, jBlob)
	})

	opt := &ListOptions{PageNum: 2}
	_, resp, err := client.APIKeys.List(ctx, "1", opt)
	if err != nil {
		t.Fatal(err)
	}

	checkCurrentPage(t, resp, 2)
}

func TestAPIKeys_Create(t *testing.T) {
	setup()
	defer teardown()

	orgID := "1"

	createRequest := &APIKeyInput{
		Desc:  "test-apiKey",
		Roles: []string{"GROUP_OWNER"},
	}

	mux.HandleFunc(fmt.Sprintf("/orgs/%s/apiKeys", orgID), func(w http.ResponseWriter, r *http.Request) {
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

	apiKey, _, err := client.APIKeys.Create(ctx, orgID, createRequest)
	if err != nil {
		t.Errorf("APIKeys.Create returned error: %v", err)
	}

	if desc := apiKey.Desc; desc != "test-apikey" {
		t.Errorf("expected username '%s', received '%s'", "test-apikeye", desc)
	}

	if pk := apiKey.PublicKey; pk != "ewmaqvdo" {
		t.Errorf("expected publicKey '%s', received '%s'", orgID, pk)
	}

}

func TestAPIKeys_GetAPIKey(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/1/apiKeys/5c47503320eef5699e1cce8d", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"desc":"test-desc"}`)
	})

	apiKeys, _, err := client.APIKeys.Get(ctx, "1", "5c47503320eef5699e1cce8d")
	if err != nil {
		t.Errorf("APIKey.Get returned error: %v", err)
	}

	expected := &APIKey{Desc: "test-desc"}

	if diff := deep.Equal(apiKeys, expected); diff != nil {
		t.Errorf("Clusters.Get = %v", diff)
	}

	if !reflect.DeepEqual(apiKeys, expected) {
		t.Errorf("APIKeys.Get\n got=%#v\nwant=%#v", apiKeys, expected)
	}
}

func TestAPIKeys_Update(t *testing.T) {
	setup()
	defer teardown()

	orgID := "1"

	updateRequest := &APIKeyInput{
		Desc:  "test-apiKey",
		Roles: []string{"GROUP_OWNER"},
	}

	mux.HandleFunc(fmt.Sprintf("/orgs/%s/apiKeys/%s", orgID, "5c47503320eef5699e1cce8d"), func(w http.ResponseWriter, r *http.Request) {
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

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprint(w, jsonBlob)
	})

	apiKey, _, err := client.APIKeys.Update(ctx, orgID, "5c47503320eef5699e1cce8d", updateRequest)
	if err != nil {
		t.Errorf("APIKeys.Create returned error: %v", err)
	}

	if desc := apiKey.Desc; desc != "test-apikey" {
		t.Errorf("expected username '%s', received '%s'", "test-apikeye", desc)
	}

	if pk := apiKey.PublicKey; pk != "ewmaqvdo" {
		t.Errorf("expected publicKey '%s', received '%s'", orgID, pk)
	}

}

func TestAPIKeys_Delete(t *testing.T) {
	setup()
	defer teardown()

	orgID := "1"
	apiKeyID := "5c47503320eef5699e1cce8d"

	mux.HandleFunc(fmt.Sprintf("/orgs/%s/apiKeys/%s", orgID, apiKeyID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.APIKeys.Delete(ctx, orgID, apiKeyID)
	if err != nil {
		t.Errorf("APIKey.Delete returned error: %v", err)
	}
}
