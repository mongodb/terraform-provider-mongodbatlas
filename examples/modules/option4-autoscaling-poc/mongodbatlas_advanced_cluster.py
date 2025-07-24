# codegen atlas-init-marker-start
import json
import sys
from dataclasses import asdict, dataclass
from typing import Optional, List, Dict, Any, Set, ClassVar, Union, Iterable
from typing import Tuple


@dataclass
class AdvancedConfiguration:
    NESTED_ATTRIBUTES: ClassVar[Set[str]] = {"custom_openssl_cipher_config_tls12"}
    REQUIRED_ATTRIBUTES: ClassVar[Set[str]] = set()
    COMPUTED_ONLY_ATTRIBUTES: ClassVar[Set[str]] = set()
    DEFAULTS_HCL_STRINGS: ClassVar[dict[str, str]] = {
        "default_write_concern": '"majority"',
        "javascript_enabled": "false",
        "minimum_enabled_tls_protocol": '"TLS1_2"',
    }
    change_stream_options_pre_and_post_images_expire_after_seconds: Optional[float] = None
    custom_openssl_cipher_config_tls12: Optional[List[str]] = None
    default_max_time_ms: Optional[float] = None
    default_write_concern: Optional[str] = None
    javascript_enabled: Optional[bool] = None
    minimum_enabled_tls_protocol: Optional[str] = None
    no_table_scan: Optional[bool] = None
    oplog_min_retention_hours: Optional[float] = None
    oplog_size_mb: Optional[float] = None
    sample_refresh_interval_bi_connector: Optional[float] = None
    sample_size_bi_connector: Optional[float] = None
    tls_cipher_config_mode: Optional[str] = None
    transaction_lifetime_limit_seconds: Optional[float] = None


@dataclass
class BiConnectorConfig:
    NESTED_ATTRIBUTES: ClassVar[Set[str]] = set()
    REQUIRED_ATTRIBUTES: ClassVar[Set[str]] = set()
    COMPUTED_ONLY_ATTRIBUTES: ClassVar[Set[str]] = set()
    DEFAULTS_HCL_STRINGS: ClassVar[dict[str, str]] = {}
    enabled: Optional[bool] = None
    read_preference: Optional[str] = None


@dataclass
class Endpoint:
    NESTED_ATTRIBUTES: ClassVar[Set[str]] = set()
    REQUIRED_ATTRIBUTES: ClassVar[Set[str]] = set()
    COMPUTED_ONLY_ATTRIBUTES: ClassVar[Set[str]] = {"endpoint_id", "provider_name", "region"}
    DEFAULTS_HCL_STRINGS: ClassVar[dict[str, str]] = {}
    endpoint_id: Optional[str] = None
    provider_name: Optional[str] = None
    region: Optional[str] = None


@dataclass
class PrivateEndpoint:
    NESTED_ATTRIBUTES: ClassVar[Set[str]] = {"endpoints"}
    REQUIRED_ATTRIBUTES: ClassVar[Set[str]] = set()
    COMPUTED_ONLY_ATTRIBUTES: ClassVar[Set[str]] = {
        "connection_string",
        "endpoints",
        "srv_connection_string",
        "srv_shard_optimized_connection_string",
        "type",
    }
    DEFAULTS_HCL_STRINGS: ClassVar[dict[str, str]] = {}
    connection_string: Optional[str] = None
    endpoints: Optional[List[Endpoint]] = None
    srv_connection_string: Optional[str] = None
    srv_shard_optimized_connection_string: Optional[str] = None
    type: Optional[str] = None

    def __post_init__(self):
        if self.endpoints is not None:
            self.endpoints = [x if isinstance(x, Endpoint) else Endpoint(**x) for x in self.endpoints]


