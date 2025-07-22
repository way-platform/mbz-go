package mbz

// Region is a region of the Mercedes-Benz Management API.
type Region string

// Regions of the Mercedes-Benz Management API.
const (
	// RegionECE is the ECE region.
	RegionECE = "ECE"
	// RegionAMAPNA is the AMAP/NA region.
	RegionAMAPNA = "AMAP/NA"
)

// Known base URLs for the Mercedes-Benz Management API.
const (
	// BaseURLECE is the base URL for the ECE region.
	BaseURLECE = "https://service.connect-business.net/api"
	// BaseURLAMAPNA is the base URL for the AMAP/NA region.
	BaseURLAMAPNA = "https://service.amap.connect-business.net/api"
)

// Known Kafka bootstrap servers.
const (
	// KafkaBootstrapServerECE is the Kafka bootstrap server for the ECE region.
	KafkaBootstrapServerECE = "bootstrap.streaming.connect-business.net:443"
	// KafkaBootstrapServerAMAPNA is the Kafka bootstrap server for the AMAP/NA region.
	KafkaBootstrapServerAMAPNA = "bootstrap.streaming.amap.connect-business.net:443"
)

// Known token URLs.
const (
	// TokenURLECE is the token URL for the ECE region.
	TokenURLECE = "https://ssoalpha.dvb.corpinter.net/v1/token"
	// TokenURLAMAPNA is the token URL for the AMAP/NA region.
	TokenURLAMAPNA = "https://ssoalpha.am.dvb.corpinter.net/v1/token"
)

// Known audience scopes.
const (
	// AudienceScopeECE is the audience scope for the ECE region.
	AudienceScopeECE = "audience:server:client_id:95B37AC2-D501-4CFD-B853-7D299DD2D872"
	// AudienceScopeAMAPNA is the audience scope for the AMAP/NA region.
	AudienceScopeAMAPNA = "audience:server:client_id:87012BCA-0B2E-4127-BE24-97A71C1F3262"
)
