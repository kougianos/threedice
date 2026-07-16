// Package integration exercises the whole service over HTTP against a real
// PostgreSQL, mirroring the Java suite's Testcontainers + REST Assured setup.
package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/kougianos/threedice/go/internal/app"
	"github.com/kougianos/threedice/go/internal/config"
)

// One container is shared across all test files, as the Java suite shares a
// single static PostgreSQLContainer.
var (
	srv    *httptest.Server
	pool   *pgxpool.Pool
	roller = &stubRoller{}

	// containerDSN points at containerDB; the migration tests derive URLs for
	// their own throwaway databases from it.
	containerDSN string
)

const containerDB = "threedice_test"

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx, "postgres:16-alpine",
		tcpostgres.WithDatabase(containerDB),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "starting postgres:", err)
		return 1
	}
	defer func() { _ = testcontainers.TerminateContainer(container) }()

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintln(os.Stderr, "connection string:", err)
		return 1
	}
	containerDSN = dsn

	// Quiet logs so test output stays readable.
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	a, err := app.New(ctx, config.Config{
		DatabaseURL:    dsn,
		InitialBalance: decimal.RequireFromString("1000.00"),
	}, roller, log)
	if err != nil {
		fmt.Fprintln(os.Stderr, "building app:", err)
		return 1
	}
	defer a.Close()

	pool = a.Pool
	srv = httptest.NewServer(a.Handler)
	defer srv.Close()

	return m.Run()
}

// ── Deterministic dice ──

// stubRoller replaces SecureRoller in tests, as @MockitoBean replaces
// DiceEngine on the Java side.
type stubRoller struct {
	mu     sync.Mutex
	values []int32
	idx    int
}

// Roll returns the scripted values in order and then repeats the last one,
// matching Mockito's thenReturn(a, b, c) behaviour on further calls.
func (s *stubRoller) Roll() int32 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.values) == 0 {
		return 1
	}
	v := s.values[s.idx]
	if s.idx < len(s.values)-1 {
		s.idx++
	}
	return v
}

// script sets the dice sequence, equivalent to when(diceEngine.roll()).thenReturn(...).
func script(values ...int32) {
	roller.mu.Lock()
	defer roller.mu.Unlock()
	roller.values = values
	roller.idx = 0
}

// ── Fixtures ──

// resetWithPlayer clears the tables and inserts one client, mirroring the
// @BeforeEach in each Java IT class. Returns the new player id.
func resetWithPlayer(t *testing.T, username, balance string) int64 {
	t.Helper()
	ctx := context.Background()

	_, err := pool.Exec(ctx, `TRUNCATE transactions, draws, bets, clients RESTART IDENTITY CASCADE`)
	require.NoError(t, err)

	var id int64
	err = pool.QueryRow(ctx,
		`INSERT INTO clients (username, balance) VALUES ($1, $2) RETURNING id`,
		username, decimal.RequireFromString(balance),
	).Scan(&id)
	require.NoError(t, err)
	return id
}

// ── HTTP helpers ──

type response struct {
	status int
	body   []byte
}

func (r response) json(t *testing.T) map[string]any {
	t.Helper()
	var m map[string]any
	require.NoError(t, json.Unmarshal(r.body, &m), "body: %s", r.body)
	return m
}

func (r response) array(t *testing.T) []map[string]any {
	t.Helper()
	var a []map[string]any
	require.NoError(t, json.Unmarshal(r.body, &a), "body: %s", r.body)
	return a
}

// detail returns the RFC 7807 detail field.
func (r response) detail(t *testing.T) string {
	t.Helper()
	d, _ := r.json(t)["detail"].(string)
	return d
}

func get(t *testing.T, path string) response {
	t.Helper()
	return do(t, http.MethodGet, path, "")
}

func postBet(t *testing.T, body string) response {
	t.Helper()
	return do(t, http.MethodPost, "/api/bets", body)
}

func do(t *testing.T, method, path, body string) response {
	t.Helper()

	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, srv.URL+path, reader)
	require.NoError(t, err)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := srv.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return response{status: resp.StatusCode, body: raw}
}

// betBody builds a place-bet payload with a unique idempotency key.
func betBody(playerID int64, stake, predicted string) string {
	return fmt.Sprintf(
		`{"playerId":%d,"stake":%s,"predictedValue":%s,"idempotencyKey":%q}`,
		playerID, stake, predicted, uniqueKey(),
	)
}

var keyCounter atomic.Int64

func uniqueKey() string {
	return fmt.Sprintf("test-key-%d-%d", time.Now().UnixNano(), keyCounter.Add(1))
}
