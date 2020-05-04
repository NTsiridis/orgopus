package domain

import (
	"crypto/rand"
	"fmt"
	"log"
	"testing"
	"time"
)

// ---- helper types and functions ----
type mockTTEntity struct {
	id        string
	startFrom time.Time
	endAt     time.Time
}

func (m mockTTEntity) IsExistentAt(pit time.Time) bool {

	if m.startFrom.After(pit) {
		return false
	}

	if m.endAt.IsZero() {
		return true
	}

	return m.endAt.After(pit)
}

func (m mockTTEntity) ExistentFrom() time.Time {
	return m.startFrom
}

func (m mockTTEntity) ValidUntil() time.Time {
	return m.endAt
}

func (m mockTTEntity) ActiveDuration() time.Duration {

	ending := time.Now()
	if !m.endAt.IsZero() {
		ending = m.endAt
	}

	return ending.Sub(m.startFrom)
}

func (m mockTTEntity) String() string {
	endingDate := ""
	if !m.endAt.IsZero() {
		endingDate = m.endAt.Format("2006-01-02 15:04:05")
	}

	return fmt.Sprintf("%s [%s -- %s]",
		m.id[len(m.id)-4:],
		m.startFrom.Format("2006-01-02 15:04:05"),
		endingDate)
}

func createMockUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func createMockTTEntity(start time.Time, end time.Time) TimeTrackedEntity {
	return mockTTEntity{
		startFrom: start,
		endAt:     end,
		id:        createMockUUID(),
	}

}

// ------------------ Tests -------

func TestAddEntityToSlice(t *testing.T) {

	collection := TimeTrackedEntityCollection{}

	collection.AddEntity(createMockTTEntity(
		time.Now(),
		NilTime()))
	collection.AddEntity(createMockTTEntity(
		time.Date(2020, 1, 2, 15, 30, 10, 0, time.Local),
		time.Date(2020, 1, 4, 15, 30, 10, 0, time.Local)))
	collection.AddEntity(createMockTTEntity(
		time.Date(2020, 1, 6, 15, 30, 10, 0, time.Local),
		time.Date(2020, 1, 8, 15, 30, 10, 0, time.Local)))
	collection.AddEntity(createMockTTEntity(
		time.Date(2020, 1, 3, 15, 30, 10, 0, time.Local),
		NilTime()))

	fmt.Printf("Collection:\n%v\n", collection)

}
