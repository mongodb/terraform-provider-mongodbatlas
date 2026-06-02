# Atlas SDK preview (dev-latest)

The provider ships two Atlas Admin API clients:

- **Versioned SDK** (`go.mongodb.org/atlas-sdk/v<YYYYMMDD>/admin`): default for production resources. Access via `MongoDBClient.AtlasV2`.
- **Preview SDK** (`github.com/mongodb/atlas-sdk-go/admin`): generated from [release candidate](https://github.com/mongodb/openapi) OpenAPI on `mongodb/openapi` branch `dev`, published on the atlas-sdk-go branch `dev-latest`. Access via `MongoDBClient.AtlasPreview`.

Use the preview client only when an endpoint exists in the RC API but is not yet available in the pinned versioned module. Move calls to `AtlasV2` after the API ships in a tagged SDK release.

For versioned SDK bumps and major-release procedure, see [Atlas SDK](atlas-sdk.md).

## Trigger a dev preview SDK build

The [mongodb/atlas-sdk-go](https://github.com/mongodb/atlas-sdk-go) repository publishes preview code through the **Generate Preview SDK** workflow:

- **Manual**: [Actions → Generate Preview SDK → Run workflow](https://github.com/mongodb/atlas-sdk-go/actions/workflows/autoupdate-preview.yaml)
- **Scheduled**: Monday, Wednesday, and Friday (08:30 UTC)

The workflow fetches OpenAPI from `mongodb/openapi` on the `dev` branch. It regenerates the client only when `openapi/atlas-api.yaml` changes, then opens or updates a PR on branch `dev-latest` (*APIBot: GO SDK Dev Preview*).

**Do not merge the preview PR to `main`.** It is a living branch for early consumers, not a release.

After the workflow succeeds, dependents pin the module pseudo-version:

```bash
go get github.com/mongodb/atlas-sdk-go@dev-latest
go mod tidy
```

More detail: [SDK releaser — SDK Preview](https://github.com/mongodb/atlas-sdk-go/blob/main/tools/releaser/README.md#sdk-preview) and the workflow PR body in [autoupdate-preview.yaml](https://github.com/mongodb/atlas-sdk-go/blob/main/.github/workflows/autoupdate-preview.yaml).

## Update preview SDK in the provider

**Preview only** (from the provider repository root):

```bash
go get github.com/mongodb/atlas-sdk-go@dev-latest
go mod tidy
make fix
```

**Versioned and preview** (major versioned bump plus preview refresh):

```bash
make update-atlas-sdk
```

`scripts/update-sdk.sh` rewrites the latest `go.mongodb.org/atlas-sdk/v<YYYYMMDD>` module, then runs `go get github.com/mongodb/atlas-sdk-go@dev-latest`.

> NOTE: `make update-atlas-sdk` can change imports across hundreds of files. Run it on a clean `master` checkout with no uncommitted changes.

## Use the preview client in resource code

Both clients share the same HTTP transport and credentials configured in `internal/config/client.go`.

### Terraform Plugin Framework

Embed `config.RSCommon` and call APIs on `r.Client.AtlasPreview`:

```go
import adminpreview "github.com/mongodb/atlas-sdk-go/admin"

func (r *RS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	conn := r.Client.AtlasPreview // github.com/mongodb/atlas-sdk-go/admin
	// conn.SomeApi.SomeMethod(ctx).Execute()
}
```

Stable resources continue to use `r.Client.AtlasV2` with `go.mongodb.org/atlas-sdk/v<YYYYMMDD>/admin`.

### SDKv2 resources

```go
conn := meta.(*config.MongoDBClient).AtlasPreview
```

### Service account OAuth

OAuth for service accounts uses `github.com/mongodb/atlas-sdk-go/auth` in `internal/config/service_account.go`. That package is separate from `AtlasPreview` admin API calls.

## Acceptance tests

`internal/testutil/acc` exposes `ConnPreview()` for direct Atlas calls in acceptance tests (same client as `AtlasPreview` on the configured provider):

```go
conn := acc.ConnPreview()
```

## When to use which client

- **Versioned SDK**: Import `go.mongodb.org/atlas-sdk/v<YYYYMMDD>/admin`. Use `AtlasV2` by default; it matches released SDK tags and provider dependencies.
- **Preview SDK**: Import `github.com/mongodb/atlas-sdk-go/admin`. Use `AtlasPreview` for RC-only endpoints until the next versioned SDK release.
- **Legacy pins**: `AtlasV220240530` and `AtlasV220241113` remain for specific backward-compatibility cases documented in `internal/config/client.go`.
