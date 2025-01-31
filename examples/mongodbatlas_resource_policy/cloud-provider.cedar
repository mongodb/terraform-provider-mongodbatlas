forbid (
    principal,
    action == ResourcePolicy::Action::"cluster.modify",
    resource
)
when
{ context.cluster.cloudProviders.containsAny([cloud::cloudProvider::"${CLOUD_PROVIDER}"]) };