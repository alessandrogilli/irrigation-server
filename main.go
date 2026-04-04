package main

import (
	"html/template"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"schedules-app/db"
	"schedules-app/handlers"
	"schedules-app/mqtt"
	"schedules-app/scheduler"
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	err := db.InitDB()
	if err != nil {
		panic(err)
	}

	scheduler.Init()
	mqttBroker := os.Getenv("MQTT_BROKER")

	mqtt.Init(mqttBroker)
	loadScheduledJobs()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// static
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// web
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", nil)
	})

	// API
	r.Route("/api/v1", func(r chi.Router) {
		// schedules
		r.Get("/schedules", handlers.GetSchedules)
		r.Post("/schedules", handlers.CreateSchedule)
		r.Put("/schedules/{id}", handlers.UpdateSchedule)
		r.Delete("/schedules/{id}", handlers.DeleteSchedule)

		// lines
		r.Get("/lines", handlers.GetLines)
		r.Post("/lines", handlers.CreateLine)
		r.Delete("/lines/{id}", handlers.DeleteLine)
		r.Post("/lines/{id}/test", handlers.TestLine)
	})

	http.ListenAndServe(":8080", r)
}

func loadScheduledJobs() {
	rows, _ := db.DB.Query(`
        SELECT id, line_id, on_time, off_time 
        FROM schedules 
        WHERE scheduled = 1
    `)

	for rows.Next() {
		var id, lineID int
		var on, off string

		rows.Scan(&id, &lineID, &on, &off)

		scheduler.AddSchedule(id, on, off, lineID)
	}
}
