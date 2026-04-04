package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"schedules-app/db"
	"schedules-app/models"
	"schedules-app/scheduler"
)

// GET /schedules
func GetSchedules(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
	SELECT s.id, s.line_id, l.name, s.on_time, s.off_time, s.scheduled, s.description
	FROM schedules s
	LEFT JOIN lines l ON s.line_id = l.id
	`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var schedules []models.Schedule

	for rows.Next() {
		var s models.Schedule
		var scheduledInt int

		err := rows.Scan(&s.ID, &s.LineID, &s.LineName, &s.OnTime, &s.OffTime, &scheduledInt, &s.Description)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		s.Scheduled = scheduledInt == 1

		schedules = append(schedules, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedules)
}

// POST /schedules
func CreateSchedule(w http.ResponseWriter, r *http.Request) {
	var s models.Schedule

	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if s.LineID == 0 || s.OnTime == "" || s.OffTime == "" {
		http.Error(w, "missing fields", 400)
		return
	}

	res, err := db.DB.Exec(
		"INSERT INTO schedules(line_id, on_time, off_time, scheduled, description) VALUES (?, ?, ?, 0, ?)",
		s.LineID, s.OnTime, s.OffTime, s.Description,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	id, _ := res.LastInsertId()
	s.ID = int(id)
	s.Scheduled = false

	log.Printf("INPUT: %+v\n", s)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

// PUT /schedules/{id}
func UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	var s models.Schedule
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	_, err = db.DB.Exec(
		"UPDATE schedules SET line_id=?, on_time=?, off_time=?, scheduled=?, description=? WHERE id=?",
		s.LineID,
		s.OnTime,
		s.OffTime,
		boolToInt(s.Scheduled),
		s.Description,
		id,
	)

	if s.Scheduled {
		scheduler.AddSchedule(id, s.OnTime, s.OffTime, s.LineID)
	} else {
		scheduler.RemoveSchedule(id)
	}

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	s.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// DELETE /schedules/{id}
func DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	_, err := db.DB.Exec("DELETE FROM schedules WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
