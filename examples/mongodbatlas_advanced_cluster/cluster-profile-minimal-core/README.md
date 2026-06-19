# MINIMAL-CONFIG PROTOTYPE — CORE

Deploy a full cluster from just three inputs: `project_id`, `name`, and `provider_region`.
`cluster_profile` is omitted here, which behaves as CORE (an unset profile is treated the same
as `"CORE"`). The remaining previously-required inputs are filled in by the profile during
`terraform plan`:

| input               | filled value                              | source                      |
| ------------------- | ----------------------------------------- | --------------------------- |
| `cluster_type`      | `REPLICASET`                              | static default              |
| `replication_specs` | 1 × `AWS:US_EAST_1` shard, `M10`, 3 nodes | `provider_region` + profile |
| `auto_scaling`      | _(known after apply)_ — none injected     | CORE applies no defaults    |

`project_id` stays a required input (no sensible default). Replace the placeholder in `main.tf`
with a real project id.

Compare with [`cluster-profile-minimal-infinite`](../cluster-profile-minimal-infinite) to see
how INFINITE additionally injects auto-scaling.

> Prototype/demo only.
