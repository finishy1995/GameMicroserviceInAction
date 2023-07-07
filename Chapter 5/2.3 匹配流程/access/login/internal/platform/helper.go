package platform

const (
	// Invalid is the invalid value for the platform interface.
	Invalid = iota
	// Mock is the mock implementation of the platform interface.
	Mock
	// Steam is the Steam implementation of the platform interface.
	Steam
)

var (
	platformStr2Int = map[string]int{
		"mock":  Mock,
		"steam": Steam,
	}
	platform2Handler = map[int]Handler{
		Mock: &mockHandler{},
	}
)

func GetPlatformHandler(name string) Handler {
	if v, ok := platformStr2Int[name]; ok {
		return platform2Handler[v]
	}
	return nil
}

func Verify(platform string, token string) (string, error) {
	if v := GetPlatformHandler(platform); v != nil {
		return v.Verify(token), nil
	}
	return "", ErrInvalidPlatform
}
