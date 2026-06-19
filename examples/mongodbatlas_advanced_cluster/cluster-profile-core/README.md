# cluster_profile PROTOTYPE — CORE (baseline)

Demonstrates the `cluster_profile` prototype attribute set to `CORE` (the baseline).

Expected behavior: **no change from today.** No `auto_scaling` block is configured
and CORE applies no defaults, so the cluster stays a fixed `M30`.

Setting `cluster_profile = "CORE"` is equivalent to omitting the attribute entirely.

Compare with the [`cluster-profile-infinite`](../cluster-profile-infinite) example to see
how `INFINITE` injects auto-scaling defaults automatically.

> Prototype/demo only — not for production use.
