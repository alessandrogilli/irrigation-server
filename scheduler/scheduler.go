package scheduler

import (
	"fmt"
	"sync"

	"schedules-app/mqtt"

	"github.com/robfig/cron/v3"
)

type JobIds struct {
	OnID  cron.EntryID
	OffID cron.EntryID
}

var Cron *cron.Cron
var jobs = make(map[int]JobIds)
var mu sync.Mutex

func Init() {
	Cron = cron.New()
	Cron.Start()
}

func AddSchedule(id int, onTime string, offTime string, lineID int) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := jobs[id]; exists {
		return
	}

	onSpec := toCronSpec(onTime)
	offSpec := toCronSpec(offTime)

	onID, _ := Cron.AddFunc(onSpec, func() {
		fmt.Printf("[ON ] Schedule %d triggered at %s\n", id, onTime)
		mqtt.Publish(lineID, "ON")
	})

	offID, _ := Cron.AddFunc(offSpec, func() {
		fmt.Printf("[OFF] Schedule %d triggered at %s\n", id, offTime)
		mqtt.Publish(lineID, "OFF")
	})

	jobs[id] = JobIds{OnID: onID, OffID: offID}

	fmt.Printf("Scheduled job %d (ON: %s, OFF: %s)\n", id, onTime, offTime)

	for cronID, entry := range Cron.Entries() {
		fmt.Printf("Cron entry %d: next run at %s\n", cronID, entry.Next)
	}
}

func RemoveSchedule(id int) {
	mu.Lock()
	defer mu.Unlock()

	if jobID, ok := jobs[id]; ok {
		Cron.Remove(jobID.OnID)
		Cron.Remove(jobID.OffID)
		delete(jobs, id)
		fmt.Printf("Removed schedule %d\n", id)
	}
}

// "18:30" → "30 18 * * *"
func toCronSpec(t string) string {
	var h, m int
	fmt.Sscanf(t, "%d:%d", &h, &m)
	return fmt.Sprintf("%d %d * * *", m, h)
}
