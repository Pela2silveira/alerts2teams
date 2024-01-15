package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"
)

// Define a struct to hold the data
type AlertData struct {
	Status            string
	CommonLabels      map[string]string
	CommonAnnotations map[string]string
	ExternalURL       string
	Alerts            []Alert
}

type Alert struct {
	Annotations map[string]string
	Labels      map[string]string
}

// RenderJSONTemplate renders the provided JSON data using the template
func RenderJSONTemplate(data AlertData, templateFilePath string) (string, error) {
	// Load the template
	tmpl, err := template.New("teams.card").ParseFiles(templateFilePath)
	if err != nil {
		return "", fmt.Errorf("error reading template file: %v", err)
	}

	// Render the template with the data
	var renderedTemplate bytes.Buffer
	err = tmpl.ExecuteTemplate(&renderedTemplate, "teams.card", data)
	if err != nil {
		return "", fmt.Errorf("error rendering template: %v", err)
	}

	return renderedTemplate.String(), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Only allow application/json content type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Unsupported Media Type. Only application/json is accepted", http.StatusUnsupportedMediaType)
		return
	}

	// Decode the JSON request body
	var data AlertData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Render the template
	renderedTemplate, err := RenderJSONTemplate(data, "template.tmpl")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		return
	}

	// Output the rendered template
	fmt.Println("Rendered Template:\n", renderedTemplate)

	// Send the rendered JSON to another endpoint using a POST request
	destinationEndpoint := os.Getenv("DESTINATION_ENDPOINT")


	
	resp, err := http.Post(destinationEndpoint, "application/json", bytes.NewBufferString(renderedTemplate))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending POST request to destination: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Template rendered and sent successfully"))
}

func main() {
	// Read environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// HTTP server configuration
	http.HandleFunc("/", handler)
	serverAddr := fmt.Sprintf(":%s", port)

	// Start the HTTP server
	fmt.Printf("Server listening on port %s...\n", port)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
