**Context:** This is the MongoDB Atlas Terraform provider repo. I need a *prototype* branch (for experimentation/demo, not production) that adds a new `cluster_profile` attribute to the cluster resource, so it can be forked to try different default-behavior permutations. Correctness/edge-cases are explicitly out of scope — this is scaffolding to demo how a profile-driven default would look.

**Task:**

1. **Add a new optional attribute `cluster_profile`** to the cluster resource schema (the `mongodbatlas_advanced_cluster` resource — confirm the correct resource file first by inspecting the schema definitions).

   * Type: string.  
   * Allowed values: `CORE` and `INFINITE`.  
   * Optional, with default `CORE`.  
   * Add a validator restricting it to those two values.  
2. **Wire the conditional default behavior** keyed on `cluster_profile`:

   * When `cluster_profile = CORE` (or unset): no behavior change — clusters behave exactly as they do today. This is the baseline.  
   * When `cluster_profile = INFINITE`: change the autoscaling defaults so that, *when the user has not explicitly configured autoscaling*:  
     * compute auto-scaling is enabled (both scale-up and scale-down on),  
     * `min_instance_size` \= the cluster's configured instance size,  
     * `max_instance_size` \= the instance size two tiers above the configured instance size (i.e. instance\_size \+ 2 in the M-tier progression — e.g. M30 → max M50). Implement a small helper that, given an instance size, returns the size two tiers up using the provider's known ordered list of instance sizes; if there's no defined "+2" (already at the top), cap at the max available.  
   * Critical: this default only applies when the user hasn't set autoscaling explicitly. If the user provides their own autoscaling config, honor it and do not override (explicit input wins over the profile default).  
3. **Implementation approach:** put the conditional-default logic in the plan-modification / default-resolution path (wherever the provider currently computes defaults before apply), not as a hardcoded schema default — because the default is conditional on another field's value. Inspect how the resource currently handles computed/defaulted fields and follow that pattern.

4. **Keep it minimal and forkable:** the goal is a clean branch Alex can fork and tweak (different fields, different default values) easily. Favor clear, localized, well-commented changes over a polished/abstracted implementation. Don't refactor surrounding code. Add a short comment block where the profile logic lives explaining what it does, so someone forking it knows where to change values.

5. **Before writing code:** inspect the repo to confirm (a) which resource file and schema the cluster lives in, (b) how the provider represents autoscaling config (field names for compute enabled, min/max instance size), and (c) the ordered list of instance sizes used for tier progression. Report what you find, then implement.

6. **Tests/docs:** add a minimal acceptance-test or unit-test stub demonstrating CORE \= unchanged and INFINITE \= the autoscaling defaults applied (doesn't need to be exhaustive — just enough to show the behavior). Skip full docs; a comment in the schema attribute description is enough.  
7. **Add example Terraform configs** demonstrating the attribute, in the provider's examples directory (inspect where existing cluster examples live and follow that convention). Include:  
* An example with `cluster_profile = "CORE"` (or omitted, to show the default) — showing baseline behavior, no autoscaling defaults applied.  
* An example with `cluster_profile = "INFINITE"` — showing a minimal cluster config where the INFINITE autoscaling defaults (compute autoscaling on, min \= instance size, max \= instance size \+2) get applied automatically because the user hasn't set autoscaling explicitly.  
* Add a short comment in each example noting what behavior to expect (e.g. `# INFINITE: autoscaling defaults applied automatically — min=M30, max=M50`), so Alex can see at a glance what each profile does and easily copy/modify for his permutations.

**Out of scope (do not do):** reverse-compatibility handling, effective-fields/drift correctness, multi-shard behavior, migration logic. This is a single-branch prototype of the profile field \+ conditional autoscaling default only.