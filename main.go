package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

// Config is a struct describing the config parsed from cli arguments
type Config struct {
	PathToYAML string
	PathToJSON string
}

type URLMapper struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

// YAMLHandler will parse the provided YAML and then return an http.HandlerFunc (which also implements http.Handler) that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having invalid YAML data.
// See MapHandler to create a similar http.HandlerFunc via a mapping of paths to urls.

func YAMLHandler(YAML []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var mappers []URLMapper
	err := yaml.Unmarshal(YAML, &mappers)
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, r *http.Request) {
		for _, mapper := range mappers {
			if mapper.Path == r.URL.Path {
				http.Redirect(w, r, mapper.URL, 301)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}, nil
}

func JSONHandler(JSON []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var mappers []URLMapper
	err := json.Unmarshal(JSON, &mappers)
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, r *http.Request) {
		for _, mapper := range mappers {
			if mapper.Path == r.URL.Path {
				http.Redirect(w, r, mapper.URL, 301)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}, nil
}

// MapHandler will return an http.HandlerFunc (which also implements http.Handler) that will attempt to map any
// paths (keys in thge map) to their corresponding URL (values that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		originalURL, ok := pathsToUrls[r.URL.Path]
		if ok {
			http.Redirect(w, r, originalURL, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func main() {
	config := getConfig()

	yamlBytes := getFileBytes(config.PathToYAML)
	jsonBytes := getFileBytes(config.PathToJSON)

	mux := makeDefaultMux()
	mapHandler := makeMapHandler(mux)

	handler := mapHandler
	if yamlBytes != nil {
		handler = makeYAMLHandler(yamlBytes, &mapHandler)
	} else if jsonBytes != nil {
		handler = makeJSONHandler(jsonBytes, &mapHandler)
	}
	startServer(handler)
}

func getConfig() *Config {
	config := Config{}
	flag.StringVar(&config.PathToYAML, "yaml", "", "--yaml=path/to/file.yml")
	flag.StringVar(&config.PathToJSON, "json", "", "--json=path/to/file.json")
	flag.Parse()
	return &config
}

func getFileBytes(pathToFile string) []byte {
	bytes, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return nil
	}
	return bytes
}

func makeDefaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloWorldHandler)
	return mux
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func makeMapHandler(mux *http.ServeMux) http.HandlerFunc {
	return MapHandler(map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}, mux)
}

func makeYAMLHandler(yamlBytes []byte, fallbackHandler *http.HandlerFunc) http.HandlerFunc {
	handler, err := YAMLHandler(yamlBytes, fallbackHandler)
	if err != nil {
		panic(err)
	}
	return handler
}

func makeJSONHandler(jsonBytes []byte, fallbackHandler *http.HandlerFunc) http.HandlerFunc {
	handler, err := JSONHandler(jsonBytes, fallbackHandler)
	if err != nil {
		panic(err)
	}
	return handler
}

func startServer(handler http.HandlerFunc) {
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}
