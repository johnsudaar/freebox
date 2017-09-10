package client

type RequestOpts struct {
	ExpandPath  bool
	CheckResult bool
}

type Client interface {
	Get(path string, opts RequestOpts) ([]byte, error)
	Post(path string, body interface{}, opts RequestOpts) ([]byte, error)
	Version() (Version, error)
	RequestAppToken(appName string) (string, chan error, error)
}
