package payloads

type PayloadBuilder interface {
	BuildPayload(values map[string]string) ([]byte, error)
}
