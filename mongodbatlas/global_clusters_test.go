package mongodbatlas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

func TestGlobalClusters_Get(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1"
	clusterName := "appData"

	mux.HandleFunc(fmt.Sprintf("/groups/%s/clusters/%s/globalWrites", groupID, clusterName), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, ` {
			"customZoneMapping" : {
			  "AF" : "5b48f5a780eef5236f689f94",
			  "AL" : "5b48f5a780eef5236f689f94",
			  "AU" : "5b48f5cddff5220000f7f375",
			  "AU-ACT" : "5b48f5cddff5220000f7f375",
			  "AU-NSW" : "5b48f5cddff5220000f7f375",
			  "AU-NT" : "5b48f5cddff5220000f7f375",
			  "AU-QLD" : "5b48f5cddff5220000f7f375",
			  "AU-SA" : "5b48f5cddff5220000f7f375",
			  "AU-TAS" : "5b48f5cddff5220000f7f375",
			  "AU-VIC" : "5b48f5cddff5220000f7f375",
			  "AU-WA" : "5b48f5cddff5220000f7f375",
			  "DZ" : "5b48f5a780eef5236f689f94"
		 },
			"managedNamespaces" : [ {
			  "collection" : "zips",
			  "customShardKey" : "city",
			  "db" : "mydata"
			},{
			  "collection" : "stores",
			  "customShardKey" : "store_number",
			  "db" : "mydata"
			} ]
		  }`)
	})

	cluster, _, err := client.GlobalClusters.Get(ctx, groupID, clusterName)
	if err != nil {
		t.Errorf("GlobalClusters.Get returned error: %v", err)
	}

	expected := &GlobalCluster{
		CustomZoneMapping: map[string]string{
			"AF":     "5b48f5a780eef5236f689f94",
			"AL":     "5b48f5a780eef5236f689f94",
			"AU":     "5b48f5cddff5220000f7f375",
			"AU-ACT": "5b48f5cddff5220000f7f375",
			"AU-NSW": "5b48f5cddff5220000f7f375",
			"AU-NT":  "5b48f5cddff5220000f7f375",
			"AU-QLD": "5b48f5cddff5220000f7f375",
			"AU-SA":  "5b48f5cddff5220000f7f375",
			"AU-TAS": "5b48f5cddff5220000f7f375",
			"AU-VIC": "5b48f5cddff5220000f7f375",
			"AU-WA":  "5b48f5cddff5220000f7f375",
			"DZ":     "5b48f5a780eef5236f689f94",
		},
		ManagedNamespaces: []ManagedNamespace{
			{
				Collection:     "zips",
				CustomShardKey: "city",
				Db:             "mydata",
			}, {
				Collection:     "stores",
				CustomShardKey: "store_number",
				Db:             "mydata",
			},
		},
	}

	if !reflect.DeepEqual(cluster, expected) {
		t.Errorf("GlobalClusters.Get\n got=%#v\nwant=%#v", cluster, expected)
	}
}

