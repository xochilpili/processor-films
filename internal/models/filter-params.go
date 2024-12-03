package models

type FilterParams struct {
	Provider   string `json:"provider,omitempty"`
	Term       string `json:"term"`
	Resolution string `json:"resolution"`
}
