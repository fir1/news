package http

import "net/http"

// GetHealth example
//
//	@Summary		Get health of server
//	@Description	Get health of server
//	@Tags Health-Server
//	@ID				get-health
//	@Accept			json
//	@Produce		json
//
// @Success      200
// @Failure      500
// @Router			/health [get].
func (s *Service) GetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		s.logger.Errorf("health write error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
