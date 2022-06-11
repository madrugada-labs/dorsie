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
	MinSalary int         `json:"minSalary"`
	Fields    []FieldEnum `json:"fields"`
}

func NewUserPreferences() *UserPreferences {
	return &UserPreferences{
		preferencesPath: "",
		state: &UserPreferencesState{
			MinSalary: 0,
		},
	}
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
	}

	return err
}

func (up *UserPreferences) LoadPreferences(flags Flags) (*UserPreferencesState, error) {

	jsonFile, err := os.Open(up.preferencesPath)
	defer jsonFile.Close()

	if err != nil {
		return nil, err
	}

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

	if len(flags.Fields) > 0 {
		up.state.Fields = flags.Fields
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
