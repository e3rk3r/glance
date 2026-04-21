package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/glanceapp/glance/internal/config"
	"github.com/glanceapp/glance/internal/server"
)

const (
	defaultConfigPath = "glance.yml"
	version           = "0.1.0"
)

func main() {
	var (
		configPath  string
		showVersion bool
		showHelp    bool
	)

	flag.StringVar(&configPath, "config", defaultConfigPath, "Path to the configuration file")
	flag.StringVar(&configPath, "c", defaultConfigPath, "Path to the configuration file (shorthand)")
	flag.BoolVar(&showVersion, "version", false, "Print version and exit")
	flag.BoolVar(&showVersion, "v", false, "Print version and exit (shorthand)")
	flag.BoolVar(&showHelp, "help", false, "Print usage information and exit")
	flag.BoolVar(&showHelp, "h", false, "Print usage information and exit (shorthand)")
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if showVersion {
		fmt.Printf("glance v%s\n", version)
		os.Exit(0)
	}

	// Load configuration from the specified file
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration from %q: %v", configPath, err)
	}

	// Initialize and start the HTTP server
	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Log the config file path being used, helpful when running multiple instances
	log.Printf("Using config file: %s", configPath)

	// Include hostname in startup log to make it easier to identify which machine
	// is running this instance when checking logs remotely
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	log.Printf("Starting glance v%s on %s:%d (host: %s)", version, cfg.Server.Host, cfg.Server.Port, hostname)

	if err := srv.Start(); err != nil {
		log.Fatalf("Server encountered a fatal error: %v", err)
	}
}
