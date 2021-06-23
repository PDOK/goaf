package provider

type CommonProvider struct {
	ServiceSpecPath    string
	ServiceEndpoint    string
	MaxReturnLimit     uint64
	DefaultReturnLimit uint64
}

func NewCommonProvider(serviceEndpoint, serviceSpecPath string, defaultReturnLimit uint64, maxReturnLimit uint64) CommonProvider {
	return CommonProvider{
		ServiceEndpoint:    serviceEndpoint,
		ServiceSpecPath:    serviceSpecPath,
		DefaultReturnLimit: defaultReturnLimit,
		MaxReturnLimit:     maxReturnLimit,
	}
}
