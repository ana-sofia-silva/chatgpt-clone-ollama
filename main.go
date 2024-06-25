package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}
}

func main() {
	r := setupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("Warning: PORT environment variable not set. Defaulting to 8080")
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	gracefulShutdown(server)

}

// setupRouter initializes and returns a new chi Mux
// router with configured routes and middleware.
func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.Get("/", index)
	r.Post("/run", run)
	return r
}

// index serves the main page.
func index(w http.ResponseWriter, r *http.Request) {
	templates, err := template.ParseFiles("static/index.html")
	if err != nil {
		log.Fatal("Error parsing template with error: ", err)
	}

	templates.Execute(w, nil)
}

// run handles the POST request to /run, decodes the input prompt,
// generates a response using an LLM, and sends back the result in JSON format.
func run(w http.ResponseWriter, r *http.Request) {

	var prompt Prompt

	if err := json.NewDecoder(r.Body).Decode(&prompt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	llm, err := ollama.New(ollama.WithModel("llama3"))
	if err != nil {
		log.Fatal("Failed to initialize LLM with error: ", err)
	}

	ctx := context.Background()

	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt.Input)
	if err != nil {
		log.Fatal("LLM generation failed with error: ", err)
	}

	response := PromptResponse{
		Input:    prompt.Input,
		Response: completion,
	}

	json.NewEncoder(w).Encode(response)
}

// gracefulShutdown listens for interrupt signals and gracefully
// shuts down the server, allowing ongoing operations to complete.
func gracefulShutdown(server *http.Server) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan // wait for SIGINT
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	log.Println("Server gracefully stopped")
}

// PromptResponse represents the structure for storing input from the user and the response from the LLM.
type PromptResponse struct {
	Input    string `json:"input"`
	Response string `json:"response"`
}

// Prompt represents the structure for storing input from the user.
type Prompt struct {
	Input string `json:"input"`
}
