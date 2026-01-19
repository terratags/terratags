package parser

// AliCloud taggable resources
// This list includes common AliCloud resources that support the 'tags' attribute
var alicloudTaggableResources = map[string]bool{
	// ECS (Elastic Compute Service)
	"alicloud_instance":                    true,
	"alicloud_reserved_instance":           true,
	"alicloud_ecs_instance_set":           true,
	"alicloud_simple_application_server_instance": true,

	// Storage
	"alicloud_oss_bucket":                 true,
	"alicloud_oss_bucket_object":          true,

	// Database
	"alicloud_db_instance":                true,
	"alicloud_db_readonly_instance":       true,
	"alicloud_rds_clone_db_instance":      true,
	"alicloud_rds_upgrade_db_instance":    true,
	"alicloud_rds_ddr_instance":           true,
	"alicloud_rds_ai_instance":            true,
	"alicloud_mongodb_instance":           true,
	"alicloud_mongodb_sharding_instance":  true,
	"alicloud_mongodb_serverless_instance": true,
	"alicloud_redis_tair_instance":        true,
	"alicloud_kvstore_instance":           true,
	"alicloud_gpdb_instance":              true,
	"alicloud_gpdb_elastic_instance":      true,
	"alicloud_hbase_instance":             true,
	"alicloud_lindorm_instance":           true,
	"alicloud_lindorm_instance_v2":        true,
	"alicloud_ocean_base_instance":        true,
	"alicloud_selectdb_db_instance":       true,
	"alicloud_star_rocks_instance":        true,
	"alicloud_tsdb_instance":              true,
	"alicloud_graph_database_db_instance": true,
	"alicloud_hologram_instance":          true,
	"alicloud_milvus_instance":            true,

	// Networking
	"alicloud_vpc":                        true,
	"alicloud_vswitch":                    true,
	"alicloud_security_group":             true,
	"alicloud_nat_gateway":                true,
	"alicloud_eip":                        true,
	"alicloud_slb":                        true,
	"alicloud_slb_load_balancer":          true,
	"alicloud_alb_load_balancer":          true,
	"alicloud_nlb_load_balancer":          true,

	// Container Services
	"alicloud_cs_kubernetes_cluster":      true,
	"alicloud_cs_managed_kubernetes":      true,
	"alicloud_cs_serverless_kubernetes":   true,
	"alicloud_cs_edge_kubernetes":         true,

	// Message Queue
	"alicloud_ons_instance":               true,
	"alicloud_alikafka_instance":          true,
	"alicloud_rocketmq_instance":          true,
	"alicloud_amqp_instance":              true,

	// CDN & DNS
	"alicloud_cdn_domain":                 true,
	"alicloud_dns_instance":               true,
	"alicloud_alidns_instance":            true,
	"alicloud_alidns_gtm_instance":        true,

	// Security
	"alicloud_kms_key":                    true,
	"alicloud_kms_instance":               true,
	"alicloud_bastionhost_instance":       true,
	"alicloud_waf_instance":               true,
	"alicloud_wafv3_instance":             true,
	"alicloud_cloud_firewall_instance":    true,
	"alicloud_threat_detection_instance":  true,
	"alicloud_sddp_instance":              true,
	"alicloud_yundun_dbaudit_instance":    true,

	// Analytics & Big Data
	"alicloud_elasticsearch_instance":     true,
	"alicloud_realtime_compute_vvp_instance": true,

	// API Gateway
	"alicloud_api_gateway_instance":       true,

	// Others
	"alicloud_ots_instance":               true,
	"alicloud_cr_ee_instance":             true,
	"alicloud_dts_instance":               true,
	"alicloud_dts_migration_instance":     true,
	"alicloud_dts_synchronization_instance": true,
	"alicloud_drds_instance":              true,
	"alicloud_drds_polardbx_instance":     true,
	"alicloud_eais_instance":              true,
	"alicloud_ebs_solution_instance":      true,
	"alicloud_ecp_instance":               true,
	"alicloud_ens_instance":               true,
	"alicloud_esa_cache_reserve_instance": true,
	"alicloud_esa_rate_plan_instance":     true,
	"alicloud_compute_nest_service_instance": true,
	"alicloud_dbfs_instance":              true,
	"alicloud_ddosbgp_instance":           true,
	"alicloud_ddoscoo_instance":           true,
	"alicloud_dms_enterprise_instance":    true,
	"alicloud_hbr_hana_instance":          true,
	"alicloud_cloud_phone_instance":       true,
	"alicloud_cloud_phone_instance_group": true,
}
