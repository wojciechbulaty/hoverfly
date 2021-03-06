package journal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"
)

type JournalEntry struct {
	Request     *models.RequestDetails
	Response    *models.ResponseDetails
	Mode        string
	TimeStarted time.Time
	Latency     time.Duration
}

type Journal struct {
	entries    []JournalEntry
	EntryLimit int
}

func NewJournal() *Journal {
	return &Journal{
		entries:    []JournalEntry{},
		EntryLimit: 1000,
	}
}

func (this *Journal) NewEntry(request *http.Request, response *http.Response, mode string, started time.Time) error {
	if this.EntryLimit == 0 {
		return fmt.Errorf("Journal disabled")
	}

	payloadRequest, _ := models.NewRequestDetailsFromHttpRequest(request)

	respBody, _ := util.GetResponseBody(response)

	payloadResponse := &models.ResponseDetails{
		Status:  response.StatusCode,
		Body:    string(respBody),
		Headers: response.Header,
	}

	if len(this.entries) >= this.EntryLimit {
		this.entries = append(this.entries[:0], this.entries[1:]...)
	}

	this.entries = append(this.entries, JournalEntry{
		Request:     &payloadRequest,
		Response:    payloadResponse,
		Mode:        mode,
		TimeStarted: started,
		Latency:     time.Since(started),
	})

	return nil
}

func (this Journal) GetEntries() ([]v2.JournalEntryView, error) {
	if this.EntryLimit == 0 {
		return []v2.JournalEntryView{}, fmt.Errorf("Journal disabled")
	}

	journalEntryViews := []v2.JournalEntryView{}
	for _, journalEntry := range this.entries {
		journalEntryViews = append(journalEntryViews, v2.JournalEntryView{
			Request:     journalEntry.Request.ConvertToRequestDetailsView(),
			Response:    journalEntry.Response.ConvertToResponseDetailsView(),
			Mode:        journalEntry.Mode,
			TimeStarted: journalEntry.TimeStarted.Format(time.RFC3339),
			Latency:     (journalEntry.Latency / time.Millisecond),
		})
	}
	return journalEntryViews, nil
}

func (this *Journal) DeleteEntries() error {
	if this.EntryLimit == 0 {
		return fmt.Errorf("Journal disabled")
	}

	this.entries = []JournalEntry{}

	return nil
}
