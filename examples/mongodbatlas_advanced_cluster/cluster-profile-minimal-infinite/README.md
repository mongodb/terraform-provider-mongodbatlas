# MINIMAL-CONFIG PROTOTYPE — INFINITE

Deploy a full cluster from just four inputs: `project_id`, `name`, `provider_region`, and
`cluster_profile = "INFINITE"`. The remaining previously-required inputs are filled in by the
profile during `terraform plan`:

| input               | filled value                                  | source                        |
| ------------------- | --------------------------------------------- | ----------------------------- |
| `cluster_type`      | `REPLICASET`                                  | static default                |
| `replication_specs` | 1 × `AWS:US_EAST_1` shard, `M30`, 3 nodes     | `provider_region` + profile   |
| `auto_scaling`      | compute on (up+down), min `M30`, max `M50`    | INFINITE "+2" rule (branch 1) |

`project_id` stays a required input (no sensible default). Replace the placeholder in `main.tf`
with a real project id.

> Prototype/demo only.
