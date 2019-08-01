package mongodbatlas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

func TestWhitelistAPIKeys_List(t *testing.T) {
	setup()
	defer teardown()

	orgID := "ORG-ID"
	apiKeyID := "API-KEY-ID"

	mux.HandleFunc(fmt.Sprintf("/"+whitelistAPIKeysPath, orgID, apiKeyID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"links": [
				{
					"href": "https://cloud.mongodb.com/api/atlas/v1.0/orgs/599c510c80eef518f3b63fe1/apiKeys/5c49e72980eef544a218f8f8/whitelist/?pretty=true&pageNum=1&itemsPerPage=100",
					"rel": "self"
				}
			],
			"results": [
				{
					"cidrBlock": "147.58.184.16/32",
					"count": 0,
					"created": "2019-01-24T16:34:57Z",
					"ipAddress": "147.58.184.16",
					"lastUsed": "2019-01-24T20:18:25Z",
					"lastUsedAddress": "147.58.184.16",
					"links": [
						{
							"href": "https://cloud.mongodb.com/api/atlas/v1.0/orgs/{ORG-ID}/apiKeys/{API-KEY-ID}/whitelist/147.58.184.16",
							"rel": "self"
						}
					]
				},
				{
					"cidrBlock": "84.255.48.125/32",
					"count": 0,
					"created": "2019-01-24T16:26:37Z",
					"ipAddress": "84.255.48.125",
					"lastUsed": "2019-01-24T20:18:25Z",
					"lastUsedAddress": "84.255.48.125",
					"links": [
						{
							"href": "https://cloud.mongodb.com/api/atlas/v1.0/orgs/{ORG-ID}/apiKeys/{API-KEY-ID}/whitelist/206.252.195.126",
							"rel": "self"
						}
					]
				}
			],
			"totalCount": 2
		}`)
	})

	whitelistAPIKeys, _, err := client.WhitelistAPIKeys.List(ctx, orgID, apiKeyID)
	if err != nil {
		t.Errorf("WhitelistAPIKeys.List returned error: %v", err)
	}

	expected := &WhitelistAPIKeys{
		Links: []*Link{
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/orgs/599c510c80eef518f3b63fe1/apiKeys/5c49e72980eef544a218f8f8/whitelist/?pretty=true&pageNum=1&itemsPerPage=100",
				Rel:  "self",
			},
		},
		Results: []*WhitelistAPIKey{
			{
				CidrBlock:       "147.58.184.16/32",
				Count:           0,
				Created:         "2019-01-24T16:34:57Z",
				IPAddress:       "147.58.184.16",
				LastUsed:        "2019-01-24T20:18:25Z",
				LastUsedAddress: "147.58.184.16",
				Links: []*Link{
					{
						Href: "https://cloud.mongodb.com/api/atlas/v1.0/orgs/{ORG-ID}/apiKeys/{API-KEY-ID}/whitelist/147.58.184.16",
						Rel:  "self",
					},
				},
			},
			{
				CidrBlock:       "84.255.48.125/32",
				Count:           0,
				Created:         "2019-01-24T16:26:37Z",
				IPAddress:       "84.255.48.125",
				LastUsed:        "2019-01-24T20:18:25Z",
				LastUsedAddress: "84.255.48.125",
				Links: []*Link{
					{
						Href: "https://cloud.mongodb.com/api/atlas/v1.0/orgs/{ORG-ID}/apiKeys/{API-KEY-ID}/whitelist/206.252.195.126",
						Rel:  "self",
					},
				},
			},
		},
		TotalCount: 2,
	}

	if diff := deep.Equal(whitelistAPIKeys, expected); diff != nil {
		t.Error(diff)
	}
	if !reflect.DeepEqual(whitelistAPIKeys, expected) {
		t.Errorf("WhitelistAPIKeys.List\n got=%#v\nwant=%#v", whitelistAPIKeys, expected)
	}
}

func TestWhitelistAPIKeys_Get(t *testing.T) {
	setup()
	defer teardown()

	orgID := "ORG-ID"
	apiKeyID := "API-KEY-ID"
	ipAddress := "IP-ADDRESS"

	mux.HandleFunc(fmt.Sprintf("/"+whitelistAPIKeysPath+"/%s", orgID, apiKeyID, ipAddress), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"cidrBlock": "147.58.184.16/32",
			"count": 0,
			"created": "2019-01-24T16:34:57Z",
			"ipAddress": "147.58.184.16",
			"links": [
				{
					"href": "https://cloud.mongodb.com/api/atlas/v1.0/orgs/{ORG-ID}/apiKeys/{API-KEY-ID}/whitelist/147.58.184.16",
					"rel": "self"
				}
			]
		}`)
	})

	whitelistAPIKey, _, err := client.WhitelistAPIKeys.Get(ctx, orgID, apiKeyID, ipAddress)
	if err != nil {
		t.Errorf("WhitelistAPIKeys.Get returned error: %v", err)
	}

	expected := &WhitelistAPIKey{
		CidrBlock: "147.58.184.16/32",
		Count:     0,
		Created:   "2019-01-24T16:34:57Z",
		IPAddress: "147.58.184.16",
		Links: []*Link{
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/orgs/{ORG-ID}/apiKeys/{API-KEY-ID}/whitelist/147.58.184.16",
				Rel:  "self",
			},
		},
	}

	if diff := deep.Equal(whitelistAPIKey, expected); diff != nil {
		t.Error(diff)
	}
	if !reflect.DeepEqual(whitelistAPIKey, expected) {
		t.Errorf("WhitelistAPIKeys.Get\n got=%#v\nwant=%#v", whitelistAPIKey, expected)
	}
}

