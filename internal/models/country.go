package models

type Country struct {
	IpAdders string `json:"ip_adders,omitempty"`
	Country  string `json:"country"`
	City     string `json:"city"`
}
