package unisphere

type HealthEnum int

const (
	HealthEnumUnknown        HealthEnum = 0
	HealthEnumOk             HealthEnum = 5
	HealthEnumOkBut          HealthEnum = 7
	HealthEnumDegraded       HealthEnum = 10
	HealthEnumMinor          HealthEnum = 15
	HealthEnumMajor          HealthEnum = 20
	HealthEnumCritical       HealthEnum = 25
	HealthEnumNonRecoverable HealthEnum = 30
)

func (h *HealthEnum) String() string {
	switch *h {
	case HealthEnumUnknown:
		return "Unknown"
	case HealthEnumOk:
		return "Ok"
	case HealthEnumOkBut:
		return "OkBut"
	case HealthEnumDegraded:
		return "Degraded"
	case HealthEnumMinor:
		return "Minor"
	case HealthEnumMajor:
		return "Major"
	case HealthEnumCritical:
		return "Critical"
	case HealthEnumNonRecoverable:
		return "NonRecoverable"
	default:
		return "Unknown"
	}
}

type Health struct {
	Value          HealthEnum `json:"value"`
	DescriptionIds []string   `json:"descriptionIds"`
	Descriptions   []string   `json:"descriptions"`
	ResolutionIds  []string   `json:"resolutionIds"`
	Resolutions    []string   `json:"resolutions"`
}