@dataclass
class ConnectionString:
    NESTED_ATTRIBUTES: ClassVar[Set[str]] = {"private_endpoint"}
    REQUIRED_ATTRIBUTES: ClassVar[Set[str]] = set()
    COMPUTED_ONLY_ATTRIBUTES: ClassVar[Set[str]] = {
        "private",
        "private_endpoint",
        "private_srv",
        "standard",
        "standard_srv",
    }
    DEFAULTS_HCL_STRINGS: ClassVar[dict[str, str]] = {}
    private: Optional[str] = None
    private_endpoint: Optional[List[PrivateEndpoint]] = None
    private_srv: Optional[str] = None
    standard: Optional[str] = None
    standard_srv: Optional[str] = None

    def __post_init__(self):
        if self.private_endpoint is not None:
            self.private_endpoint = [
                x if isinstance(x, PrivateEndpoint) else PrivateEndpoint(**x) for x in self.private_endpoint
            ]


@dataclass
class PinnedFcv:
    NESTED_ATTRIBUTES: ClassVar[Set[str]] = set()
    REQUIRED_ATTRIBUTES: ClassVar[Set[str]] = {"expiration_date"}
    COMPUTED_ONLY_ATTRIBUTES: ClassVar[Set[str]] = {"version"}
    DEFAULTS_HCL_STRINGS: ClassVar[dict[str, str]] = {}
    expiration_date: Optional[str] = None
    version: Optional[str] = None


@dataclass
class Autoscaling:
    NESTED_ATTRIBUTES: ClassVar[Set[str]] = set()
    REQUIRED_ATTRIBUTES: ClassVar[Set[str]] = set()
    COMPUTED_ONLY_ATTRIBUTES: ClassVar[Set[str]] = set()
    DEFAULTS_HCL_STRINGS: ClassVar[dict[str, str]] = {}
    compute_enabled: Optional[bool] = None
    compute_max_instance_size: Optional[str] = None
    compute_min_instance_size: Optional[str] = None
    compute_scale_down_enabled: Optional[bool] = None
    disk_gb_enabled: Optional[bool] = None


@dataclass
class Spec:
    NESTED_ATTRIBUTES: ClassVar[Set[str]] = set()
    REQUIRED_ATTRIBUTES: ClassVar[Set[str]] = set()
    COMPUTED_ONLY_ATTRIBUTES: ClassVar[Set[str]] = set()
    DEFAULTS_HCL_STRINGS: ClassVar[dict[str, str]] = {}
    disk_iops: Optional[float] = None
    disk_size_gb: Optional[float] = None
    ebs_volume_type: Optional[str] = None
    instance_size: Optional[str] = None
    node_count: Optional[float] = None


@dataclass
class RegionConfig:
    NESTED_ATTRIBUTES: ClassVar[Set[str]] = {
        "analytics_auto_scaling",
        "analytics_specs",
        "auto_scaling",
        "electable_specs",
        "read_only_specs",
    }
    REQUIRED_ATTRIBUTES: ClassVar[Set[str]] = {"priority", "provider_name", "region_name"}
    COMPUTED_ONLY_ATTRIBUTES: ClassVar[Set[str]] = set()
    DEFAULTS_HCL_STRINGS: ClassVar[dict[str, str]] = {}
    analytics_auto_scaling: Optional[Autoscaling] = None
    analytics_specs: Optional[Spec] = None
    auto_scaling: Optional[Autoscaling] = None
    backing_provider_name: Optional[str] = None
    electable_specs: Optional[Spec] = None
    priority: Optional[float] = None
    provider_name: Optional[str] = None
    read_only_specs: Optional[Spec] = None
    region_name: Optional[str] = None

    def __post_init__(self):
        if self.analytics_auto_scaling is not None and not isinstance(self.analytics_auto_scaling, Autoscaling):
            assert isinstance(self.analytics_auto_scaling, dict), (
                f"Expected analytics_auto_scaling to be a Autoscaling or a dict, got {type(self.analytics_auto_scaling)}"
            )
            self.analytics_auto_scaling = Autoscaling(**self.analytics_auto_scaling)
        if self.analytics_specs is not None and not isinstance(self.analytics_specs, Spec):
            assert isinstance(self.analytics_specs, dict), (
                f"Expected analytics_specs to be a Spec or a dict, got {type(self.analytics_specs)}"
            )
            self.analytics_specs = Spec(**self.analytics_specs)
        if self.auto_scaling is not None and not isinstance(self.auto_scaling, Autoscaling):
            assert isinstance(self.auto_scaling, dict), (
                f"Expected auto_scaling to be a Autoscaling or a dict, got {type(self.auto_scaling)}"
            )
            self.auto_scaling = Autoscaling(**self.auto_scaling)
        if self.electable_specs is not None and not isinstance(self.electable_specs, Spec):
            assert isinstance(self.electable_specs, dict), (
                f"Expected electable_specs to be a Spec or a dict, got {type(self.electable_specs)}"
            )
            self.electable_specs = Spec(**self.electable_specs)
        if self.read_only_specs is not None and not isinstance(self.read_only_specs, Spec):
            assert isinstance(self.read_only_specs, dict), (
                f"Expected read_only_specs to be a Spec or a dict, got {type(self.read_only_specs)}"
            )
            self.read_only_specs = Spec(**self.read_only_specs)


