package localization

// Translator defines the interface for a translator
type Translator interface {
	Tr(key string, args ...interface{}) string
}
