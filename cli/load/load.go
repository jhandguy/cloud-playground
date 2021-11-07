package load

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jhandguy/devops-playground/cli/message"
)

// Cmd load command
var Cmd = &cobra.Command{
	Use:   "load",
	Short: "Load sub commands",
	Long:  "Load sub commands to test system load",
}

var testLoadCmd = &cobra.Command{
	Use:   "test",
	Short: "Test load",
	Long:  "Test load and push metrics to Prometheus",
	Run:   testLoad,
}

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "devops_playground_cli_requests_count",
			Help: "Request counter per path, method and deployment",
		},
		[]string{"path", "method", "deployment", "success"},
	)

	latencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "devops_playground_cli_requests_latency",
			Help: "Request latency histogram per path and method",
		},
		[]string{"path", "method"},
	)
)

var (
	rounds int
)

func handleMissingFlag(err error) {
	if err != nil {
		log.Fatalf("missing required flag: %v", err)
	}
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	Cmd.AddCommand(testLoadCmd)

	Cmd.PersistentFlags().StringP("token", "t", "", "gateway auth token")
	handleMissingFlag(viper.BindPFlag("gateway-token", Cmd.PersistentFlags().Lookup("token")))

	Cmd.PersistentFlags().StringP("url", "u", "", "gateway URL")
	handleMissingFlag(viper.BindPFlag("gateway-url", Cmd.PersistentFlags().Lookup("url")))

	Cmd.PersistentFlags().StringP("host", "o", "", "gateway host")
	handleMissingFlag(viper.BindPFlag("gateway-host", Cmd.PersistentFlags().Lookup("host")))

	Cmd.PersistentFlags().StringP("push", "p", "", "push gateway url")
	handleMissingFlag(viper.BindPFlag("pushgateway-url", Cmd.PersistentFlags().Lookup("push")))

	testLoadCmd.Flags().IntVarP(&rounds, "rounds", "r", 100, "number of test rounds")

	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(latencyHistogram)
}

func createMessage() (*message.Message, error) {
	startTime := time.Now()
	res, err := message.Create(message.Message{Content: "content"})

	latencyHistogram.
		WithLabelValues("message", "create").
		Observe(time.Since(startTime).Seconds())

	if err != nil {
		return nil, err
	}
	deployment := res.Header().Get("x-debug")

	if res.IsError() {
		requestCounter.
			WithLabelValues("message", "create", deployment, "false").
			Inc()

		return nil, fmt.Errorf("failed to create message: %d", res.StatusCode())
	}

	requestCounter.
		WithLabelValues("message", "create", deployment, "true").
		Inc()

	return res.Result().(*message.Message), nil
}

func getMessage(id string) (*message.Message, error) {
	startTime := time.Now()
	res, err := message.Get(id)

	latencyHistogram.
		WithLabelValues("message", "get").
		Observe(time.Since(startTime).Seconds())

	if err != nil {
		return nil, err
	}
	deployment := res.Header().Get("x-debug")

	if res.IsError() {
		requestCounter.
			WithLabelValues("message", "get", deployment, "false").
			Inc()

		return nil, fmt.Errorf("failed to get message: %d", res.StatusCode())
	}

	requestCounter.
		WithLabelValues("message", "get", deployment, "true").
		Inc()

	return res.Result().(*message.Message), nil
}

func deleteMessage(id string) error {
	startTime := time.Now()
	res, err := message.Delete(id)

	latencyHistogram.
		WithLabelValues("message", "delete").
		Observe(time.Since(startTime).Seconds())

	if err != nil {
		return err
	}
	deployment := res.Header().Get("x-debug")

	if res.IsError() {
		requestCounter.
			WithLabelValues("message", "delete", deployment, "false").
			Inc()

		return fmt.Errorf("failed to delete message: %d", res.StatusCode())
	}

	requestCounter.
		WithLabelValues("message", "delete", deployment, "true").
		Inc()

	return nil
}

func pushMetrics() {
	url := fmt.Sprintf("http://%s", viper.GetString("pushgateway-url"))
	if err := push.New(url, "cli").
		Collector(requestCounter).
		Collector(latencyHistogram).
		Push(); err != nil {
		log.Fatal(err)
	}
}

func testLoad(cmd *cobra.Command, _ []string) {
	var wg sync.WaitGroup
	failures := 0

	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Starting load test with %d rounds\n", rounds); err != nil {
		log.Fatal(err)
	}

	sleep := func(sec int) {
		duration := time.Duration(rand.Intn(sec))
		time.Sleep(duration * time.Second)
	}

	for r := 0; r < rounds; r++ {
		wg.Add(1)

		go func() {
			sleep(40)

			msg, err := createMessage()
			if err != nil {
				failures++
				wg.Done()
				return
			}

			sleep(10)

			msg, err = getMessage(msg.ID)
			if err != nil {
				failures++
				wg.Done()
				return
			}

			sleep(10)

			err = deleteMessage(msg.ID)
			if err != nil {
				failures++
				wg.Done()
				return
			}

			wg.Done()
		}()
	}

	wg.Wait()

	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Finished load test with %d failures\n", failures); err != nil {
		log.Fatal(err)
	}

	pushMetrics()
}