func TestGlobalClusters_AddManagedNamespace(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1"
	clusterName := "appData"

	createRequest := &ManagedNamespace{
		Db:             "mydata",
		Collection:     "publishers",
		CustomShardKey: "city",
	}

	mux.HandleFunc(fmt.Sprintf("/groups/%s/clusters/%s/globalWrites/managedNamespaces", groupID, clusterName), func(w http.ResponseWriter, r *http.Request) {
		expectedRequest := map[string]interface{}{
			"db":             "mydata",
			"collection":     "publishers",
			"customShardKey": "city",
		}

		jsonBlob := `
		{
			"customZoneMapping" : {
			  "AF" : "5b48f5a780eef5236f689f94",
			  "AL" : "5b48f5a780eef5236f689f94",
			  "AU" : "5b48f5cddff5220000f7f375",
			  "AU-ACT" : "5b48f5cddff5220000f7f375",
			  "AU-NSW" : "5b48f5cddff5220000f7f375",
			  "AU-NT" : "5b48f5cddff5220000f7f375",
			  "AU-QLD" : "5b48f5cddff5220000f7f375",
			  "AU-SA" : "5b48f5cddff5220000f7f375",
			  "AU-TAS" : "5b48f5cddff5220000f7f375",
			  "AU-VIC" : "5b48f5cddff5220000f7f375",
			  "AU-WA" : "5b48f5cddff5220000f7f375",
			  "DZ" : "5b48f5a780eef5236f689f94"
		 },
			"managedNamespaces" : [ {
			  "collection" : "publishers",
			  "customShardKey" : "city",
			  "db" : "mydata"
			},{
			  "collection" : "stores",
			  "customShardKey" : "store_number",
			  "db" : "mydata"
			} ]
		  }
		`

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}
		if diff := deep.Equal(v, expectedRequest); diff != nil {
			t.Errorf("GlobalClusters.AddManagedNamespace Request Body = %v", diff)
		}

		if !reflect.DeepEqual(v, expectedRequest) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expectedRequest)
		}

		fmt.Fprint(w, jsonBlob)
	})

	globalCluster, _, err := client.GlobalClusters.AddManagedNamespace(ctx, groupID, clusterName, createRequest)
	if err != nil {
		t.Errorf("GlobalClusters.AddManagedNamespace returned error: %v", err)
	}

	if namespacesCount := len(globalCluster.ManagedNamespaces); namespacesCount != 2 {
		t.Errorf("expected name '%d', received '%d'", 2, namespacesCount)
	}

	expectedCustomZoneMapping := map[string]string{
		"AF":     "5b48f5a780eef5236f689f94",
		"AL":     "5b48f5a780eef5236f689f94",
		"AU":     "5b48f5cddff5220000f7f375",
		"AU-ACT": "5b48f5cddff5220000f7f375",
		"AU-NSW": "5b48f5cddff5220000f7f375",
		"AU-NT":  "5b48f5cddff5220000f7f375",
		"AU-QLD": "5b48f5cddff5220000f7f375",
		"AU-SA":  "5b48f5cddff5220000f7f375",
		"AU-TAS": "5b48f5cddff5220000f7f375",
		"AU-VIC": "5b48f5cddff5220000f7f375",
		"AU-WA":  "5b48f5cddff5220000f7f375",
		"DZ":     "5b48f5a780eef5236f689f94",
	}

	if diff := deep.Equal(globalCluster.CustomZoneMapping, expectedCustomZoneMapping); diff != nil {
		t.Errorf("expected globalCluster.CustomZoneMapping = %v", diff)
	}

}

func TestGlobalClusters_DeleteManagedNamespace(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1"
	name := "appData"

	mn := ManagedNamespace{
		Db:         "data",
		Collection: "distributors",
	}

	mux.HandleFunc(fmt.Sprintf("/groups/%s/clusters/%s/globalWrites/managedNamespaces", groupID, name), func(w http.ResponseWriter, r *http.Request) {

		if collection := r.URL.Query().Get("collection"); collection != mn.Collection {
			t.Errorf("expected query param collection = '%s', received '%s'", mn.Collection, collection)
		}

		if db := r.URL.Query().Get("db"); db != mn.Db {
			t.Errorf("expected query param db = '%s', received '%s'", mn.Collection, db)
		}

		jsonBlob := `
		{
			"customZoneMapping" : {
			  "AF" : "5b48f5a780eef5236f689f94",
			  "AL" : "5b48f5a780eef5236f689f94",
			  "AU" : "5b48f5cddff5220000f7f375",
			  "AU-ACT" : "5b48f5cddff5220000f7f375",
			  "AU-NSW" : "5b48f5cddff5220000f7f375",
			  "AU-NT" : "5b48f5cddff5220000f7f375",
			  "AU-QLD" : "5b48f5cddff5220000f7f375",
			  "AU-SA" : "5b48f5cddff5220000f7f375",
			  "AU-TAS" : "5b48f5cddff5220000f7f375",
			  "AU-VIC" : "5b48f5cddff5220000f7f375",
			  "AU-WA" : "5b48f5cddff5220000f7f375",
			  "DZ" : "5b48f5a780eef5236f689f94"
		 },
			"managedNamespaces" : [ {
			  "collection" : "zip-codes",
			  "customShardKey" : "city",
			  "db" : "data"
			}, {
			  "collection" : "retail-stores",
			  "customShardKey" : "store-number",
			  "db" : "mydata"
			} ]
		  }
		`
		testMethod(t, r, http.MethodDelete)

		fmt.Fprint(w, jsonBlob)
	})

	globalCluster, _, err := client.GlobalClusters.DeleteManagedNamespace(ctx, groupID, name, &mn)

	if err != nil {
		t.Errorf("GlobalClusters.DeleteManagedNamespace returned error: %v", err)
	}

	expected := &GlobalCluster{
		CustomZoneMapping: map[string]string{
			"AF":     "5b48f5a780eef5236f689f94",
			"AL":     "5b48f5a780eef5236f689f94",
			"AU":     "5b48f5cddff5220000f7f375",
			"AU-ACT": "5b48f5cddff5220000f7f375",
			"AU-NSW": "5b48f5cddff5220000f7f375",
			"AU-NT":  "5b48f5cddff5220000f7f375",
			"AU-QLD": "5b48f5cddff5220000f7f375",
			"AU-SA":  "5b48f5cddff5220000f7f375",
			"AU-TAS": "5b48f5cddff5220000f7f375",
			"AU-VIC": "5b48f5cddff5220000f7f375",
			"AU-WA":  "5b48f5cddff5220000f7f375",
			"DZ":     "5b48f5a780eef5236f689f94",
		},
		ManagedNamespaces: []ManagedNamespace{
			{
				Collection:     "zip-codes",
				CustomShardKey: "city",
				Db:             "data",
			}, {
				Collection:     "retail-stores",
				CustomShardKey: "store-number",
				Db:             "mydata",
			},
		},
	}

	if diff := deep.Equal(globalCluster, expected); diff != nil {
		t.Errorf("GlobalClusters.DeleteCustomZoneMappings Reponse Body = %v", diff)
	}
}

