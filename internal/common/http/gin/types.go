package gin

import "net/url"

type RequestLog struct {
	Request  RequestMetadata  `json:"request"`
	Response ResponseMetadata `json:"response"`
}

type RequestMetadata struct {
	Method   string     `json:"method"`
	URL      string     `json:"url"`
	Path     string     `json:"path"`
	Query    url.Values `json:"query"`
	RemoteIP string     `json:"remote_ip"`
	UserID   string     `json:"user_id,omitempty"`
	Body     string     `json:"body,omitempty"`
}

type ResponseMetadata struct {
	Status    int    `json:"status"`
	LatencyMS int64  `json:"latency_ms"`
	Body      string `json:"body,omitempty"`
}
