package platform

type Handler interface {
	Verify(token string) string
}
