package dynamic

// +k8s:deepcopy-gen=true

// TCPMiddleware holds the TCPMiddleware configuration.
type TCPMiddleware struct {
	InFlightConn *TCPInFlightConn `json:"InFlightConn,omitempty" toml:"InFlightConn,omitempty" yaml:"InFlightConn,omitempty" export:"true"`
	IPAllowList  *TCPIPAllowList  `json:"ipAllowList,omitempty" toml:"ipAllowList,omitempty" yaml:"ipAllowList,omitempty" export:"true"`
}

// +k8s:deepcopy-gen=true

// TCPInFlightConn holds the TCP in flight connection configuration.
type TCPInFlightConn struct {
	Amount int64 `json:"amount,omitempty" toml:"amount,omitempty" yaml:"amount,omitempty" export:"true"`
}

// +k8s:deepcopy-gen=true

// TCPIPAllowList holds the TCP ip white list configuration.
type TCPIPAllowList struct {
	SourceRange []string `json:"sourceRange,omitempty" toml:"sourceRange,omitempty" yaml:"sourceRange,omitempty"`
}
