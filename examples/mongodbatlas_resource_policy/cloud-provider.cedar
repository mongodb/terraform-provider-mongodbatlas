forbid (
    principal,
    action == ResourcePolicy::Action::"cluster.modify",
    resource
)
when
{ context.cluster.cloudProviders.containsAny([ResourcePolicy::CloudProvider::"${CLOUD_PROVIDER}"]) };