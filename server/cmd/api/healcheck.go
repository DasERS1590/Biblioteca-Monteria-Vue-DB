package main

import (
	"net/http"
)

// healthcheckHandler godoc
// @Summary      Check system status
// @Description  Provides information about the current system status and version
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  object{status=string,system_info=object{environment=string,version=string}}
// @Router       /healthcheck [get]
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
