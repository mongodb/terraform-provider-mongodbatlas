forbid (
    principal,
    action == cloud::Action::"cluster.createEdit",
    resource
)
when
{ context.cluster.cloudProviders.containsAny([cloud::cloudProvider::"${CLOUD_PROVIDER}"]) };