package insight

func (im *InsightManager) ValidateRequiredFields() error {
	if err := im.BaseManager.ValidateRequiredFields(); err != nil {
		return err
	}

	return nil
}