@dataclass
class ReplicationSpec:
    NESTED_ATTRIBUTES: ClassVar[Set[str]] = {"container_id", "region_configs"}
    REQUIRED_ATTRIBUTES: ClassVar[Set[str]] = {"region_configs"}
    COMPUTED_ONLY_ATTRIBUTES: ClassVar[Set[str]] = {"container_id", "external_id", "zone_id"}
    DEFAULTS_HCL_STRINGS: ClassVar[dict[str, str]] = {}
    container_id: Optional[Dict[str, Any]] = None
    external_id: Optional[str] = None
    region_configs: Optional[List[RegionConfig]] = None
    zone_id: Optional[str] = None
    zone_name: Optional[str] = None

    def __post_init__(self):
        if self.region_configs is not None:
            self.region_configs = [x if isinstance(x, RegionConfig) else RegionConfig(**x) for x in self.region_configs]


@dataclass
class Timeout:
    NESTED_ATTRIBUTES: ClassVar[Set[str]] = set()
    REQUIRED_ATTRIBUTES: ClassVar[Set[str]] = set()
    COMPUTED_ONLY_ATTRIBUTES: ClassVar[Set[str]] = set()
    DEFAULTS_HCL_STRINGS: ClassVar[dict[str, str]] = {}
    create: Optional[str] = None
    delete: Optional[str] = None
    update: Optional[str] = None


