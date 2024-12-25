package application

import "net/http"

func (app *Application) getAllScheduleTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	templates, err := app.models.ScheduleTemplate.SelectAll()
	if err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	app.successResponse(w, r, "获取班表模板成功", templates)
}

func (app *Application) createScheduleTemplateHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := app.readJSON(r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}
	if err := app.validate.Struct(&payload); err != nil {
		app.validateError(w, r, err)
		return
	}

	// check whether the name exists
	exists, err := app.models.ScheduleTemplate.CheckTemplateNameValid(payload.Name)
	if err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	if exists {
		app.errorResponse(w, r, errScheduleTemplateExists)
		return
	}

	// insert the template
	template, err := app.models.ScheduleTemplate.InsertScheduleTemplate(payload.Name, payload.Description)
	if err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	app.successResponse(w, r, "创建模板成功", template)
}
