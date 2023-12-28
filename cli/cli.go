package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"github.com/ssd0247/static-site-generator/config"
	"github.com/ssd0247/static-site-generator/datasource"
	"github.com/ssd0247/static-site-generator/generator"
)

// Run runs the application
func Run() {
	cfg, err := readConfig()
	if err != nil {
		log.Fatal("There was an error while reading the configurationn file: ", err)
	}
	ds := datasource.New()
	dirs, err := ds.Fetch(cfg)
	if err != nil {
		log.Fatal(err)
	}

	g := generator.New(&generator.SiteConfig{
		Sources:     dirs,
		Destination: cfg.Generator.Dest,
		Config:      cfg,
	})
	err = g.Generate()

	if err != nil {
		log.Fatal(err)
	}
}

// Start a local HTTP server for development/testing purposes
func Serve() {
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		log.Fatal(cwdErr)
	}

	staticPath := filepath.Join(cwd, "www")
	log.Printf("Serving directory %s", staticPath)
	log.Println("HTTP Server listening on 127.0.0.1:8000")
	fileSystemDir := http.FileServer(http.Dir(staticPath))
	http.Handle("/", fileSystemDir)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func readConfig() (*config.Config, error) {
	data, err := ioutil.ReadFile("bloggen.yml")
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %v", err)
	}
	cfg := config.Config{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("could not parse config: %v", err)
	}
	if cfg.Generator.Repo == "" {
		return nil, fmt.Errorf("please provide a repository URL, e.g.: https://github.com/ssd0247/blog")
	}
	if cfg.Generator.Tmp == "" {
		cfg.Generator.Tmp = "tmp"
	}
	if cfg.Generator.Dest == "" {
		cfg.Generator.Dest = "www"
	}
	if cfg.Blog.URL == "" {
		return nil, fmt.Errorf("please provide a Blog URL, e.g.: https://www.ssd0247.github.io/blog")
	}
	if cfg.Blog.Language == "" {
		cfg.Blog.Language = "en-us"
	}
	if cfg.Blog.Description == "" {
		return nil, fmt.Errorf("please provide a Blog Description, e.g.: A blog about Statistics, First-Order Logic, Programming, Software Development, and all that brings ML/AI to life")
	}
	if cfg.Blog.Dateformat == "" {
		cfg.Blog.Dateformat = "14.07.2023"
	}
	if cfg.Blog.Title == "" {
		return nil, fmt.Errorf("please provide a Blog Title, e.g.: What does 'support-vectors' mean in SVM?")
	}
	if cfg.Blog.Author == "" {
		return nil, fmt.Errorf("please provide a Blog author, e.g.: Shubham Dhapola")
	}
	if cfg.Blog.Frontpageposts == 0 {
		cfg.Blog.Frontpageposts = 10
	}
	return &cfg, nil
}
