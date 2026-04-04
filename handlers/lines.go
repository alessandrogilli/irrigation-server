package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"schedules-app/db"
	"schedules-app/models"
	"schedules-app/mqtt"
)

// GET /lines
func GetLines(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name, is_on FROM lines")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var lines []models.Line

	for rows.Next() {
		var l models.Line
		var isOnInt int
		rows.Scan(&l.ID, &l.Name, &isOnInt)
		l.IsOn = isOnInt == 1
		lines = append(lines, l)
	}

	json.NewEncoder(w).Encode(lines)
}

// POST /lines
func CreateLine(w http.ResponseWriter, r *http.Request) {
	var l models.Line
	json.NewDecoder(r.Body).Decode(&l)

	if l.Name == "" {
		http.Error(w, "name required", 400)
		return
	}

	if l.ID == 0 { // if ID is not provided, let sqlite auto-increment it. Otherwise, use the provided ID (useful for testing)
		db.DB.Exec("INSERT INTO lines(name, is_on) VALUES (?, 0)", l.Name)
	} else {
		db.DB.Exec("INSERT INTO lines(id, name, is_on) VALUES (?, ?, 0)", l.ID, l.Name)
	}

	json.NewEncoder(w).Encode(l)
}

// DELETE /lines/{id}
func DeleteLine(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	db.DB.Exec("DELETE FROM lines WHERE id=?", id)
	w.WriteHeader(http.StatusNoContent)
}

// POST /lines/{id}/test
func TestLine(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var req struct {
		Cmd string `json:"cmd"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// Determine the new state
	isOn := 0
	if req.Cmd == "ON" {
		isOn = 1
	}

	// Update database
	db.DB.Exec("UPDATE lines SET is_on=? WHERE id=?", isOn, id)

	// Publish to MQTT
	mqtt.Publish(id, req.Cmd)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"is_on": isOn})
}
