package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Configuration represents the configuration structure
type Configuration struct {
	PodmanURL     string
	ServerAddress string
	ServerPort    int
}

// Container represents a container object
type Container struct {
	Image string `json:"Image"`
	// other fields if needed
}

var (
	jobRunning = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "Nosana",
			Name:      "running_job",
			Help:      "Number of running jobs",
		})

	jobImage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "Nosana",
			Name:      "running_image",
			Help:      "Name of the image running",
		},
		[]string{
			"imageName",
		},
	)

	defaultPodmanURL     = "http://127.0.0.1:8080/v3.4.2/libpod/containers/json"
	defaultServerAddress string
	defaultServerPort    int
)

func init() {
	// Get the first IPv4 address of the machine
	defaultServerAddress, _ = getFirstIPv4()
	defaultServerPort = 8995
}

func collector(job prometheus.Gauge, jobImage *prometheus.GaugeVec, url string) {
	// Test URL response
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("[%s] - Error accessing URL: %s . Error: %s\n", time.Now().Format("2006-01-02 15:04:05"), url, err)
		job.Set(0)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[%s] - Podman URL %s returned non-OK status code: %d\n", time.Now().Format("2006-01-02 15:04:05"), url, resp.StatusCode)
		job.Set(0)
		return
	}

	// Read JSON data
	var containers []Container
	err = json.NewDecoder(resp.Body).Decode(&containers)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		job.Set(0)
		return
	}

	// Counter for containers with Nosana image
	nosanaContainers := 0

	// Loop through containers and check if "Image" contains "/nosana/"
	for _, container := range containers {
		if strings.Contains(container.Image, "/nosana/") {
			// Increment the counter
			nosanaContainers++
			jobImage.WithLabelValues(container.Image).Set(1)
			fmt.Printf("[%s] - Found container with Nosana image at  %s\n", time.Now().Format("2006-01-02 15:04:05"), container.Image)
		}
	}
	job.Set(float64(nosanaContainers))
}

func main() {
	podmanURLPtr := flag.String("podman-url", defaultPodmanURL, "URL for Podman API")
	serverAddressPtr := flag.String("server-address", defaultServerAddress, "Server address")
	serverPortPtr := flag.Int("server-port", defaultServerPort, "Server port")
	flag.Parse()

	config := Configuration{
		PodmanURL:     *podmanURLPtr,
		ServerAddress: *serverAddressPtr,
		ServerPort:    *serverPortPtr,
	}

	runApp(config)
}

func runApp(config Configuration) {
	jobRunning.Set(0) // Initialize the gauge

	prometheus.MustRegister(jobRunning)
	prometheus.MustRegister(jobImage)

	go func() {
		for {
			collector(jobRunning, jobImage, config.PodmanURL)
			time.Sleep(5 * time.Minute) // Adjust the interval as needed
		}
	}()

	http.Handle("/metrics", promhttp.Handler())

	listenAddress := fmt.Sprintf("%s:%d", config.ServerAddress, config.ServerPort)
	fmt.Printf("Exporter server: http://%s/metrics\n", listenAddress) // Print server info
	http.ListenAndServe(listenAddress, nil)
}

// getFirstIPv4 returns the first IPv4 address of the machine
func getFirstIPv4() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), nil
		}
	}
	return "", fmt.Errorf("no IPv4 address found")
}
