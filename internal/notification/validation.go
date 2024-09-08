package notification

func validateCategories(category string) bool {

	// Validate the category
	validCategories := map[string]bool{
		"Sports":  true,
		"Finance": true,
		"Films":   true,
	}
	return validCategories[category]
}
