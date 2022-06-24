package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type UserPreferences struct {
	preferencesPath string
	state           *UserPreferencesState
}

type UserPreferencesState struct {
	MinSalary   int              `json:"minSalary"`
	Fields      []FieldEnum      `json:"fields"`
	Experiences []ExperienceEnum `json:"experiences"`
	SkipIntro   bool             `json:"skipIntro"`
}

func NewUserPreferences() *UserPreferences {
	return &UserPreferences{
		preferencesPath: "",
		state: &UserPreferencesState{
			MinSalary:   0,
			Fields:      nil,
			Experiences: []ExperienceEnum{"early_career", "mid_level", "senior"},
			SkipIntro:   false,
		},
	}
}

func (up *UserPreferences) SkipIntroEnabled() bool {
	return up.state.SkipIntro
}

// creates the preferences file only if it does not exist
func (up *UserPreferences) CreatePreferencesFile() error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	newpath := filepath.Join(dirname, ".dorse")
	err = os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		return err
	}

	// if the file does not exist, then it creates it with the current state
	up.preferencesPath = filepath.Join(newpath, "preferences.json")
	_, err = os.Stat(filepath.Join(up.preferencesPath))
	if errors.Is(err, os.ErrNotExist) {
		_, err = os.Create(filepath.Join(up.preferencesPath))
		preferencesBytes, err := json.Marshal(up.state)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(up.preferencesPath, preferencesBytes, 0)
		if err != nil {
			return err
		}
	}

	return err
}

func (up *UserPreferences) LoadPreferences(flags Flags) (*UserPreferencesState, error) {

	jsonFile, err := os.Open(up.preferencesPath)

	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	var ups UserPreferencesState
	jsonFileBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonFileBytes, &ups)
	if err != nil {
		return nil, err
	}
	up.state = &ups

	// update all non nil flags into config
	if flags.MinSalary != nil && *flags.MinSalary > 0 {
		up.state.MinSalary = *flags.MinSalary
	}

	if len(flags.Experiences) > 0 {
		up.state.Experiences = flags.Experiences
	}

	if len(flags.Fields) > 0 {
		up.state.Fields = flags.Fields
	}

	switch *flags.SkipIntro {
	case true:
		ups.SkipIntro = true
	case false:
		if *flags.MinSalary != -1 || flags.Fields != nil {
			ups.SkipIntro = true
		} else {
			ups.SkipIntro = false
		}
	}

	return &ups, nil
}

func (up *UserPreferences) PersistPreferences(newState *UserPreferencesState) error {
	preferencesBytes, err := json.Marshal(newState)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(up.preferencesPath, preferencesBytes, 0)
	up.state = newState
	return err
}
