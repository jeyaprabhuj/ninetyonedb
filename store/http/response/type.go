package response

type ErrorResponse struct {
	Code    int
	Message string
}

type ServerInfo struct {
	Version    string
	Package    string
	Author     string
	Storages   []string
	StorageURL string
}