func TestWhitelistAPIKeys_Create(t *testing.T) {
	setup()
	defer teardown()

	orgID := "ORG-ID"
	apiKeyID := "API-KEY-ID"

	createRequest := &[]WhitelistAPIKeysReq{
		{
			IPAddress: "77.54.32.11",
			CidrBlock: "77.54.32.11/32",
		},
	}

	mux.HandleFunc(fmt.Sprintf("/"+whitelistAPIKeysPath, orgID, apiKeyID), func(w http.ResponseWriter, r *http.Request) {
		expected := []map[string]interface{}{
			{
				"ipAddress": "77.54.32.11",
				"cidrBlock": "77.54.32.11/32",
			},
		}

		var v []map[string]interface{}
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

		fmt.Fprintf(w, `{
			"links": [
				{
					"href": "https://cloud.mongodb.com/api/atlas/v1.0/orgs/599c510c80eef518f3b63fe1/apiKeys/5c49e72980eef544a218f8f8/whitelist/?pretty=true&pageNum=1&itemsPerPage=100",
					"rel": "self"
				}
			],
			"results": [
				{
					"cidrBlock": "147.58.184.16/32",
					"count": 0,
					"created": "2019-01-24T16:34:57Z",
					"ipAddress": "147.58.184.16",
					"lastUsed": "2019-01-24T20:18:25Z",
					"lastUsedAddress": "147.58.184.16",
					"links": [
						{
							"href": "https://cloud.mongodb.com/api/atlas/v1.0/orgs/{ORG-ID}/apiKeys/{API-KEY-ID}/whitelist/147.58.184.16",
							"rel": "self"
						}
					]
				},
				{
					"cidrBlock": "77.54.32.11/32",
					"count": 0,
					"created": "2019-01-24T16:26:37Z",
					"ipAddress": "77.54.32.11",
					"lastUsed": "2019-01-24T20:18:25Z",
					"lastUsedAddress": "77.54.32.11",
					"links": [
						{
							"href": "https://cloud.mongodb.com/api/atlas/v1.0/orgs/{ORG-ID}/apiKeys/{API-KEY-ID}/whitelist/77.54.32.11",
							"rel": "self"
						}
					]
				}
			],
			"totalCount": 2
		}`)
	})

	whitelistAPIKey, _, err := client.WhitelistAPIKeys.Create(ctx, orgID, apiKeyID, createRequest)
	if err != nil {
		t.Errorf("WhitelistAPIKeys.Create returned error: %v", err)
	}

	expected := &WhitelistAPIKeys{
		Links: []*Link{
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/orgs/599c510c80eef518f3b63fe1/apiKeys/5c49e72980eef544a218f8f8/whitelist/?pretty=true&pageNum=1&itemsPerPage=100",
				Rel:  "self",
			},
		},
		Results: []*WhitelistAPIKey{
			{
				CidrBlock:       "147.58.184.16/32",
				Count:           0,
				Created:         "2019-01-24T16:34:57Z",
				IPAddress:       "147.58.184.16",
				LastUsed:        "2019-01-24T20:18:25Z",
				LastUsedAddress: "147.58.184.16",
				Links: []*Link{
					{
						Href: "https://cloud.mongodb.com/api/atlas/v1.0/orgs/{ORG-ID}/apiKeys/{API-KEY-ID}/whitelist/147.58.184.16",
						Rel:  "self",
					},
				},
			},
			{
				CidrBlock:       "77.54.32.11/32",
				Count:           0,
				Created:         "2019-01-24T16:26:37Z",
				IPAddress:       "77.54.32.11",
				LastUsed:        "2019-01-24T20:18:25Z",
				LastUsedAddress: "77.54.32.11",
				Links: []*Link{
					{
						Href: "https://cloud.mongodb.com/api/atlas/v1.0/orgs/{ORG-ID}/apiKeys/{API-KEY-ID}/whitelist/77.54.32.11",
						Rel:  "self",
					},
				},
			},
		},
		TotalCount: 2,
	}

	if diff := deep.Equal(whitelistAPIKey, expected); diff != nil {
		t.Error(diff)
	}
	if !reflect.DeepEqual(whitelistAPIKey, expected) {
		t.Errorf("WhitelistAPIKeys.Create\n got=%#v\nwant=%#v", whitelistAPIKey, expected)
	}
}

func TestWhitelistAPIKeys_Delete(t *testing.T) {
	setup()
	defer teardown()

	orgID := "ORG-ID"
	apiKeyID := "API-KEY-ID"
	ipAddress := "IP-ADDRESS"

	mux.HandleFunc(fmt.Sprintf("/"+whitelistAPIKeysPath+"/%s", orgID, apiKeyID, ipAddress), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.WhitelistAPIKeys.Delete(ctx, orgID, apiKeyID, ipAddress)
	if err != nil {
		t.Errorf("WhitelistAPIKeys.Delete returned error: %v", err)
	}
}
