package utils

import "github.com/mssola/user_agent"

type UaResult struct {
	Os      string `json:"os"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Model   string `json:"model"`
}

// Ua
func Ua(userAgent string) UaResult {
	ua := user_agent.New(userAgent)
	browser, s := ua.Browser()
	return UaResult{
		Os:      ua.OS(),
		Name:    browser,
		Version: s,
		Model:   ua.Platform(),
	}
}
