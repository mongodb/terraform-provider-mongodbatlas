package mongodbatlas

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

func TestProject_GetAllProjects(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"links" : [ {
				"href" : "https://cloud.mongodb.com/api/atlas/v1.0/groups",
				"rel" : "self"
			} ],
			"results" : [ {
				"clusterCount" : 2,
				"created" : "2016-07-14T14:19:33Z",
				"id" : "5a0a1e7e0f2912c554080ae6",
				"links" : [ {
					"href" : "https://cloud.mongodb.com/api/atlas/v1.0/groups/5a0a1e7e0f2912c554080ae6",
					"rel" : "self"
				} ],
				"name" : "ProjectBar",
				"orgId" : "5a0a1e7e0f2912c554080adc"
			}, {
				"clusterCount" : 0,
				"created" : "2017-10-16T15:24:01Z",
				"id" : "5a0a1e7e0f2912c554080ae7",
				"links" : [ {
					"href" : "https://cloud.mongodb.com/api/atlas/v1.0/groups/5a0a1e7e0f2912c554080ae7",
					"rel" : "self"
				} ],
				"name" : "Project Foo",
				"orgId" : "5a0a1e7e0f2912c554080adc"
			} ],
			"totalCount" : 2
		}`)
	})

	projects, _, err := client.Projects.GetAllProjects(ctx)
	if err != nil {
		t.Errorf("Projects.GetAllProjects returned error: %v", err)
	}

	expected := &Projects{
		Links: []*Link{
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups",
				Rel:  "self",
			},
		},
		Results: []*Project{
			{
				ClusterCount: 2,
				Created:      "2016-07-14T14:19:33Z",
				ID:           "5a0a1e7e0f2912c554080ae6",
				Links: []*Link{
					{
						Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5a0a1e7e0f2912c554080ae6",
						Rel:  "self",
					},
				},
				Name:  "ProjectBar",
				OrgID: "5a0a1e7e0f2912c554080adc",
			},
			{
				ClusterCount: 0,
				Created:      "2017-10-16T15:24:01Z",
				ID:           "5a0a1e7e0f2912c554080ae7",
				Links: []*Link{
					{
						Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5a0a1e7e0f2912c554080ae7",
						Rel:  "self",
					},
				},
				Name:  "Project Foo",
				OrgID: "5a0a1e7e0f2912c554080adc",
			},
		},
		TotalCount: 2,
	}

	if !reflect.DeepEqual(projects, expected) {
		t.Errorf("Projects.GetAllProjects\n got=%#v\nwant=%#v", projects, expected)
	}
}

func TestProject_GetOneProject(t *testing.T) {
	setup()
	defer teardown()

	projectID := "5a0a1e7e0f2912c554080adc"

	mux.HandleFunc(fmt.Sprintf("/%s/%s", projectBasePath, projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"id" : "5a0a1e7e0f2912c554080ae6",
			"orgId" : "5a0a1e7e0f2912c554080adc",
			"name" : "DocsFeedbackGroup",
			"clusterCount" : 2,
			"created" : "2016-07-14T14:19:33Z",
			"links" : [ {
				"href" : "https://cloud.mongodb.com/api/atlas/v1.0/groups/5a0a1e7e0f2912c554080ae6",
				"rel" : "self"
			} ]
		}`)
	})

	projectResponse, _, err := client.Projects.GetOneProject(ctx, projectID)
	if err != nil {
		t.Errorf("Projects.GetOneProject returned error: %v", err)
	}

	expected := &Project{
		ID:           "5a0a1e7e0f2912c554080ae6",
		OrgID:        "5a0a1e7e0f2912c554080adc",
		Name:         "DocsFeedbackGroup",
		ClusterCount: 2,
		Created:      "2016-07-14T14:19:33Z",
		Links: []*Link{
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5a0a1e7e0f2912c554080ae6",
				Rel:  "self",
			},
		},
	}

	if !reflect.DeepEqual(projectResponse, expected) {
		t.Errorf("Projects.GetOneProject\n got=%#v\nwant=%#v", projectResponse, expected)
	}
}