func TestGlobalClusters_AddCustomZoneMappings(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1"
	clusterName := "appData"

	createRequest := &CustomZoneMappingsRequest{
		CustomZoneMappings: []CustomZoneMapping{
			{
				Location: "CA",
				Zone:     "Zone 1",
			},
		},
	}

	mux.HandleFunc(fmt.Sprintf("/groups/%s/clusters/%s/globalWrites/customZoneMapping", groupID, clusterName), func(w http.ResponseWriter, r *http.Request) {
		expectedRequest := map[string]interface{}{
			"customZoneMappings": []interface{}{
				map[string]interface{}{"location": "CA", "zone": "Zone 1"},
			},
		}

		jsonBlob := `
		{
			"customZoneMapping" : {
			   "CA" : "5b50bf4180eef547653df4d0"
			},
			"managedNamespaces" : [ {
			   "collection" : "publishers",
			   "customShardKey" : "city",
			   "db" : "myData"
			} ]
		 }
		`

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}
		if diff := deep.Equal(v, expectedRequest); diff != nil {
			t.Errorf("GlobalClusters.AddCustomZoneMappings Request Body = %v", diff)
		}

		if !reflect.DeepEqual(v, expectedRequest) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expectedRequest)
		}

		fmt.Fprint(w, jsonBlob)
	})

	globalCluster, _, err := client.GlobalClusters.AddCustomZoneMappings(ctx, groupID, clusterName, createRequest)
	if err != nil {
		t.Errorf("GlobalClusters.AddCustomZoneMappings returned error: %v", err)
	}

	if namespacesCount := len(globalCluster.ManagedNamespaces); namespacesCount != 1 {
		t.Errorf("expected name '%d', received '%d'", 1, namespacesCount)
	}

	if id := globalCluster.CustomZoneMapping["CA"]; id != "5b50bf4180eef547653df4d0" {
		t.Errorf("expected CustomZoneMapping.CA '%s', received '%s'", "5b50bf4180eef547653df4d0", id)
	}

}

func TestGlobalClusters_DeleteCustomZoneMappings(t *testing.T) {
	setup()
	defer teardown()

	groupID := "1"
	name := "appData"

	mux.HandleFunc(fmt.Sprintf("/groups/%s/clusters/%s/globalWrites/customZoneMapping", groupID, name), func(w http.ResponseWriter, r *http.Request) {

		jsonBlob := `
		{
			"customZoneMapping" : { },
			"managedNamespaces" : [ {
			  "collection" : "publishers",
			  "customShardKey" : "city",
			  "db" : "myData"
			} ]
		  }
		`
		testMethod(t, r, http.MethodDelete)

		fmt.Fprint(w, jsonBlob)
	})

	globalCluster, _, err := client.GlobalClusters.DeleteCustomZoneMappings(ctx, groupID, name)

	if err != nil {
		t.Errorf("Cluster.Delete returned error: %v", err)
	}

	expected := &GlobalCluster{
		CustomZoneMapping: map[string]string{},
		ManagedNamespaces: []ManagedNamespace{
			{
				Collection:     "publishers",
				CustomShardKey: "city",
				Db:             "myData",
			},
		},
	}

	if diff := deep.Equal(globalCluster, expected); diff != nil {
		t.Errorf("GlobalClusters.AddCustomZoneMappings Reponse Body = %v", diff)
	}
}
