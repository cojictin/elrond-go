package factory

// SetSkPkProviderHandler updates the handler for testing reasons
func (cspf *cryptoSigningParamsFactory) SetSkPkProviderHandler(handler func() ([]byte, []byte, error)) {
	cspf.skPkProviderHandler = handler
}

// GetSkPk will call the inner function
func (cspf *cryptoSigningParamsFactory) GetSkPk() ([]byte, []byte, error) {
	return cspf.getSkPk()
}
