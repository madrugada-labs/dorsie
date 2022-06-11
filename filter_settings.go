package main

type FilterSettings struct {
	MinSalary int
	Fields    *[]FieldEnum
}

func NewFilterSettings(userPreferences *UserPreferences) *FilterSettings {
	return &FilterSettings{
		MinSalary: userPreferences.state.MinSalary,
		Fields:    &userPreferences.state.Fields,
	}
}

func (fs *FilterSettings) UpdateFilters(userPreferences *UserPreferences) {
	fs.MinSalary = userPreferences.state.MinSalary
}