@dataclass
class Resource:
    NESTED_ATTRIBUTES: ClassVar[Set[str]] = {
        "advanced_configuration",
        "bi_connector_config",
        "connection_strings",
        "labels",
        "pinned_fcv",
        "replication_specs",
        "tags",
        "timeouts",
    }
    REQUIRED_ATTRIBUTES: ClassVar[Set[str]] = {"cluster_type", "name", "project_id", "replication_specs"}
    COMPUTED_ONLY_ATTRIBUTES: ClassVar[Set[str]] = {
        "cluster_id",
        "config_server_type",
        "connection_strings",
        "create_date",
        "mongo_db_version",
        "state_name",
    }
    DEFAULTS_HCL_STRINGS: ClassVar[dict[str, str]] = {
        "backup_enabled": "true",
        "retain_backups_enabled": "true",
        "termination_protection_enabled": "true",
    }
    accept_data_risks_and_force_replica_set_reconfig: Optional[str] = None
    advanced_configuration: Optional[AdvancedConfiguration] = None
    backup_enabled: Optional[bool] = None
    bi_connector_config: Optional[BiConnectorConfig] = None
    cluster_id: Optional[str] = None
    cluster_type: Optional[str] = None
    config_server_management_mode: Optional[str] = None
    config_server_type: Optional[str] = None
    connection_strings: Optional[ConnectionString] = None
    create_date: Optional[str] = None
    delete_on_create_timeout: Optional[bool] = None
    encryption_at_rest_provider: Optional[str] = None
    global_cluster_self_managed_sharding: Optional[bool] = None
    labels: Optional[Dict[str, Any]] = None
    mongo_db_major_version: Optional[str] = None
    mongo_db_version: Optional[str] = None
    name: Optional[str] = None
    paused: Optional[bool] = None
    pinned_fcv: Optional[PinnedFcv] = None
    pit_enabled: Optional[bool] = None
    project_id: Optional[str] = None
    redact_client_log_data: Optional[bool] = None
    replica_set_scaling_strategy: Optional[str] = None
    replication_specs: Optional[List[ReplicationSpec]] = None
    retain_backups_enabled: Optional[bool] = None
    root_cert_type: Optional[str] = None
    state_name: Optional[str] = None
    tags: Optional[Dict[str, Any]] = None
    termination_protection_enabled: Optional[bool] = None
    timeouts: Optional[Timeout] = None
    version_release_system: Optional[str] = None

    def __post_init__(self):
        if self.advanced_configuration is not None and not isinstance(
            self.advanced_configuration, AdvancedConfiguration
        ):
            assert isinstance(self.advanced_configuration, dict), (
                f"Expected advanced_configuration to be a AdvancedConfiguration or a dict, got {type(self.advanced_configuration)}"
            )
            self.advanced_configuration = AdvancedConfiguration(**self.advanced_configuration)
        if self.bi_connector_config is not None and not isinstance(self.bi_connector_config, BiConnectorConfig):
            assert isinstance(self.bi_connector_config, dict), (
                f"Expected bi_connector_config to be a BiConnectorConfig or a dict, got {type(self.bi_connector_config)}"
            )
            self.bi_connector_config = BiConnectorConfig(**self.bi_connector_config)
        if self.connection_strings is not None and not isinstance(self.connection_strings, ConnectionString):
            assert isinstance(self.connection_strings, dict), (
                f"Expected connection_strings to be a ConnectionString or a dict, got {type(self.connection_strings)}"
            )
            self.connection_strings = ConnectionString(**self.connection_strings)
        if self.pinned_fcv is not None and not isinstance(self.pinned_fcv, PinnedFcv):
            assert isinstance(self.pinned_fcv, dict), (
                f"Expected pinned_fcv to be a PinnedFcv or a dict, got {type(self.pinned_fcv)}"
            )
            self.pinned_fcv = PinnedFcv(**self.pinned_fcv)
        if self.replication_specs is not None:
            self.replication_specs = [
                x if isinstance(x, ReplicationSpec) else ReplicationSpec(**x) for x in self.replication_specs
            ]
        if self.timeouts is not None and not isinstance(self.timeouts, Timeout):
            assert isinstance(self.timeouts, dict), (
                f"Expected timeouts to be a Timeout or a dict, got {type(self.timeouts)}"
            )
            self.timeouts = Timeout(**self.timeouts)
        if self.auto_scaling is not None and not isinstance(self.auto_scaling, Autoscaling):
            assert isinstance(self.auto_scaling, dict), (
                f"Expected auto_scaling to be a Autoscaling or a dict, got {type(self.auto_scaling)}"
            )
            self.auto_scaling = Autoscaling(**self.auto_scaling)
        if self.auto_scaling_analytics is not None and not isinstance(self.auto_scaling_analytics, Autoscaling):
            assert isinstance(self.auto_scaling_analytics, dict), (
                f"Expected auto_scaling_analytics to be a Autoscaling or a dict, got {type(self.auto_scaling_analytics)}"
            )
            self.auto_scaling_analytics = Autoscaling(**self.auto_scaling_analytics)

        if self.regions is not None:
            self.regions = [x if isinstance(x, Region) else Region(**x) for x in self.regions]


def format_primitive(value: Union[str, float, bool, int, None]):
    if value is None:
        return None
    if value is True:
        return "true"
    if value is False:
        return "false"
    return str(value)


def main():
    input_data = sys.stdin.read()
    # Parse the input as JSON
    params = json.loads(input_data)
    input_json = params["input_json"]
    resource = ResourceExt(**json.loads(input_json))
    error_message = "\n".join(errors(resource))
    primitive_types = (str, float, bool, int)
    resource = modify_out(resource)
    output = {
        key: format_primitive(value) if value is None or isinstance(value, primitive_types) else json.dumps(value)
        for key, value in asdict(resource).items()
    }
    output["error_message"] = error_message
    json_str = json.dumps(output)
    print(json_str)


