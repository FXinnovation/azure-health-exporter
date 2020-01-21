package main

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"

	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

var (
	configFile     = kingpin.Flag("config.file", "Exporter configuration file.").Default("config/config.yml").String()
	listenAddress  = kingpin.Flag("web.listen-address", "The address to listen on for HTTP requests.").Default(":9613").String()
	metricsPath    = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	config         Config
	azureErrorDesc = prometheus.NewDesc("azure_error", "Error collecting metrics", nil, nil)
)

// Config of the exporter
type Config struct {
}

func init() {
	prometheus.MustRegister(version.NewCollector("azure_health_exporter"))
}

func main() {
	kingpin.Version(version.Print("azure-health-exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Info("Starting exporter", version.Info())
	log.Info("Build context", version.BuildContext())

	_, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	session, err := NewAzureSession(os.Getenv("AZURE_SUBSCRIPTION_ID"))
	if err != nil {
		log.Fatalf("Error creating Azure session: %v", err)
	}

	resourceHealthCollector := NewResourceHealthCollector(session)
	prometheus.MustRegister(resourceHealthCollector)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>azure-health-exporter</title></head>
			<body>
			<h1>azure-health-exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Info("Beginning to serve on address ", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func loadConfig(configFile string) (Config, error) {
	if fileExists(configFile) {
		log.Infof("Loading config file %v", configFile)

		// Load config from file
		configData, err := ioutil.ReadFile(configFile)
		if err != nil {
			return Config{}, err
		}

		return loadConfigContent(configData)
	}

	log.Infof("Config file %v does not exist, using default values", configFile)
	return Config{}, nil

}

func loadConfigContent(configData []byte) (Config, error) {
	config = Config{}
	var err error

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return config, err
	}

	log.Info("Config loaded")
	return config, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
