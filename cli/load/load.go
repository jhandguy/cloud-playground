package load

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jhandguy/devops-playground/cli/message"
)

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
	totalReqCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "devops_playground_cli_requests_total",
			Help: "Total requests counter per path and method",
		},
		[]string{"path", "method"},
	)

	successReqCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "devops_playground_cli_requests_success",
			Help: "Successful requests counter per path and method",
		},
		[]string{"path", "method"},
	)

	latencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "devops_playground_cli_requests_latency",
			Help: "Requests latency histogram per path and method",
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
	handleMissingFlag(viper.BindPFlag("gateway-api-key", Cmd.PersistentFlags().Lookup("token")))

	Cmd.PersistentFlags().StringP("url", "u", "", "gateway URL")
	handleMissingFlag(viper.BindPFlag("gateway-url", Cmd.PersistentFlags().Lookup("url")))

	Cmd.PersistentFlags().StringP("push", "p", "", "push gateway url")
	handleMissingFlag(viper.BindPFlag("pushgateway-url", Cmd.PersistentFlags().Lookup("push")))

	testLoadCmd.Flags().IntVarP(&rounds, "rounds", "r", 100, "number of test rounds")

	prometheus.MustRegister(totalReqCounter)
	prometheus.MustRegister(successReqCounter)
	prometheus.MustRegister(latencyHistogram)
}

func createMessage() (*message.Message, error) {
	startTime := time.Now()
	res, err := message.Create(message.Message{Content: "content"})

	totalReqCounter.
		WithLabelValues("message", "create").
		Inc()

	latencyHistogram.
		WithLabelValues("message", "create").
		Observe(time.Since(startTime).Seconds())

	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to create message: %d", res.StatusCode())
	}

	successReqCounter.
		WithLabelValues("message", "create").
		Inc()

	return res.Result().(*message.Message), nil
}

func getMessage(id string) (*message.Message, error) {
	startTime := time.Now()
	res, err := message.Get(id)

	totalReqCounter.
		WithLabelValues("message", "get").
		Inc()

	latencyHistogram.
		WithLabelValues("message", "get").
		Observe(time.Since(startTime).Seconds())

	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to get message: %d", res.StatusCode())
	}

	successReqCounter.
		WithLabelValues("message", "get").
		Inc()

	return res.Result().(*message.Message), nil
}

func deleteMessage(id string) error {
	startTime := time.Now()
	res, err := message.Delete(id)

	totalReqCounter.
		WithLabelValues("message", "delete").
		Inc()

	latencyHistogram.
		WithLabelValues("message", "delete").
		Observe(time.Since(startTime).Seconds())

	if err != nil {
		return err
	}

	if res.IsError() {
		return fmt.Errorf("failed to delete message: %d", res.StatusCode())
	}

	successReqCounter.
		WithLabelValues("message", "delete").
		Inc()

	return nil
}

func pushMetrics() {
	url := fmt.Sprintf("http://%s", viper.GetString("pushgateway-url"))
	if err := push.New(url, "cli").
		Collector(totalReqCounter).
		Collector(successReqCounter).
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

	for r := 0; r < rounds; r++ {
		wg.Add(1)

		go func() {
			msg, err := createMessage()
			if err != nil {
				failures++
				wg.Done()
				return
			}

			msg, err = getMessage(msg.ID)
			if err != nil {
				failures++
				wg.Done()
				return
			}

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