# codegen atlas-init-marker-end


@dataclass
class Region:
    name: str
    node_count: int
    shard_index: Optional[int] = None
    provider_name: Optional[str] = None
    node_count_read_only: Optional[int] = None
    node_count_analytics: Optional[int] = None
    instance_size: Optional[str] = None
    instance_size_analytics: Optional[str] = None
    zone_name: Optional[str] = None


@dataclass
class ResourceExt(Resource):
    regions: Optional[List[Region]] = None
    provider_name: Optional[str] = None
    instance_size: Optional[str] = None
    disk_size_gb: Optional[float] = None
    instance_size_analytics: Optional[str] = None
    auto_scaling: Optional[Autoscaling] = None
    auto_scaling_analytics: Optional[Autoscaling] = None
    old_cluster: Optional[Resource] = None

    DEFAULT_INSTANCE_SIZE: ClassVar[str] = "M10"
    MUTUALLY_EXCLUSIVE: ClassVar[dict[str, list[str]]] = {
        "regions": ["replication_specs"],
        "auto_scaling": ["instance_size"],
        "auto_scaling_analytics": ["instance_size_analytics"],
    }
    REQUIRES_OTHER: ClassVar[dict[str, list[str]]] = {
        "instance_size": ["regions"],
        "num_shards": ["regions"],
        "auto_scaling": ["regions"],
        "auto_scaling_analytics": ["regions"],
    }
    SKIP_VARIABLES: ClassVar[set[str]] = {"old_cluster"}

    def __post_init__(self):
        super().__post_init__()
        if self.old_cluster is not None and not isinstance(self.old_cluster, Resource):
            assert isinstance(self.old_cluster, dict), (
                f"Expected old_cluster to be a Resource or a dict, got {type(self.old_cluster)}"
            )
            self.old_cluster = Resource(
                **{k: v for k, v in self.old_cluster.items() if k not in {"use_replication_spec_per_shard"}}
            )
        if self.auto_scaling is not None and not isinstance(self.auto_scaling, Autoscaling):
            assert isinstance(self.auto_scaling, dict), (
                f"Expected auto_scaling to be an Autoscaling or a dict, got {type(self.auto_scaling)}"
            )
            self.auto_scaling = Autoscaling(**self.auto_scaling)
        if self.auto_scaling_analytics is not None and not isinstance(self.auto_scaling_analytics, Autoscaling):
            assert isinstance(self.auto_scaling_analytics, dict), (
                f"Expected auto_scaling_analytics to be an Autoscaling or a dict, got {type(self.auto_scaling_analytics)}"
            )
            self.auto_scaling_analytics = Autoscaling(**self.auto_scaling_analytics)
        if self.regions is not None:
            self.regions = [Region(**region) if isinstance(region, dict) else region for region in self.regions]

    @property
    def num_shards(self) -> int:
        if self.infer_cluster_type() == "GEOSHARDED":
            return len({region.zone_name for region in self.regions or [] if region.zone_name is not None})
        return (
            max((region.shard_index or 0 for region in self.regions or [] if region.shard_index is not None), default=0)
            + 1
        )

    def infer_cluster_type(self) -> str:
        if self.cluster_type:
            return self.cluster_type
        if all(region.zone_name is not None for region in self.regions or []):
            return "GEOSHARDED"
        if all(region.shard_index is not None for region in self.regions or []):
            return "SHARDED"
        return "REPLICASET"

    def iterate_rep_spec_region_configs(self) -> Iterable[Tuple[int, list[Region]]]:
        regions = self.regions or []
        if self.infer_cluster_type() == "REPLICASET":
            yield 0, regions
        elif self.infer_cluster_type() == "SHARDED":
            num_shards = self.num_shards or 1
            shard_regions = {index: [] for index in range(num_shards)}
            for region in regions:
                shard_regions[region.shard_index or 0].append(region)
            yield from shard_regions.items()
        else:
            zone_shard_indexes = {}
            zone_regions = {index: [] for index in range(self.num_shards or 1)}
            current_shard_index = 0
            for region in regions:
                assert region.zone_name
                if region.zone_name not in zone_shard_indexes:
                    zone_shard_indexes[region.zone_name] = current_shard_index
                    current_shard_index += 1
                zone_regions[zone_shard_indexes[region.zone_name]].append(region)
            yield from zone_regions.items()

    def get_instance_size_electable(self, region: Region, shard_index: int, region_config_index: int) -> str:
        if self.auto_scaling is None:
            return region.instance_size or self.instance_size or self.DEFAULT_INSTANCE_SIZE
        default_min_size = self.auto_scaling.compute_min_instance_size
        assert default_min_size is not None, (
            f"{self.auto_scaling.compute_min_instance_size} is not a valid instance size"
        )
        if self.old_cluster is not None:
            return self.current_instance_size_electable(shard_index, region_config_index) or default_min_size
        else:
            return default_min_size

    def current_instance_size_electable(self, shard_index: int, region_config_index: int) -> Optional[str]:
        old_cluster = self.old_cluster
        assert old_cluster is not None
        specs = old_cluster.replication_specs
        if not specs or len(specs) <= shard_index:
            return ""
        shard = specs[shard_index]
        region_configs = shard.region_configs
        if not region_configs or len(region_configs) <= region_config_index:
            return ""
        region_config = region_configs[region_config_index]
        if region_config.electable_specs is None:
            return ""
        return region_config.electable_specs.instance_size

    def get_instance_size_analytics(self, region: Region, shard_index: int, region_config_index: int) -> str:
        if self.auto_scaling_analytics is None:
            return region.instance_size_analytics or self.instance_size_analytics or self.DEFAULT_INSTANCE_SIZE
        default_min_size = self.auto_scaling_analytics.compute_min_instance_size
        assert default_min_size is not None, (
            f"{self.auto_scaling_analytics.compute_min_instance_size} is not a valid instance size"
        )
        if self.old_cluster is not None:
            return self.current_instance_size_analytics(shard_index, region_config_index) or default_min_size
        else:
            return default_min_size

    def current_instance_size_analytics(self, shard_index: int, region_config_index: int) -> Optional[str]:
        old_cluster = self.old_cluster
        assert old_cluster is not None
        specs = old_cluster.replication_specs
        if not specs or len(specs) <= shard_index:
            return ""
        shard = specs[shard_index]
        region_configs = shard.region_configs
        if not region_configs or len(region_configs) <= region_config_index:
            return ""
        region_config = region_configs[region_config_index]
        if region_config.analytics_specs is None:
            return ""
        return region_config.analytics_specs.instance_size


