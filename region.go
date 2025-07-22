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
