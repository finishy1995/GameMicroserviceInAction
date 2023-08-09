package platform

type mockHandler struct {
}

func (m mockHandler) Verify(token string) string {
	return token
}
