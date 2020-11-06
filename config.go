package main

import (
	"fmt"
	"log"

	"gopkg.in/ini.v1"
)

// RecordUpdater ...
type RecordUpdater struct {
	configFile string
	config     *ini.File
}

// NewRecordUpdater ...
func NewRecordUpdater(file string) RecordUpdater {
	cfg, err := ini.Load(file)
	if err != nil {
		log.Fatal(err)
	}

	c := RecordUpdater{
		configFile: file,
		config:     cfg,
	}
	return c
}

// GetString ...
func (ru *RecordUpdater) GetString(key string, section string) (string, error) {
	if ru.config.Section(section).HasKey(key) {
		return ru.config.Section(section).Key(key).String(), nil
	}

	return "", fmt.Errorf("Missing key '%s'", key)
}

// GetStrings ...
func (ru *RecordUpdater) GetStrings(key string, del string, section string) ([]string, error) {
	if ru.config.Section(section).HasKey(key) {
		return ru.config.Section(section).Key(key).Strings(del), nil
	}

	return nil, fmt.Errorf("Missing key '%s'", key)
}
