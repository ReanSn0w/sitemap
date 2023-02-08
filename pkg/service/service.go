package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ReanSn0w/sitemap/pkg/sitemap"
)

func New(host string, static bool, builders ...sitemap.SitemapBuilder) *Service {
	return &Service{
		sm:     sitemap.NewSitemap(host, builders...),
		static: static,
	}
}

type Service struct {
	srv *http.Server
	sm  sitemap.Sitemap

	static bool
}

// run service with graceful shutdown
func (s *Service) Run() {
	s.sm.Run(time.Minute * 10)

	log.Println("Запуск сервера")
	s.srv = &http.Server{
		Addr:    ":8080",
		Handler: s.Handler(),
	}

	go s.gracefulShutdown()

	if err := s.srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (s *Service) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sitemap.xml" {
			http.FileServer(http.Dir("static")).ServeHTTP(w, r)
			return
		}

		s.sm.Handler().ServeHTTP(w, r)
	})
}

func (s *Service) gracefulShutdown() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	registretSignal := <-quit
	log.Printf("Зарегистрирован системный сигнал: %s", registretSignal.String())

	log.Println("Производится отключение сервера")
	err := s.srv.Shutdown(context.Background())
	if err != nil {
		log.Println(err)
	}
}