func TestProject_GetOneProjectByName(t *testing.T) {
	setup()
	defer teardown()

	projectName := "5a0a1e7e0f2912c554080adc"

	mux.HandleFunc(fmt.Sprintf("/%s/byName/%s", projectBasePath, projectName), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"id" : "5a0a1e7e0f2912c554080ae6",
			"orgId" : "5a0a1e7e0f2912c554080adc",
			"name" : "ProjectBar",
			"clusterCount" : 2,
			"created" : "2016-07-14T14:19:33Z",
			"links" : [ {
				"href" : "https://cloud.mongodb.com/api/atlas/v1.0/groups/5a0a1e7e0f2912c554080ae6",
				"rel" : "self"
			} ]
		}`)
	})

	projectResponse, _, err := client.Projects.GetOneProjectByName(ctx, projectName)
	if err != nil {
		t.Errorf("Projects.GetOneProject returned error: %v", err)
	}

	expected := &Project{
		ID:           "5a0a1e7e0f2912c554080ae6",
		OrgID:        "5a0a1e7e0f2912c554080adc",
		Name:         "ProjectBar",
		ClusterCount: 2,
		Created:      "2016-07-14T14:19:33Z",
		Links: []*Link{
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5a0a1e7e0f2912c554080ae6",
				Rel:  "self",
			},
		},
	}

	if diff := deep.Equal(projectResponse, expected); diff != nil {
		t.Error(diff)
	}
	if !reflect.DeepEqual(projectResponse, expected) {
		t.Errorf("Projects.GetOneProject\n got=%#v\nwant=%#v", projectResponse, expected)
	}
}

func TestProject_Create(t *testing.T) {
	setup()
	defer teardown()

	createRequest := &Project{
		OrgID: "5a0a1e7e0f2912c554080adc",
		Name:  "ProjectFoobar",
	}

	mux.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"clusterCount": 2,
			"created": "2016-07-14T14:19:33Z",
			"id": "5a0a1e7e0f2912c554080ae6",
			"links": [{
				"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/5a0a1e7e0f2912c554080ae6",
				"rel": "self"
			}],
			"name": "ProjectFoobar",
			"orgId": "5a0a1e7e0f2912c554080adc"
		}`)
	})

	project, _, err := client.Projects.Create(ctx, createRequest)
	if err != nil {
		t.Errorf("Projects.Create returned error: %v", err)
	}

	expected := &Project{
		ClusterCount: 2,
		Created:      "2016-07-14T14:19:33Z",
		ID:           "5a0a1e7e0f2912c554080ae6",
		Links: []*Link{
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/5a0a1e7e0f2912c554080ae6",
				Rel:  "self",
			},
		},
		Name:  "ProjectFoobar",
		OrgID: "5a0a1e7e0f2912c554080adc",
	}

	if diff := deep.Equal(project, expected); diff != nil {
		t.Error(diff)
	}
	if !reflect.DeepEqual(project, expected) {
		t.Errorf("DatabaseUsers.Get\n got=%#v\nwant=%#v", project, expected)
	}
}

func TestProject_Delete(t *testing.T) {
	setup()
	defer teardown()

	projectID := "5a0a1e7e0f2912c554080adc"

	mux.HandleFunc(fmt.Sprintf("/groups/%s", projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Projects.Delete(ctx, projectID)
	if err != nil {
		t.Errorf("Projects.Delete returned error: %v", err)
	}
}

func TestProject_GetProjectTeamsAssigned(t *testing.T) {
	setup()
	defer teardown()

	projectID := "5a0a1e7e0f2912c554080adc"

	mux.HandleFunc(fmt.Sprintf("/%s/%s/teams", projectBasePath, projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"links": [
				{
					"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/{GROUP-ID}/teams",
					"rel": "self"
				}
			],
			"results": [
				{
					"links": [
						{
							"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/{GROUP-ID}/teams/{TEAM-ID}",
							"rel": "self"
						}
					],
					"roleNames": [
						"GROUP_READ_ONLY"
					],
					"teamId": "{TEAM-ID}"
				}
			],
			"totalCount": 1
		}`)
	})

	teamsAssigned, _, err := client.Projects.GetProjectTeamsAssigned(ctx, projectID)
	if err != nil {
		t.Errorf("Projects.GetProjectTeamsAssigned returned error: %v", err)
	}

	expected := &TeamsAssigned{
		Links: []*Link{
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/{GROUP-ID}/teams",
				Rel:  "self",
			},
		},
		Results: []*Result{
			{
				Links: []*Link{
					{
						Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/{GROUP-ID}/teams/{TEAM-ID}",
						Rel:  "self",
					},
				},
				RoleNames: []string{"GROUP_READ_ONLY"},
				TeamID:    "{TEAM-ID}",
			},
		},
		TotalCount: 1,
	}

	if diff := deep.Equal(teamsAssigned, expected); diff != nil {
		t.Error(diff)
	}

	if !reflect.DeepEqual(teamsAssigned, expected) {
		t.Errorf("Projects.GetProjectTeamsAssigned\n got=%+v\nwant=%+v", teamsAssigned, expected)
	}
}

func TestProject_AddTeamsToProject(t *testing.T) {
	setup()
	defer teardown()

	projectID := "5a0a1e7e0f2912c554080adc"

	createRequest := &ProjectTeam{
		TeamID: "{TEAM-ID}",
		Roles: []*RoleName{
			{RoleName: GROUP_OWNER},
			{RoleName: GROUP_READ_ONLY},
		},
	}

	mux.HandleFunc(fmt.Sprintf("/groups/%s/teams", projectID), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"links": [{
				"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/{GROUP-ID}/teams",
				"rel": "self"
			}],
			"results": [{
				"links": [{
					"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/{GROUP-ID}/teams/{TEAM-ID}",
					"rel": "self"
				}],
				"roleNames": ["GROUP_OWNER"],
				"teamId": "{TEAM-ID}"
			}],
			"totalCount": 1
		}`)
	})

	team, _, err := client.Projects.AddTeamsToProject(ctx, projectID, createRequest)
	if err != nil {
		t.Errorf("Projects.AddTeamsToProject returned error: %v", err)
	}

	expected := &TeamsAssigned{
		Links: []*Link{
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/{GROUP-ID}/teams",
				Rel:  "self",
			},
		},
		Results: []*Result{
			{
				Links: []*Link{
					{
						Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/{GROUP-ID}/teams/{TEAM-ID}",
						Rel:  "self",
					},
				},
				RoleNames: []string{"GROUP_OWNER"},
				TeamID:    "{TEAM-ID}",
			},
		},
		TotalCount: 1,
	}

	if diff := deep.Equal(team, expected); diff != nil {
		t.Error(diff)
	}
	if !reflect.DeepEqual(team, expected) {
		t.Errorf("DatabaseUsers.Get\n got=%#v\nwant=%#v", team, expected)
	}
}
