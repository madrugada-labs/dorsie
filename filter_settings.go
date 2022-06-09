package main

type FilterSettings struct {
	MinSalary int
}

func NewFilterSettings(userPreferences *UserPreferences) *FilterSettings {
	return &FilterSettings{
		MinSalary: userPreferences.state.MinSalary,
	}
}

func (fs *FilterSettings) UpdateFilters(userPreferences *UserPreferences) {
	fs.MinSalary = userPreferences.state.MinSalary
}
