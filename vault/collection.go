package vault


type SecretCollection struct {
	Secrets map[string]*Secret
}

func NewSecretCollection() *SecretCollection {
	return &SecretCollection{
		Secrets: make(map[string]*Secret),
	}
}

func (s *SecretCollection) HasSecret(secretPath string) bool {
	if _, ok := s.Secrets[secretPath]; ok {
		return true
	} else {
		return false
	}
}

func (s *SecretCollection) AddSecret(secretPath string, secretValue *Secret) {
	s.Secrets[secretPath] = secretValue
}