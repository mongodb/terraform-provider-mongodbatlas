package mongodbatlas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/mwielbut/pointy"
)

func TestEncryptionsAtRest_Create(t *testing.T) {
	setup()
	defer teardown()

	createRequest := &EncryptionAtRest{
		GroupID: "6d2065c687d9d64ae7acdg41",
		AwsKms: AwsKms{
			Enabled:             pointy.Bool(true),
			AccessKeyID:         "AKIAIOSFODNN7EXAMPLE",
			SecretAccessKey:     "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			CustomerMasterKeyID: "030gce02-586d-48d2-a966-05ea954fde0g",
			Region:              CaCentral1,
		},
		AzureKeyVault: AzureKeyVault{
			Enabled:           pointy.Bool(true),
			ClientID:          "g54f9e2-89e3-40fd-8188-EXAMPLEID",
			AzureEnvironment:  Azure,
			SubscriptionID:    "0ec944e3-g725-44f9-a147-EXAMPLEID",
			ResourceGroupName: "ExampleRGName",
			KeyVaultName:      "EXAMPLEKeyVault",
			KeyIdentifier:     "https://EXAMPLEKeyVault.vault.azure.net/keys/EXAMPLEKey/d891821e3d364e9eb88fbd3d11807b86",
			Secret:            "EXAMPLESECRET",
			TenantID:          "e8e4b6ba-ff32-4c88-a9af-EXAMPLEID",
		},
		GoogleCloudKms: GoogleCloudKms{
			Enabled:              pointy.Bool(true),
			ServiceAccountKey:    "{\"type\": \"service_account\",\"project_id\": \"my-project-common-0\",\"private_key_id\": \"e120598ea4f88249469fcdd75a9a785c1bb3\",\"private_key\": \"-----BEGIN PRIVATE KEY-----\\nMIIEuwIBA(truncated)SfecnS0mT94D9\\n-----END PRIVATE KEY-----\\n\",\"client_email\": \"my-email-kms-0@my-project-common-0.iam.gserviceaccount.com\",\"client_id\": \"10180967717292066\",\"auth_uri\": \"https://accounts.google.com/o/oauth2/auth\",\"token_uri\": \"https://accounts.google.com/o/oauth2/token\",\"auth_provider_x509_cert_url\": \"https://www.googleapis.com/oauth2/v1/certs\",\"client_x509_cert_url\": \"https://www.googleapis.com/robot/v1/metadata/x509/my-email-kms-0%40my-project-common-0.iam.gserviceaccount.com\"}",
			KeyVersionResourceID: "projects/my-project-common-0/locations/us-east4/keyRings/my-key-ring-0/cryptoKeys/my-key-0/cryptoKeyVersions/1",
		},
	}

	mux.HandleFunc(fmt.Sprintf("/"+encryptionsAtRestBasePath, createRequest.GroupID), func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"awsKms": map[string]interface{}{
				"enabled":             true,
				"accessKeyID":         "AKIAIOSFODNN7EXAMPLE",
				"secretAccessKey":     "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				"customerMasterKeyID": "030gce02-586d-48d2-a966-05ea954fde0g",
				"region":              CaCentral1,
			},
			"azureKeyVault": map[string]interface{}{
				"enabled":           true,
				"clientID":          "g54f9e2-89e3-40fd-8188-EXAMPLEID",
				"azureEnvironment":  Azure,
				"subscriptionID":    "0ec944e3-g725-44f9-a147-EXAMPLEID",
				"resourceGroupName": "ExampleRGName",
				"keyVaultName":      "EXAMPLEKeyVault",
				"keyIdentifier":     "https://EXAMPLEKeyVault.vault.azure.net/keys/EXAMPLEKey/d891821e3d364e9eb88fbd3d11807b86",
				"secret":            "EXAMPLESECRET",
				"tenantID":          "e8e4b6ba-ff32-4c88-a9af-EXAMPLEID",
			},
			"googleCloudKms": map[string]interface{}{
				"enabled":              true,
				"serviceAccountKey":    "{\"type\": \"service_account\",\"project_id\": \"my-project-common-0\",\"private_key_id\": \"e120598ea4f88249469fcdd75a9a785c1bb3\",\"private_key\": \"-----BEGIN PRIVATE KEY-----\\nMIIEuwIBA(truncated)SfecnS0mT94D9\\n-----END PRIVATE KEY-----\\n\",\"client_email\": \"my-email-kms-0@my-project-common-0.iam.gserviceaccount.com\",\"client_id\": \"10180967717292066\",\"auth_uri\": \"https://accounts.google.com/o/oauth2/auth\",\"token_uri\": \"https://accounts.google.com/o/oauth2/token\",\"auth_provider_x509_cert_url\": \"https://www.googleapis.com/oauth2/v1/certs\",\"client_x509_cert_url\": \"https://www.googleapis.com/robot/v1/metadata/x509/my-email-kms-0%40my-project-common-0.iam.gserviceaccount.com\"}",
				"keyVersionResourceID": "projects/my-project-common-0/locations/us-east4/keyRings/my-key-ring-0/cryptoKeys/my-key-0/cryptoKeyVersions/1",
			},
		}

		var v map[string]interface{}
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
			"awsKms": {
				"enabled": true,
				"accessKeyID": "AKIAIOSFODNN7EXAMPLE",
				"customerMasterKeyID": "030gce02-586d-48d2-a966-05ea954fde0g",
				"region": "US_EAST_1"
			},
			"azureKeyVault": {
				"enabled": true,
				"clientID": "g54f9e2-89e3-40fd-8188-EXAMPLEID",
				"azureEnvironment": "AZURE",
				"subscriptionID": "0ec944e3-g725-44f9-a147-EXAMPLEID",
				"resourceGroupName": "ExampleRGName",
				"keyVaultName": "EXAMPLEKeyVault",
				"keyIdentifier": "https://EXAMPLEKeyVault.vault.azure.net/keys/EXAMPLEKey/d891821e3d364e9eb88fbd3d11807b86",
				"tenantID": "e8e4b6ba-ff32-4c88-a9af-dc17efegdf63"
			},
			"googleCloudKms": {
				"enabled": true,
				"keyVersionResourceID": "projects/my-project-common-0/locations/us-east4/keyRings/my-key-ring-0/cryptoKeys/my-key-0/cryptoKeyVersions/1"
			}
		}`)
	})

	cloudProviderSnapshot, _, err := client.EncryptionsAtRest.Create(ctx, createRequest)
	if err != nil {
		t.Errorf("EncryptionsAtRest.Create returned error: %v", err)
	}

	expected := &EncryptionAtRest{
		AwsKms: AwsKms{
			Enabled:             pointy.Bool(true),
			AccessKeyID:         "AKIAIOSFODNN7EXAMPLE",
			CustomerMasterKeyID: "030gce02-586d-48d2-a966-05ea954fde0g",
			Region:              "US_EAST_1",
		},
		AzureKeyVault: AzureKeyVault{
			Enabled:           pointy.Bool(true),
			ClientID:          "g54f9e2-89e3-40fd-8188-EXAMPLEID",
			AzureEnvironment:  "AZURE",
			SubscriptionID:    "0ec944e3-g725-44f9-a147-EXAMPLEID",
			ResourceGroupName: "ExampleRGName",
			KeyVaultName:      "EXAMPLEKeyVault",
			KeyIdentifier:     "https://EXAMPLEKeyVault.vault.azure.net/keys/EXAMPLEKey/d891821e3d364e9eb88fbd3d11807b86",
			TenantID:          "e8e4b6ba-ff32-4c88-a9af-dc17efegdf63",
		},
		GoogleCloudKms: GoogleCloudKms{
			Enabled:              pointy.Bool(true),
			KeyVersionResourceID: "projects/my-project-common-0/locations/us-east4/keyRings/my-key-ring-0/cryptoKeys/my-key-0/cryptoKeyVersions/1",
		},
	}

	if diff := deep.Equal(cloudProviderSnapshot, expected); diff != nil {
		t.Error(diff)
	}
	if !reflect.DeepEqual(cloudProviderSnapshot, expected) {
		t.Errorf("EncryptionsAtRest.Create\n got=%#v\nwant=%#v", cloudProviderSnapshot, expected)
	}
}

func TestEncryptionsAtRest_Get(t *testing.T) {
	setup()
	defer teardown()

	groupID := "6d2065c687d9d64ae7acdg41"

	mux.HandleFunc(fmt.Sprintf("/"+encryptionsAtRestBasePath, groupID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"awsKms": {
				"enabled": true,
				"accessKeyID": "AKIAIOSFODNN7EXAMPLE",
				"customerMasterKeyID": "030gce02-586d-48d2-a966-05ea954fde0g",
				"region": "US_EAST_1"
			},
			"azureKeyVault": {
				"enabled": true,
				"azureEnvironment": "AZURE",
				"clientID": "g54f9e2-89e3-40fd-8188-EXAMPLEID",
				"keyIdentifier": "https://EXAMPLEKeyVault.vault.azure.net/keys/EXAMPLEKey/d891821e3d364e9eb88fbd3d11807b86",
				"keyVaultName": "EXAMPLEKeyVault",
				"resourceGroupName": "ExampleRGName",
				"subscriptionID": "0ec944e3-g725-44f9-a147-EXAMPLEID",
				"tenantID": "e8e4b6ba-ff32-4c88-a9af-EXAMPLEID"
			},
			"googleCloudKms": {
				"enabled": true,
				"keyVersionResourceID": "projects/my-project/locations/us-east4/keyRings/my-key-ring-0/cryptoKeys/my-key-0/cryptoKeyVersions/1"
			}
		}`)
	})

	cloudProviderSnapshot, _, err := client.EncryptionsAtRest.Get(ctx, groupID)
	if err != nil {
		t.Errorf("EncryptionsAtRest.Get returned error: %v", err)
	}

	expected := &EncryptionAtRest{
		AwsKms: AwsKms{
			Enabled:             pointy.Bool(true),
			AccessKeyID:         "AKIAIOSFODNN7EXAMPLE",
			CustomerMasterKeyID: "030gce02-586d-48d2-a966-05ea954fde0g",
			Region:              "US_EAST_1",
		},
		AzureKeyVault: AzureKeyVault{
			Enabled:           pointy.Bool(true),
			AzureEnvironment:  "AZURE",
			ClientID:          "g54f9e2-89e3-40fd-8188-EXAMPLEID",
			KeyIdentifier:     "https://EXAMPLEKeyVault.vault.azure.net/keys/EXAMPLEKey/d891821e3d364e9eb88fbd3d11807b86",
			KeyVaultName:      "EXAMPLEKeyVault",
			ResourceGroupName: "ExampleRGName",
			SubscriptionID:    "0ec944e3-g725-44f9-a147-EXAMPLEID",
			TenantID:          "e8e4b6ba-ff32-4c88-a9af-EXAMPLEID",
		},
		GoogleCloudKms: GoogleCloudKms{
			Enabled:              pointy.Bool(true),
			KeyVersionResourceID: "projects/my-project/locations/us-east4/keyRings/my-key-ring-0/cryptoKeys/my-key-0/cryptoKeyVersions/1",
		},
	}

	if diff := deep.Equal(cloudProviderSnapshot, expected); diff != nil {
		t.Error(diff)
	}
	if !reflect.DeepEqual(cloudProviderSnapshot, expected) {
		t.Errorf("EncryptionsAtRest.Get\n got=%#v\nwant=%#v", cloudProviderSnapshot, expected)
	}
}
