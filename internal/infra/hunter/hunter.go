package hunter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HunterClient struct {
	apiKey string
	client *http.Client
}

func NewHunterClient(apiKey string) *HunterClient {
	return &HunterClient{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type HunterEmail struct {
	Value       string `json:"value"`
	Type        string `json:"type"`
	Confidence  int    `json:"confidence"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Position    string `json:"position"`
	Seniority   string `json:"seniority"`
	Department  string `json:"department"`
	LinkedinURL string `json:"linkedin_url"`
}

type hunterData struct {
	Domain  string        `json:"domain"`
	Emails  []HunterEmail `json:"emails"`
	Pattern string        `json:"pattern"`
}

type hunterResponse struct {
	Data hunterData `json:"data"`
	Meta struct {
		Results int `json:"results"`
	} `json:"meta"`
}

func (h *HunterClient) ExtractDomain(websiteURL string) string {
	if websiteURL == "" {
		return ""
	}

	if !strings.HasPrefix(websiteURL, "http") {
		websiteURL = "http://" + websiteURL
	}

	u, err := url.Parse(websiteURL)
	if err != nil {
		return ""
	}

	domain := u.Hostname()
	domain = strings.TrimPrefix(domain, "www.")

	return domain
}

// GetRelevantContacts: Busca e-mails filtrando por decisores (Executivos/Gerentes)
func (h *HunterClient) GetRelevantContacts(domain string, limit int) ([]HunterEmail, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain is empty")
	}

	endpoint := "https://api.hunter.io/v2/domain-search"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("domain", domain)
	q.Add("api_key", h.apiKey)
	q.Add("limit", fmt.Sprintf("%d", limit))

	// --- FILTRO DE RELEVÂNCIA ---
	// Aqui pedimos apenas cargos de liderança.
	// Se quiser tudo, remova essa linha.
	q.Add("seniority", "executive,manager")

	// Opcional: Se quiser filtrar por departamento (ex: TI e Engenharia)
	// q.Add("department", "it,engineering")
	req.URL.RawQuery = q.Encode()

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Tratamento de Erros da API
	if resp.StatusCode != http.StatusOK {
		// 401 = API Key inválida
		// 429 = Rate Limit excedido (muitas requisições)
		// 404 = Domínio não encontrado (não gasta crédito)
		if resp.StatusCode == 404 {
			return []HunterEmail{}, nil // Retorna vazio sem erro
		}
		return nil, fmt.Errorf("hunter api error: status %d", resp.StatusCode)
	}

	var result hunterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data.Emails, nil
}
