package load

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/jhandguy/cloud-playground/cli/message"
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
			Name: "cloud_playground_cli_requests_count",
			Help: "Request counter per path and method",
		},
		[]string{"path", "method", "success"},
	)

	latencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "cloud_playground_cli_requests_latency",
			Help: "Request latency histogram per path and method",
		},
		[]string{"path", "method"},
	)
)

var (
	rounds int
)

func handleUnbindableFlag(err error) {
	if err != nil {
		zap.S().Fatalw("could not bind flag", "error", err.Error())
	}
}

func init() {
	Cmd.AddCommand(testLoadCmd)

	Cmd.PersistentFlags().StringP("push", "p", "", "push gateway url")
	handleUnbindableFlag(viper.BindPFlag("pushgateway-url", Cmd.PersistentFlags().Lookup("push")))

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
		requestCounter.
			WithLabelValues("message", "create", "false").
			Inc()

		return nil, err
	}

	requestCounter.
		WithLabelValues("message", "create", "true").
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
		requestCounter.
			WithLabelValues("message", "get", "false").
			Inc()

		return nil, err
	}

	requestCounter.
		WithLabelValues("message", "get", "true").
		Inc()

	return res.Result().(*message.Message), nil
}

func deleteMessage(id string) error {
	startTime := time.Now()
	_, err := message.Delete(id)

	latencyHistogram.
		WithLabelValues("message", "delete").
		Observe(time.Since(startTime).Seconds())

	if err != nil {
		requestCounter.
			WithLabelValues("message", "delete", "false").
			Inc()

		return err
	}

	requestCounter.
		WithLabelValues("message", "delete", "true").
		Inc()

	return nil
}

func pushMetrics() {
	url := fmt.Sprintf("http://%s", viper.GetString("pushgateway-url"))
	if err := push.New(url, "cli").
		Collector(requestCounter).
		Collector(latencyHistogram).
		Push(); err != nil {
		zap.S().Errorw("failed to push metrics", "error", err.Error())
	}
}

func testLoad(cmd *cobra.Command, _ []string) {
	zap.S().Infow("starting load test", "rounds", rounds)

	sleep := func(sec int) {
		duration := time.Duration(rand.Intn(sec))
		time.Sleep(duration * time.Second)
	}

	failures := 0
	channel := make(chan error, rounds)

	for r := 0; r < rounds; r++ {
		go func() {
			sleep(40)

			msg, err := createMessage()
			if err != nil {
				channel <- err
				return
			}

			sleep(10)

			msg, err = getMessage(msg.ID)
			if err != nil {
				channel <- err
				return
			}

			sleep(10)

			err = deleteMessage(msg.ID)
			channel <- err
		}()
	}

	for i := 0; i < rounds; i++ {
		err := <-channel
		if err != nil {
			failures++
		}
	}

	zap.S().Infow("finished load test", "failures", failures)

	pushMetrics()
}
