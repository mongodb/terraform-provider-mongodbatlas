# cluster_profile PROTOTYPE — INFINITE

Demonstrates the `cluster_profile` prototype attribute set to `INFINITE`.

The cluster configures a single `M30` region with **no** `auto_scaling` block. Because
`cluster_profile = "INFINITE"`, the provider automatically injects compute auto-scaling
defaults during plan:

| field                        | value  | meaning                          |
| ---------------------------- | ------ | -------------------------------- |
| `compute_enabled`            | `true` | scale up enabled                 |
| `compute_scale_down_enabled` | `true` | scale down enabled               |
| `compute_min_instance_size`  | `M30`  | = configured instance size       |
| `compute_max_instance_size`  | `M50`  | = two tiers above (M30→M40→M50)   |
| `disk_gb_enabled`            | `false`| only compute auto-scaling is set |

Run `terraform plan` to see these appear on `auto_scaling` automatically.

**Explicit input always wins:** if you add your own `auto_scaling` block to a region
config, the profile leaves it untouched.

Compare with the [`cluster-profile-core`](../cluster-profile-core) example for the baseline.

> Prototype/demo only — not for production use.
