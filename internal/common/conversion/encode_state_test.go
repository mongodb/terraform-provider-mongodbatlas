package conversion_test

import (
	"reflect"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func TestEncodeDecodeID(t *testing.T) {
	expected := map[string]string{
		"project_id":   "5cf5a45a9ccf6400e60981b6",
		"cluster_name": "test-acc-q4y272zo9y",
		"snapshot_id":  "5e42e646553855a5aee40138",
	}

	got := conversion.DecodeStateID(conversion.EncodeStateID(expected))

	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("Bad testEncodeDecodeID return \n got = %#v\nwant = %#v", got, expected)
	}
}

func TestDecodeID(t *testing.T) {
	expected := "Y2x1c3Rlcl9uYW1l:dGVzdC1hY2MtcTR5Mjcyem85eQ==-c25hcHNob3RfaWQ=:NWU0MmU2NDY1NTM4NTVhNWFlZTQwMTM4-cHJvamVjdF9pZA==:NWNmNWE0NWE5Y2NmNjQwMGU2MDk4MWI2"
	expected2 := "c25hcHNob3RfaWQ=:NWU0MmU2NDY1NTM4NTVhNWFlZTQwMTM4-cHJvamVjdF9pZA==:NWNmNWE0NWE5Y2NmNjQwMGU2MDk4MWI2-Y2x1c3Rlcl9uYW1l:dGVzdC1hY2MtcTR5Mjcyem85eQ=="

	got := conversion.DecodeStateID(expected)
	got2 := conversion.DecodeStateID(expected2)

	if !reflect.DeepEqual(got, got2) {
		t.Fatalf("Bad TestDecodeID return \n got = %#v\nwant = %#v", got, got2)
	}
}
