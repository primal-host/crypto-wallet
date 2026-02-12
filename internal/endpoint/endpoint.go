package endpoint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Endpoint represents a named EVM RPC endpoint.
type Endpoint struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	Symbol string `json:"symbol"` // native token symbol (e.g. "AVAX", "ETH")
}

// Status is the live health info for an endpoint.
type Status struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Symbol      string `json:"symbol"`
	Online      bool   `json:"online"`
	ChainID     string `json:"chain_id,omitempty"`
	BlockNumber string `json:"block_number,omitempty"`
	Latency     int64  `json:"latency_ms"`
}

// Store manages endpoints loaded from a JSON file.
type Store struct {
	mu        sync.RWMutex
	endpoints []Endpoint
	path      string
}

// NewStore loads endpoints from a JSON file. If the file doesn't exist, starts empty.
func NewStore(path string) (*Store, error) {
	s := &Store{path: path}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			s.endpoints = []Endpoint{}
			return s, nil
		}
		return nil, fmt.Errorf("read endpoints: %w", err)
	}
	if err := json.Unmarshal(data, &s.endpoints); err != nil {
		return nil, fmt.Errorf("parse endpoints: %w", err)
	}
	return s, nil
}

// List returns all configured endpoints.
func (s *Store) List() []Endpoint {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Endpoint, len(s.endpoints))
	copy(out, s.endpoints)
	return out
}

var slugRe = regexp.MustCompile(`[^a-z0-9-]+`)

// slugify converts a name to a URL-safe ID.
func slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = strings.ReplaceAll(s, " ", "-")
	s = slugRe.ReplaceAllString(s, "")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "endpoint"
	}
	return s
}

// Add creates a new endpoint, generating an ID from the name.
func (s *Store) Add(ep Endpoint) (Endpoint, error) {
	if strings.TrimSpace(ep.Name) == "" {
		return Endpoint{}, fmt.Errorf("name is required")
	}
	if strings.TrimSpace(ep.URL) == "" {
		return Endpoint{}, fmt.Errorf("url is required")
	}
	if _, err := url.ParseRequestURI(ep.URL); err != nil {
		return Endpoint{}, fmt.Errorf("invalid url: %w", err)
	}
	if strings.TrimSpace(ep.Symbol) == "" {
		return Endpoint{}, fmt.Errorf("symbol is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate unique ID from name.
	base := slugify(ep.Name)
	id := base
	for n := 2; s.findLocked(id) != nil; n++ {
		id = fmt.Sprintf("%s-%d", base, n)
	}
	ep.ID = id

	s.endpoints = append(s.endpoints, ep)
	if err := s.save(); err != nil {
		// Roll back.
		s.endpoints = s.endpoints[:len(s.endpoints)-1]
		return Endpoint{}, err
	}
	return ep, nil
}

// Update replaces an existing endpoint's fields by ID.
func (s *Store) Update(id string, ep Endpoint) (Endpoint, error) {
	if strings.TrimSpace(ep.Name) == "" {
		return Endpoint{}, fmt.Errorf("name is required")
	}
	if strings.TrimSpace(ep.URL) == "" {
		return Endpoint{}, fmt.Errorf("url is required")
	}
	if _, err := url.ParseRequestURI(ep.URL); err != nil {
		return Endpoint{}, fmt.Errorf("invalid url: %w", err)
	}
	if strings.TrimSpace(ep.Symbol) == "" {
		return Endpoint{}, fmt.Errorf("symbol is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for i, existing := range s.endpoints {
		if existing.ID == id {
			ep.ID = id
			old := s.endpoints[i]
			s.endpoints[i] = ep
			if err := s.save(); err != nil {
				s.endpoints[i] = old
				return Endpoint{}, err
			}
			return ep, nil
		}
	}
	return Endpoint{}, fmt.Errorf("endpoint %q not found", id)
}

// Delete removes an endpoint by ID.
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, ep := range s.endpoints {
		if ep.ID == id {
			old := s.endpoints
			s.endpoints = append(s.endpoints[:i], s.endpoints[i+1:]...)
			if err := s.save(); err != nil {
				s.endpoints = old
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("endpoint %q not found", id)
}

// findLocked finds an endpoint by ID. Must be called with mu held.
func (s *Store) findLocked(id string) *Endpoint {
	for i := range s.endpoints {
		if s.endpoints[i].ID == id {
			return &s.endpoints[i]
		}
	}
	return nil
}

// save writes the current endpoints to disk. Must be called with mu held.
func (s *Store) save() error {
	data, err := json.MarshalIndent(s.endpoints, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal endpoints: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(s.path, data, 0644); err != nil {
		return fmt.Errorf("write endpoints: %w", err)
	}
	return nil
}

// Poll checks each endpoint with eth_chainId and eth_blockNumber, returning live status.
func (s *Store) Poll() []Status {
	eps := s.List()
	results := make([]Status, len(eps))
	var wg sync.WaitGroup
	for i, ep := range eps {
		wg.Add(1)
		go func(i int, ep Endpoint) {
			defer wg.Done()
			results[i] = poll(ep)
		}(i, ep)
	}
	wg.Wait()
	return results
}

func poll(ep Endpoint) Status {
	st := Status{
		ID:     ep.ID,
		Name:   ep.Name,
		URL:    ep.URL,
		Symbol: ep.Symbol,
	}

	start := time.Now()

	// Get chain ID.
	chainID, err := rpcCall(ep.URL, "eth_chainId", nil)
	if err != nil {
		st.Latency = time.Since(start).Milliseconds()
		return st
	}
	st.ChainID = chainID

	// Get block number.
	blockNum, err := rpcCall(ep.URL, "eth_blockNumber", nil)
	if err != nil {
		st.Latency = time.Since(start).Milliseconds()
		st.Online = true // chain ID worked, so it's partially online
		return st
	}
	st.BlockNumber = blockNum

	st.Latency = time.Since(start).Milliseconds()
	st.Online = true
	return st
}

// RPCCall makes a JSON-RPC call and returns the result string.
func RPCCall(url, method string, params []any) (json.RawMessage, error) {
	body := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, fmt.Errorf("rpc error %d: %s", result.Error.Code, result.Error.Message)
	}
	return result.Result, nil
}

// rpcCall is the internal helper returning a string result.
func rpcCall(url, method string, params []any) (string, error) {
	raw, err := RPCCall(url, method, params)
	if err != nil {
		return "", err
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return string(raw), nil
	}
	return s, nil
}
