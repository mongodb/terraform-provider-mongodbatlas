package mongodbatlas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-test/deep"
	"github.com/mwielbut/pointy"
)

func TestMaintenanceWindows_UpdateWithSheduleTime(t *testing.T) {
	setup()
	defer teardown()

	groupID := "6d2065c687d9d64ae7acdg41"

	mux.HandleFunc(fmt.Sprintf("/"+maintenanceWindowsPath, groupID), func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"dayOfWeek":         float64(2),
			"hourOfDay":         float64(3),
			"numberOfDeferrals": float64(4),
			"startASAP":         false,
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if diff := deep.Equal(v, expected); diff != nil {
			t.Errorf("Request body\n got=%#v\nwant=%#v \n\ndiff=%#v", v, expected, diff)
		}

		fmt.Fprint(w, `{}`)
	})

	maintenanceRequest := &MaintenanceWindow{
		DayOfWeek:         2,
		HourOfDay:         pointy.Int(3),
		NumberOfDeferrals: 4,
		StartASAP:         pointy.Bool(false),
	}

	_, err := client.MaintenanceWindows.Update(ctx, groupID, maintenanceRequest)
	if err != nil {
		t.Errorf("MaintenanceWindow.Update returned error: %v", err)
		return
	}
}

func TestMaintenanceWindows_UpdateWithStartNow(t *testing.T) {
	setup()
	defer teardown()

	groupID := "6d2065c687d9d64ae7acdg41"

	mux.HandleFunc(fmt.Sprintf("/"+maintenanceWindowsPath, groupID), func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"startASAP": true,
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if diff := deep.Equal(v, expected); diff != nil {
			t.Errorf("Request body\n got=%#v\nwant=%#v \n\ndiff=%#v", v, expected, diff)
		}

		fmt.Fprint(w, `{}`)
	})

	maintenanceRequest := &MaintenanceWindow{
		StartASAP: pointy.Bool(true),
	}

	_, err := client.MaintenanceWindows.Update(ctx, groupID, maintenanceRequest)
	if err != nil {
		t.Errorf("MaintenanceWindow.Update returned error: %v", err)
		return
	}
}

func TestMaintenanceWindows_Get(t *testing.T) {
	setup()
	defer teardown()

	groupID := "6d2065c687d9d64ae7acdg41"

	mux.HandleFunc(fmt.Sprintf("/"+maintenanceWindowsPath, groupID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, fmt.Sprintf(`{
				"dayOfWeek": 2,
				"hourOfDay": 3,
				"numberOfDeferrals": 4
		}`))
	})

	maintenanceWindow, _, err := client.MaintenanceWindows.Get(ctx, groupID)
	if err != nil {
		t.Errorf("MaintenanceWindows.Get returned error: %v", err)
	}

	expected := &MaintenanceWindow{
		DayOfWeek:         2,
		HourOfDay:         pointy.Int(3),
		NumberOfDeferrals: 4,
	}

	if diff := deep.Equal(maintenanceWindow, expected); diff != nil {
		t.Errorf("Request body\n got=%#v\nwant=%#v \n\ndiff=%#v", maintenanceWindow, expected, diff)
	}
}

func TestMaintenanceWindows_Defer(t *testing.T) {
	setup()
	defer teardown()

	groupID := "6d2065c687d9d64ae7acdg41"

	mux.HandleFunc(fmt.Sprintf("/"+maintenanceWindowsPath+"/defer", groupID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{}`)
	})

	_, err := client.MaintenanceWindows.Defer(ctx, groupID)
	if err != nil {
		t.Errorf("MaintenanceWindows.Defer returned error: %v", err)
	}
}

func TestMaintenanceWindows_Delete(t *testing.T) {
	setup()
	defer teardown()

	groupID := "6d2065c687d9d64ae7acdg41"

	mux.HandleFunc(fmt.Sprintf("/"+maintenanceWindowsPath, groupID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		fmt.Fprint(w, `{}`)
	})

	_, err := client.MaintenanceWindows.Reset(ctx, groupID)
	if err != nil {
		t.Errorf("MaintenanceWindows.Get returned error: %v", err)
	}
}