def errors(resource: ResourceExt) -> Iterable[str]:
    for var, mutually_exclusive_fields in resource.MUTUALLY_EXCLUSIVE.items():
        if not getattr(resource, var):
            continue
        for incompatible in mutually_exclusive_fields:
            if getattr(resource, incompatible):
                yield f"Cannot use var.{var} and var.{incompatible} together"
    for var, required_vars in resource.REQUIRES_OTHER.items():
        if not getattr(resource, var):
            continue
        missing_required = [required for required in required_vars if not getattr(resource, required)]
        if missing_required:
            yield f"Cannot use {var} without {','.join(sorted(missing_required))}"
    if resource.auto_scaling is not None:
        invalid_instance_sizes = [
            f"instance_size @ index {index} = {region.instance_size}"
            for index, region in enumerate(resource.regions or [])
            if region.instance_size is not None
        ]
        if invalid_instance_sizes:
            yield f"Cannot use `regions.*.instance_size` when auto_scaling is used: {','.join(invalid_instance_sizes)}"
    if resource.auto_scaling_analytics is not None:
        invalid_instance_sizes = [
            f"instance_size @ index {index} = {region.instance_size_analytics}"
            for index, region in enumerate(resource.regions or [])
            if region.instance_size_analytics is not None
        ]
        if invalid_instance_sizes:
            yield f"Cannot use `regions.*.instance_size_analytics` when auto_scaling_analytics is used: {','.join(invalid_instance_sizes)}"
    found_cluster_type = resource.infer_cluster_type()
    zone_names_found = [
        f"regions[{i}].zone_name={region.zone_name}"
        for i, region in enumerate(resource.regions or [])
        if region.zone_name is not None
    ]
    shard_indexes_found = [
        f"regions[{i}].shard_index={region.shard_index}"
        for i, region in enumerate(resource.regions or [])
        if region.shard_index is not None
    ]
    if found_cluster_type == "GEOSHARDED":
        missing_zone_names = [
            f"zone_name missing @ index {index}"
            for index, region in enumerate(resource.regions or [])
            if region.zone_name is None
        ]
        if missing_zone_names:
            yield f"Must use `regions.*.zone_name` when cluster_type is GEOSHARDED: {','.join(missing_zone_names)}"
        if shard_indexes_found:
            yield f"Geosharded cluster should not define shard_index: {','.join(shard_indexes_found)}"
    if found_cluster_type == "SHARDED":
        missing_shard_indexes = [
            f"shard_index missing @ index {index}"
            for index, region in enumerate(resource.regions or [])
            if region.shard_index is None
        ]
        if missing_shard_indexes:
            yield f"Must use `regions.*.shard_index` when cluster_type is SHARDED: {','.join(missing_shard_indexes)}"
        if zone_names_found:
            yield f"Sharded cluster should not define zone_name: {','.join(zone_names_found)}"
    if found_cluster_type == "REPLICASET":
        if shard_indexes_found:
            yield f"Replicaset cluster should not define shard_index: {','.join(shard_indexes_found)}"
        if zone_names_found:
            yield f"Replicaset cluster should not define zone_name: {','.join(zone_names_found)}"
    if missing_cloud_provider := [
        f"regions[{i}].provider_name is missing"
        for i, region in enumerate(resource.regions or [])
        if region.provider_name is None
    ]:
        if resource.provider_name is None:
            yield f"Must use `regions.*.provider_name` when root `provider_name` is not specified: {','.join(missing_cloud_provider)}"


