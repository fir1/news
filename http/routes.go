package http

func (s *Service) routes() {
	s.router.Get("/health", s.GetHealth)
	s.router.Get("/news", s.listNews)
	s.router.Get("/article", s.getArticle)
}
