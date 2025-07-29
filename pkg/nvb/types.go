package nvb

// FunctionTitlesResponse represents the response from the function titles API
type FunctionTitlesResponse struct {
	Suggestions []string `json:"suggestions"`
}

// Coordinates represents geographic coordinates with proper numeric types
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// GeoLocationResponse represents the response from the geolocation API
type GeoLocationResponse struct {
	CityCenter Coordinates `json:"cityCenter"`
	CityName   string      `json:"cityName"`
}

// JobsResponse represents the response from the jobs search API
type JobsResponse struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Pages int `json:"pages"`
	Total int `json:"total"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		First struct {
			Href string `json:"href"`
		} `json:"first"`
		Last struct {
			Href string `json:"href"`
		} `json:"last"`
		Next struct {
			Href string `json:"href"`
		} `json:"next"`
	} `json:"_links"`
	Embedded struct {
		Jobs []Job `json:"jobs"`
	} `json:"_embedded"`
}

// Job represents a single job listing
type Job struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	DcoTitle    string `json:"dcoTitle"`
	Description string `json:"description"`
	Company     struct {
		Name    string `json:"name"`
		Website string `json:"website"`
		Slug    string `json:"slug"`
		Type    string `json:"type"`
	} `json:"company"`
	Salary struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"salary"`
	ContractType string   `json:"contractType"`
	CareerLevel  string   `json:"careerLevel"`
	Categories   []string `json:"categories"`
	Industries   []string `json:"industries"`
	StartDate    string   `json:"startDate"`
	EndDate      string   `json:"endDate"`
	Status       string   `json:"status"`
	WorkingHours struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"workingHours"`
}

// JobSearchParams holds all parameters for job search
type JobSearchParams struct {
	Page      int
	Limit     int
	Sort      string
	City      string
	JobTitle  string
	Latitude  float64
	Longitude float64
	Distance  float64
}

// NationalevacaturebankAPI defines the interface for the National Evacature Bank API
type NationalevacaturebankAPI interface {
	SearchFunctionTitles(query string) (*FunctionTitlesResponse, error)
	SearchCities(startsWith string) ([]string, error)
	GetGeoLocation(cityName string) (*GeoLocationResponse, error)
	FindJobs(params JobSearchParams) (*JobsResponse, error)
}