def generate_replication_specs(resource: ResourceExt) -> list[ReplicationSpec]:
    specs = []
    auto_scaling = resource.auto_scaling
    auto_scaling_analytics = resource.auto_scaling_analytics
    for rep_spec_index, regions in resource.iterate_rep_spec_region_configs():
        spec = ReplicationSpec(
            region_configs=[],
            zone_name=regions[0].zone_name,
        )
        specs.append(spec)
        current_priority = 7
        for region_config_index, region in enumerate(regions):
            electable = (
                Spec(
                    disk_size_gb=resource.disk_size_gb,
                    instance_size=resource.get_instance_size_electable(
                        region, rep_spec_index, region_config_index=region_config_index
                    ),
                    node_count=region.node_count,
                )
                if region.node_count
                else None
            )
            analytics = (
                Spec(
                    disk_size_gb=resource.disk_size_gb,
                    instance_size=resource.get_instance_size_analytics(
                        region, rep_spec_index, region_config_index=region_config_index
                    ),
                    node_count=region.node_count_analytics,
                )
                if region.node_count_analytics
                else None
            )
            read_only = (
                Spec(
                    disk_size_gb=resource.disk_size_gb,
                    instance_size=resource.get_instance_size_electable(
                        region, rep_spec_index, region_config_index=region_config_index
                    ),
                    node_count=region.node_count_read_only,
                )
                if region.node_count_read_only
                else None
            )
            region_config = RegionConfig(
                provider_name=region.provider_name or resource.provider_name,
                region_name=region.name,
                priority=current_priority,
                electable_specs=electable,
                read_only_specs=read_only,
                analytics_specs=analytics,
                auto_scaling=auto_scaling,
                analytics_auto_scaling=auto_scaling_analytics,
            )
            current_priority -= 1
            assert spec.region_configs is not None
            spec.region_configs.append(region_config)
    return specs


def modify_out(resource: ResourceExt) -> ResourceExt:
    if resource.regions:
        resource.replication_specs = generate_replication_specs(resource)
        resource.cluster_type = resource.infer_cluster_type()
    return resource


if __name__ == "__main__":
    main()
