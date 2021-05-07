package controllers

import (
	"go-app/app/models"
	"go-app/bmonster"
)

func UpdateLatestPrograms(apiKey string) error {
	apiClient := bmonster.New(apiKey)
	programList, err := apiClient.GetLatestProgramList()
	if err != nil {
		return err
	}

	for _, program := range programList {
		if p := models.GetProgram(program.StudioName, program.StartTime); p != nil {
			program.Save()
		} else {
			program.Create()
		}
	}
	return nil
}
