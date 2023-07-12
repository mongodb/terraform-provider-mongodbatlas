package provider

import (
	"os"
	"testing"
)

func testAccPreCheckBasic(tb testing.TB) {
	if os.Getenv("MONGODB_ATLAS_PUBLIC_KEY") == "" ||
		os.Getenv("MONGODB_ATLAS_PRIVATE_KEY") == "" ||
		os.Getenv("MONGODB_ATLAS_ORG_ID") == "" {
		tb.Fatal("`MONGODB_ATLAS_PUBLIC_KEY`, `MONGODB_ATLAS_PRIVATE_KEY`, and `MONGODB_ATLAS_ORG_ID` must be set for acceptance testing")
	}
}
