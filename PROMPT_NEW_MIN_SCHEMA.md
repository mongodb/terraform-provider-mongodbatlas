**Context:** This is the MongoDB Atlas Terraform provider repo. This is a *prototype* branch (for experimentation/demo, not production) building on a prior branch that added a `cluster_profile` attribute (`CORE`/`INFINITE`) with conditional autoscaling defaults. This second branch extends that to a **minimal-config experience**: the goal is that a user can deploy a working cluster by supplying only a few inputs (cluster name, provider region, and `cluster_profile`), with the profile supplying sensible defaults for everything else that's currently required. Correctness/edge-cases and effective-fields are explicitly out of scope — this is scaffolding to demo the minimal-config experience.

**Important — what this is NOT:** this does NOT involve effective fields, response-side propagation, computed-value surfacing. It is purely about the *input* side: making currently-required inputs optional and having the profile fill in defaults when they're omitted, so a minimal config deploys a full cluster.

**Task:**

1. **Start from the first branch** (the one adding `cluster_profile` \+ conditional autoscaling defaults). Build on top of it; don't redo that work.  
     
2. **Identify the currently-required inputs** on the cluster resource (the `mongodbatlas_advanced_cluster` resource — confirm by inspecting the schema). Report the list of required fields before changing anything, so we can confirm which ones to make optional.  
     
3. **Make those required inputs optional**, so the resource can be created with a minimal config. The target minimal config is roughly:

```
resource "mongodbatlas_advanced_cluster" "example" {
  name            = "my-cluster"
  provider_region = "AWS:US_EAST_1"
  cluster_profile = "INFINITE"
}
```

   - Change the relevant `Required: true` fields to `Optional: true` (except `project_id` and `name`).  
   - For each field made optional, supply a default value when the user omits it. Where the default depends on the profile, key it on `cluster_profile` (CORE defaults vs INFINITE defaults); where it's a universal sensible default, use a static one.  
   - Pick reasonable default values (e.g. a default instance size, default node count, default cloud provider) and **comment each default clearly** so it's obvious what was chosen and easy to change. This is a prototype — reasonable choices with clear comments are the goal, not "correct" production defaults. Only exception is `cluster_type` — that should be "REPLICASET" by default

   

4. **Precedence (critical):** explicit user input always wins. If the user provides a field, honor it; the profile default only fills the field when the user omitted it. (Same precedence rule as the autoscaling default in the first branch.)  
     
5. **Reverse-compat note:** making a field `Required → Optional` should be non-breaking for existing configs (configs that set the field still work). Verify that an existing-style config (one that sets all the fields explicitly) still plans cleanly with no unexpected diffs — i.e. you haven't changed behavior for users who specify everything. Include a test/example showing a "full explicit config" still works unchanged.  
     
6. **Implementation approach:** put the default-filling logic in the same plan-modification / default-resolution path used for the autoscaling default in the first branch — follow that existing pattern. Check whether the resource is on SDKv2 or the Plugin Framework, since how you implement optional-with-default differs between them (note: the framework disallows `Computed: true` \+ `Default` together; conditional defaults go in plan logic, not a schema `Default`).  
     
7. **Keep it minimal and forkable:** localized, well-commented changes so Alex can fork and adjust which fields are defaulted and to what values. Don't refactor surrounding code. Add a comment block where the minimal-config defaulting logic lives.  
     
8. **Add example Terraform configs** in the examples directory:  
     
   - **Minimal INFINITE config** — just name \+ region \+ `cluster_profile = "INFINITE"`, with a comment listing what defaults get filled in.  
   - **Minimal CORE config** — just name \+ region \+ `cluster_profile = "CORE"` (or omitted), with a comment on what defaults apply.  
   - **Full explicit config** — a config that sets all the (now-optional) fields explicitly, to show existing-style usage still works unchanged.

   

9. **Validation bar:** `terraform plan` is sufficient (no `apply` needed). Confirm via plan that (a) the minimal INFINITE config resolves to a full cluster spec with profile-driven defaults, (b) the minimal CORE config does likewise with CORE defaults, and (c) the full explicit config plans with no unexpected diffs vs. how it behaves today. Report the plan output for each.

**Out of scope (do not do):** don't add any new resources or new schema structures — only modify the existing cluster resource. Don't build any system for surfacing server-computed/resolved values back into state or output (this is input-side only). No multi-shard-specific behavior. No migration logic. Keep changes limited to making existing required inputs optional with profile-driven defaults.