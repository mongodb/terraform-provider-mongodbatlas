# MINIMAL-CONFIG PROTOTYPE — FULL EXPLICIT (reverse-compatibility)

An existing-style config that sets every (now-optional) input explicitly — `project_id`,
`cluster_type`, and a full `replication_specs` — and uses neither `cluster_profile` nor
`provider_region`.

Its purpose is to confirm that making those inputs `Required → Optional` did **not** change
behavior for users who specify everything: `terraform plan` shows exactly the configured values
with no profile-driven defaulting (no synthesized `replication_specs`, no injected
`auto_scaling`). See the parent task notes for the captured plan output.

> Prototype/demo only. Replace the `project_id` placeholder with a real project id to apply.
