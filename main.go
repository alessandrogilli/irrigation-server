package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	if err := mqtt.Init(mqttBroker); err != nil {
		fmt.Printf("MQTT init error: %v\n", err)
	}
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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// start server
	go func() {
		fmt.Println("Starting HTTP server :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	// stop scheduler
	if scheduler.Cron != nil {
		scheduler.Cron.Stop()
	}

	// disconnect mqtt client if connected
	if mqtt.Client != nil {
		if mqtt.Client.IsConnected() {
			mqtt.Client.Disconnect(250)
		}
	}

	// close DB
	if db.DB != nil {
		db.DB.Close()
	}

	fmt.Println("Server stopped")
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
