package application

import "net/http"

func (app *Application) getAllScheduleTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := app.models.ScheduleTemplate.SelectAll()
	if err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	app.successResponse(w, r, "获取班表模板成功", templates)
}
