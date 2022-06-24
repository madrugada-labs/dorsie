package main

/// FilterSettings tracks which settings are persisted and loaded
/// so that clients can use them in their queries
type FilterSettings struct {
	MinSalary   int
	Experiences []ExperienceEnum
	Fields      *[]FieldEnum
}

/// NewFilterSettings creates a new FilterSettings
func NewFilterSettings(userPreferences *UserPreferences) *FilterSettings {
	return &FilterSettings{
		MinSalary:   userPreferences.state.MinSalary,
		Experiences: userPreferences.state.Experiences,
		Fields:      &userPreferences.state.Fields,
	}
}

/// UpdateFilters updates the user preferences
func (fs *FilterSettings) UpdateFilters(userPreferences *UserPreferences) {
	fs.MinSalary = userPreferences.state.MinSalary
}
