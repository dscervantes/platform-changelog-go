package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redhatinsights/platform-changelog-go/internal/config"
	"github.com/redhatinsights/platform-changelog-go/internal/db"
	"github.com/redhatinsights/platform-changelog-go/internal/endpoints"
	"github.com/redhatinsights/platform-changelog-go/internal/metrics"
)

func startAPI(cfg *config.Config) {
	dbConnector := db.NewDBConnector(cfg)

	handler := endpoints.NewHandler(dbConnector)

	r := chi.NewRouter()
	mr := chi.NewRouter()
	sub := chi.NewRouter().With(metrics.ResponseMetricsMiddleware)

	// Mount the root of the api router on /api/v1
	r.Mount("/api/v1", sub)
	r.Get("/", handler.LubdubHandler)

	mr.Get("/", handler.LubdubHandler)
	mr.Get("/healthz", handler.LubdubHandler)
	mr.Handle("/metrics", promhttp.Handler())

	sub.Post("/github", handler.Github)
	sub.Post("/github-webhook", handler.GithubWebhook)
	sub.Post("/gitlab-webhook", handler.GitlabWebhook)
	sub.Post("/tekton", handler.TektonTaskRun)

	sub.Get("/services", handler.GetServicesAll)
	sub.Get("/projects", handler.GetProjectsAll)
	sub.Get("/timelines", handler.GetTimelinesAll)
	sub.Get("/commits", handler.GetCommitsAll)
	sub.Get("/deploys", handler.GetDeploysAll)

	sub.Get("/services/{service}", handler.GetServiceByName)
	sub.Get("/services/{service}/projects", handler.GetProjectsByService)
	sub.Get("/services/{service}/timelines", handler.GetTimelinesByService)
	sub.Get("/services/{service}/commits", handler.GetCommitsByService)
	sub.Get("/services/{service}/deploys", handler.GetDeploysByService)

	sub.Get("/projects/{project}", handler.GetProjectByName)
	sub.Get("/projects/{project}/timelines", handler.GetTimelinesByProject)
	sub.Get("/projects/{project}/commits", handler.GetCommitsByProject)
	sub.Get("/projects/{project}/deploys", handler.GetDeploysByProject)

	sub.Get("/timelines/{ref}", handler.GetTimelineByRef)
	sub.Get("/commits/{ref}", handler.GetCommitByRef)
	sub.Get("/deploys/{ref}", handler.GetDeployByRef)

	sub.Get("/openapi.json", handler.OpenAPIHandler(cfg))

	srv := http.Server{
		Addr:    ":" + cfg.PublicPort,
		Handler: r,
	}

	msrv := http.Server{
		Addr:    ":" + cfg.MetricsPort,
		Handler: mr,
	}

	go func() {
		if err := msrv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
