package scraperapp

import (
	"encoding/json"
	"net/http"

	"github.com/ctfrancia/maple/app/sdk/web"
	"github.com/ctfrancia/maple/business/domain/scraperbus"
	"github.com/ctfrancia/maple/business/sdk/order"
	"github.com/ctfrancia/maple/business/sdk/page"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type app struct {
	scraperBus   scraperbus.ExtBusiness
	orchestrator Orchestrator
}

func newApp(scraperBus scraperbus.ExtBusiness, orchestrator Orchestrator) *app {
	return &app{
		scraperBus:   scraperBus,
		orchestrator: orchestrator,
	}
}

// createJob creates a new scrape job
func (a *app) createJob(w http.ResponseWriter, r *http.Request) (*web.Response, error) {
	var nj NewScrapeJob
	if err := json.NewDecoder(r.Body).Decode(&nj); err != nil {
		return &web.Response{Status: http.StatusBadRequest}, err
	}

	if nj.Federation == "" {
		return &web.Response{
			Status: http.StatusBadRequest,
			Data:   map[string]string{"error": "federation is required"},
		}, nil
	}

	job, err := a.scraperBus.CreateJob(r.Context(), scraperbus.NewScrapeJob{
		Federation: nj.Federation,
	})
	if err != nil {
		return nil, err
	}

	return &web.Response{
		Data:   toAppScrapeJob(job),
		Status: http.StatusCreated,
	}, nil
}

// queryJobs retrieves a list of scrape jobs
func (a *app) queryJobs(w http.ResponseWriter, r *http.Request) (*web.Response, error) {
	filter := scraperbus.QueryFilter{}

	if fed := r.URL.Query().Get("federation"); fed != "" {
		filter.Federation = &fed
	}
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = &status
	}

	jobs, err := a.scraperBus.QueryJobs(r.Context(), filter, order.By{}, page.MustParse("1", "50"))
	if err != nil {
		return nil, err
	}

	appJobs := make([]ScrapeJob, len(jobs))
	for i, job := range jobs {
		appJobs[i] = toAppScrapeJob(job)
	}

	return &web.Response{
		Data:   appJobs,
		Status: http.StatusOK,
	}, nil
}

// fetchJob retrieves a scrape job by ID
func (a *app) fetchJob(w http.ResponseWriter, r *http.Request) (*web.Response, error) {
	jobID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return &web.Response{Status: http.StatusBadRequest}, err
	}

	job, err := a.scraperBus.QueryJobByID(r.Context(), jobID)
	if err != nil {
		return nil, err
	}

	return &web.Response{
		Data:   toAppScrapeJob(job),
		Status: http.StatusOK,
	}, nil
}

// queryTournaments retrieves a list of scraped tournaments
func (a *app) queryTournaments(w http.ResponseWriter, r *http.Request) (*web.Response, error) {
	filter := scraperbus.QueryFilter{}

	if fed := r.URL.Query().Get("federation"); fed != "" {
		filter.Federation = &fed
	}

	tournaments, err := a.scraperBus.QueryTournaments(r.Context(), filter, order.By{}, page.MustParse("1", "50"))
	if err != nil {
		return nil, err
	}

	appTournaments := make([]ScrapedTournament, len(tournaments))
	for i, t := range tournaments {
		appTournaments[i] = toAppScrapedTournament(t)
	}

	return &web.Response{
		Data:   appTournaments,
		Status: http.StatusOK,
	}, nil
}

// triggerScrape manually triggers a scrape for a federation
func (a *app) triggerScrape(w http.ResponseWriter, r *http.Request) (*web.Response, error) {
	var req struct {
		Federation string `json:"federation"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &web.Response{Status: http.StatusBadRequest}, err
	}

	if req.Federation == "" {
		return &web.Response{
			Status: http.StatusBadRequest,
			Data:   map[string]string{"error": "federation is required"},
		}, nil
	}

	jobID, err := a.orchestrator.TriggerScrape(r.Context(), req.Federation)
	if err != nil {
		return nil, err
	}

	return &web.Response{
		Data: map[string]string{
			"job_id":     jobID.String(),
			"federation": req.Federation,
			"message":    "scrape job triggered",
		},
		Status: http.StatusAccepted,
	}, nil
}

// getStatistics returns aggregated statistics about scraping
func (a *app) getStatistics(w http.ResponseWriter, r *http.Request) (*web.Response, error) {
	// Query all jobs
	allJobs, err := a.scraperBus.QueryJobs(r.Context(), scraperbus.QueryFilter{}, order.By{}, page.MustParse("1", "1000"))
	if err != nil {
		return nil, err
	}

	// Calculate statistics
	stats := map[string]interface{}{
		"total_jobs": len(allJobs),
		"by_status":  make(map[string]int),
		"by_federation": make(map[string]int),
		"total_found": 0,
		"total_created": 0,
		"total_updated": 0,
	}

	byStatus := stats["by_status"].(map[string]int)
	byFederation := stats["by_federation"].(map[string]int)

	for _, job := range allJobs {
		byStatus[job.Status]++
		byFederation[job.Federation]++
		stats["total_found"] = stats["total_found"].(int) + job.Found
		stats["total_created"] = stats["total_created"].(int) + job.Created
		stats["total_updated"] = stats["total_updated"].(int) + job.Updated
	}

	return &web.Response{
		Data:   stats,
		Status: http.StatusOK,
	}, nil
}
