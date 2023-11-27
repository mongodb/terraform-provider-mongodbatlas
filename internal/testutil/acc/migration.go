package acc

import (
	"os"
	"testing"
)

func PreCheckMigration(tb testing.TB) {
	PreCheck(tb)
	if os.Getenv("MONGODB_ATLAS_LAST_VERSION") == "" {
		tb.Fatal("`MONGODB_ATLAS_LAST_VERSION` must be set for migration acceptance testing")
	}
}

func PreCheckBasicMigration(tb testing.TB) {
	PreCheckBasic(tb)
	if os.Getenv("MONGODB_ATLAS_LAST_VERSION") == "" {
		tb.Fatal("`MONGODB_ATLAS_LAST_VERSION` must be set for migration acceptance testing")
	}
}

func PreCheckBasicOwnerIDMigration(tb testing.TB) {
	PreCheckBasicMigration(tb)
	if os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID") == "" {
		tb.Fatal("`MONGODB_ATLAS_PROJECT_OWNER_ID` must be set ")
	}
}
