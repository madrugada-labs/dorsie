package main

type FilterSettings struct {
	MinSalary int
}

func NewFilterSettings() *FilterSettings {
	return &FilterSettings{
		MinSalary: 0,
	}
}
