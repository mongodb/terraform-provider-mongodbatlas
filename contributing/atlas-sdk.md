Most of the new features of the provider are using [atlas-sdk](https://github.com/mongodb/atlas-sdk-go)
SDK is updated automatically, tracking all new Atlas features.

### Updating Atlas SDK 

To update Atlas SDK run:

```bash
make update-atlas-sdk
```

> NOTE: The update mechanism is only needed for major releases. Any other releases will be supported by dependabot.

> NOTE: Command can make import changes to +500 files. Please make sure that you perform update on main branch without any uncommited changes.

### SDK Major Release Update Procedure

1. If the SDK update doesn’t cause any compilation issues create a new SDK update PR
   1. Review [API Changelog](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/changelog) for any deprecated fields and breaking changes.
2. For SDK updates introducing compilation issues without graceful workaround
   1. Use the previous major version of the SDK (including the old client) for the affected resource
   1. Create an issue to identify the root cause and mitigation paths based on changelog information  
   2. If applicable: Make required notice/update to the end users based on the plan.
