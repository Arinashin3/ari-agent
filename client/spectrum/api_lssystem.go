package spectrum

import "encoding/json"

type LsSystemInst struct {
	Id                            string `json:"id,omitempty"`
	Name                          string `json:"name,omitempty"`
	Location                      string `json:"location,omitempty"`
	Partnership                   string `json:"partnership,omitempty"`
	TotalMdiskCapacity            string `json:"total_mdisk_capacity,omitempty"`
	SpaceInMdiskGrps              string `json:"space_in_mdisk_grps,omitempty"`
	SpaceAllocatedToVdisks        string `json:"space_allocated_to_vdisks,omitempty"`
	TotalFreeSpace                string `json:"total_free_space,omitempty"`
	TotalVdiskcopyCapacity        string `json:"total_vdiskcopy_capacity,omitempty"`
	TotalUsedCapacity             string `json:"total_used_capacity,omitempty"`
	TotalOverallocation           string `json:"total_overallocation,omitempty"`
	TotalVdiskCapacity            string `json:"total_vdisk_capacity,omitempty"`
	TotalAllocatedExtentCapacity  string `json:"total_allocated_extent_capacity,omitempty"`
	StatisticsStatus              string `json:"statistics_status,omitempty"`
	StatisticsFrequency           string `json:"statistics_frequency,omitempty"`
	ClusterLocale                 string `json:"cluster_locale,omitempty"`
	TimeZone                      string `json:"time_zone,omitempty"`
	CodeLevel                     string `json:"code_level,omitempty"`
	ConsoleIP                     string `json:"console_IP,omitempty"`
	IdAlias                       string `json:"id_alias,omitempty"`
	GmLinkTolerance               string `json:"gm_link_tolerance,omitempty"`
	GmInterClusterDelaySimulation string `json:"gm_inter_cluster_delay_simulation,omitempty"`
	GmIntraClusterDelaySimulation string `json:"gm_intra_cluster_delay_simulation,omitempty"`
	GmMaxHostDelay                string `json:"gm_max_host_delay,omitempty"`
	EmailReply                    string `json:"email_reply,omitempty"`
	EmailContact                  string `json:"email_contact,omitempty"`
	EmailContactPrimary           string `json:"email_contact_primary,omitempty"`
	EmailContactAlternate         string `json:"email_contact_alternate,omitempty"`
	EmailContactLocation          string `json:"email_contact_location,omitempty"`
	EmailContact2                 string `json:"email_contact2,omitempty"`
	EmailContact2Primary          string `json:"email_contact2_primary,omitempty"`
	EmailContact2Alternate        string `json:"email_contact2_alternate,omitempty"`
	EmailState                    string `json:"email_state,omitempty"`
	InventoryMailInterval         string `json:"inventory_mail_interval,omitempty"`
	ClusterNtpIPAddress           string `json:"cluster_ntp_IP_address,omitempty"`
	ClusterIsnsIPAddress          string `json:"cluster_isns_IP_address,omitempty"`
	IscsiAuthMethod               string `json:"iscsi_auth_method,omitempty"`
	IscsiChapSecret               string `json:"iscsi_chap_secret,omitempty"`
	AuthServiceConfigured         string `json:"auth_service_configured,omitempty"`
	AuthServiceEnabled            string `json:"auth_service_enabled,omitempty"`
	AuthServiceUrl                string `json:"auth_service_url,omitempty"`
	AuthServiceUserName           string `json:"auth_service_user_name,omitempty"`
	AuthServicePwdSet             string `json:"auth_service_pwd_set,omitempty"`
	AuthServiceCertSet            string `json:"auth_service_cert_set,omitempty"`
	AuthServiceType               string `json:"auth_service_type,omitempty"`
	RelationshipBandwidthLimit    string `json:"relationship_bandwidth_limit,omitempty"`
	Tiers                         []struct {
		Tier             string `json:"tier,omitempty"`
		TierCapacity     string `json:"tier_capacity,omitempty"`
		TierFreeCapacity string `json:"tier_free_capacity,omitempty"`
	} `json:"tiers,omitempty"`
	EasyTierAcceleration             string `json:"easy_tier_acceleration,omitempty"`
	HasNasKey                        string `json:"has_nas_key,omitempty"`
	Layer                            string `json:"layer,omitempty"`
	RcBufferSize                     string `json:"rc_buffer_size,omitempty"`
	CompressionActive                string `json:"compression_active,omitempty"`
	CompressionVirtualCapacity       string `json:"compression_virtual_capacity,omitempty"`
	CompressionCompressedCapacity    string `json:"compression_compressed_capacity,omitempty"`
	CompressionUncompressedCapacity  string `json:"compression_uncompressed_capacity,omitempty"`
	CachePrefetch                    string `json:"cache_prefetch,omitempty"`
	EmailOrganization                string `json:"email_organization,omitempty"`
	EmailMachineAddress              string `json:"email_machine_address,omitempty"`
	EmailMachineCity                 string `json:"email_machine_city,omitempty"`
	EmailMachineState                string `json:"email_machine_state,omitempty"`
	EmailMachineZip                  string `json:"email_machine_zip,omitempty"`
	EmailMachineCountry              string `json:"email_machine_country,omitempty"`
	TotalDriveRawCapacity            string `json:"total_drive_raw_capacity,omitempty"`
	CompressionDestageMode           string `json:"compression_destage_mode,omitempty"`
	LocalFcPortMask                  string `json:"local_fc_port_mask,omitempty"`
	PartnerFcPortMask                string `json:"partner_fc_port_mask,omitempty"`
	HighTempMode                     string `json:"high_temp_mode,omitempty"`
	Topology                         string `json:"topology,omitempty"`
	TopologyStatus                   string `json:"topology_status,omitempty"`
	RcAuthMethod                     string `json:"rc_auth_method,omitempty"`
	VdiskProtectionTime              string `json:"vdisk_protection_time,omitempty"`
	VdiskProtectionEnabled           string `json:"vdisk_protection_enabled,omitempty"`
	ProductName                      string `json:"product_name,omitempty"`
	Odx                              string `json:"odx,omitempty"`
	MaxReplicationDelay              string `json:"max_replication_delay,omitempty"`
	PartnershipExclusionThreshold    string `json:"partnership_exclusion_threshold,omitempty"`
	Gen1CompatibilityModeEnabled     string `json:"gen1_compatibility_mode_enabled,omitempty"`
	IbmCustomer                      string `json:"ibm_customer,omitempty"`
	IbmComponent                     string `json:"ibm_component,omitempty"`
	IbmCountry                       string `json:"ibm_country,omitempty"`
	TierScmCompressedDataUsed        string `json:"tier_scm_compressed_data_used,omitempty"`
	Tier0FlashCompressedDataUsed     string `json:"tier0_flash_compressed_data_used,omitempty"`
	Tier1FlashCompressedDataUsed     string `json:"tier1_flash_compressed_data_used,omitempty"`
	TierEnterpriseCompressedDataUsed string `json:"tier_enterprise_compressed_data_used,omitempty"`
	TierNearlineCompressedDataUsed   string `json:"tier_nearline_compressed_data_used,omitempty"`
	TotalReclaimableCapacity         string `json:"total_reclaimable_capacity,omitempty"`
	PhysicalCapacity                 string `json:"physical_capacity,omitempty"`
	PhysicalFreeCapacity             string `json:"physical_free_capacity,omitempty"`
	UsedCapacityBeforeReduction      string `json:"used_capacity_before_reduction,omitempty"`
	UsedCapacityAfterReduction       string `json:"used_capacity_after_reduction,omitempty"`
	OverheadCapacity                 string `json:"overhead_capacity,omitempty"`
	DeduplicationCapacitySaving      string `json:"deduplication_capacity_saving,omitempty"`
	EnhancedCallhome                 string `json:"enhanced_callhome,omitempty"`
	CensorCallhome                   string `json:"censor_callhome,omitempty"`
	HostUnmap                        string `json:"host_unmap,omitempty"`
	BackendUnmap                     string `json:"backend_unmap,omitempty"`
	QuorumMode                       string `json:"quorum_mode,omitempty"`
	QuorumSiteId                     string `json:"quorum_site_id,omitempty"`
	QuorumSiteName                   string `json:"quorum_site_name,omitempty"`
	QuorumLease                      string `json:"quorum_lease,omitempty"`
}

func (c *Client) PostLsSystem() *LsSystemInst {
	body, err := c.post("/rest/lssystem", nil)
	if err != nil {
		c.lastAuth = false
		return nil
	}
	if body == nil {
		return nil
	}
	var data LsSystemInst
	err = json.Unmarshal(body, &data)

	return &data
}
