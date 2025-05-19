package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type SiteStatus struct {
	URL       string    `json:"url"`
	Status    string    `json:"status"`
	CheckedAt time.Time `json:"checked_at"`
}

var (
	sites = []string{} // Lista de sites para monitorar
	mu    sync.Mutex   // Mutex para evitar conflitos de concorrência
)

func checkSite(url string) SiteStatus {
	_, err := http.Get(url)
	status := "Online"
	if err != nil {
		status = "Offline"
	}
	return SiteStatus{
		URL:       url,
		Status:    status,
		CheckedAt: time.Now(),
	}
}

func addSiteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		URL string `json:"url"`
	}

	// Decodificar JSON recebido
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Erro ao decodificar JSON: %v\n", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validar URL
	parsedURL, err := url.Parse(data.URL)
	if err != nil || parsedURL.Host == "" {
		log.Printf("URL inválida recebida: %s\n", data.URL)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Adicionar esquema se ausente
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
	}

	finalURL := parsedURL.String()

	// Verificar duplicatas
	mu.Lock()
	defer mu.Unlock()
	for _, site := range sites {
		if site == finalURL {
			log.Printf("URL já adicionada: %s\n", finalURL)
			http.Error(w, "Site already added", http.StatusConflict)
			return
		}
	}

	// Adicionar site à lista
	sites = append(sites, finalURL)
	log.Printf("Site adicionado: %s\n", finalURL)
	w.WriteHeader(http.StatusCreated)
}

func checkSitesHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var statuses []SiteStatus
	for _, site := range sites {
		statuses = append(statuses, checkSite(site))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

func main() {
	http.HandleFunc("/api/check-sites", checkSitesHandler)
	http.HandleFunc("/api/add-site", addSiteHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
