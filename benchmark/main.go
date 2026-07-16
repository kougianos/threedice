// Command benchmark load-tests one ThreeDice endpoint and reports latency and
// throughput. It backs the numbers in the root README, so they can be re-checked
// rather than taken on trust.
//
// It exists because the bet endpoint requires a unique idempotency key per
// request, which off-the-shelf load generators cannot produce -- a fixed body
// would be rejected as a duplicate after the first request.
//
// Point it at either service:
//
//	go run . -url http://localhost:8080 -mode get -n 5000 -c 50   # Spring Boot
//	go run . -url http://localhost:8081 -mode bet -n 1000 -c 10   # Go
//
// Warm the JVM before measuring it, or you are timing JIT compilation rather
// than the code: run once with -quiet, then again for real.
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	base := flag.String("url", "http://localhost:8080", "base url of the service")
	n := flag.Int("n", 1000, "total requests")
	c := flag.Int("c", 50, "concurrent workers")
	mode := flag.String("mode", "get", "get (GET /api/players/1) | bet (POST /api/bets)")
	label := flag.String("label", "run", "label, keeps idempotency keys unique across runs")
	quiet := flag.Bool("quiet", false, "run without reporting, for warmup passes")
	flag.Parse()

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 200,
		},
	}

	var (
		counter  atomic.Int64
		failures atomic.Int64
		mu       sync.Mutex
		latency  = make([]time.Duration, 0, *n)
		statuses = map[int]int{}
	)

	jobs := make(chan struct{}, *n)
	for i := 0; i < *n; i++ {
		jobs <- struct{}{}
	}
	close(jobs)

	start := time.Now()
	var wg sync.WaitGroup
	for w := 0; w < *c; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				t0 := time.Now()
				code, err := do(client, *base, *mode, *label, counter.Add(1))
				elapsed := time.Since(t0)
				if err != nil {
					failures.Add(1)
					continue
				}
				mu.Lock()
				latency = append(latency, elapsed)
				statuses[code]++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if !*quiet {
		report(time.Since(start), latency, statuses, failures.Load(), *n, *c)
	}
}

func do(client *http.Client, base, mode, label string, seq int64) (int, error) {
	var (
		resp *http.Response
		err  error
	)
	switch mode {
	case "get":
		resp, err = client.Get(base + "/api/players/1")
	case "bet":
		body := fmt.Sprintf(
			`{"playerId":1,"stake":1,"predictedValue":12,"idempotencyKey":"%s-%d-%d"}`,
			label, time.Now().UnixNano(), seq)
		resp, err = client.Post(base+"/api/bets", "application/json", strings.NewReader(body))
	default:
		fmt.Fprintln(os.Stderr, "unknown mode:", mode)
		os.Exit(2)
	}
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode, nil
}

func report(elapsed time.Duration, latency []time.Duration, statuses map[int]int, failures int64, n, c int) {
	if len(latency) == 0 {
		fmt.Println("  no successful responses")
		return
	}
	sort.Slice(latency, func(i, j int) bool { return latency[i] < latency[j] })

	pct := func(p float64) time.Duration {
		return latency[int(float64(len(latency)-1)*p)]
	}
	var total time.Duration
	for _, d := range latency {
		total += d
	}

	codes := make([]string, 0, len(statuses))
	for code, count := range statuses {
		codes = append(codes, fmt.Sprintf("%d:%d", code, count))
	}
	sort.Strings(codes)

	fmt.Printf("  requests=%d concurrency=%d elapsed=%.2fs\n", n, c, elapsed.Seconds())
	fmt.Printf("  throughput=%.0f req/s\n", float64(len(latency))/elapsed.Seconds())
	fmt.Printf("  latency mean=%.2fms p50=%.2fms p95=%.2fms p99=%.2fms max=%.2fms\n",
		ms(total/time.Duration(len(latency))), ms(pct(0.50)), ms(pct(0.95)),
		ms(pct(0.99)), ms(latency[len(latency)-1]))
	fmt.Printf("  statuses=%s failures=%d\n", strings.Join(codes, " "), failures)
}

func ms(d time.Duration) float64 { return float64(d.Nanoseconds()) / 1e6 }
