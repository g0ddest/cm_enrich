package models

type EnrichmentMsg struct {
	ID                 string   `json:"id"`
	MP                 string   `json:"mp"`
	Organization       string   `json:"organization"`
	ShortDescription   string   `json:"short_description"`
	Event              string   `json:"event"`
	EventStart         string   `json:"event_start"`
	EventStop          *string  `json:"event_stop"`
	City               string   `json:"city"`
	StreetType         string   `json:"street_type"`
	Street             string   `json:"street"`
	Service            string   `json:"service"`
	HouseNumbers       []string `json:"house_numbers"`
	HouseRanges        []string `json:"house_ranges"`
	FullAddress        string   `json:"full_address,omitempty"`
	Kladr              string   `json:"kladr,omitempty"`
	RegionKladr        string   `json:"region_kladr,omitempty"`
	RegionName         string   `json:"region_name,omitempty"`
	RegionNameWithType string   `json:"region_name_with_type,omitempty"`
	RegionType         string   `json:"region_type,omitempty"`
	RegionTypeFull     string   `json:"region_type_full,omitempty"`
	StreetKladr        string   `json:"street_kladr,omitempty"`
	StreetName         string   `json:"street_name,omitempty"`
	StreetNameWithType string   `json:"street_name_with_type,omitempty"`
	StreetTypeFull     string   `json:"street_type_full,omitempty"`
}
