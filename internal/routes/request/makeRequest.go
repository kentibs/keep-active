package request

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
)

func MakeRequest() {
	app := models.App{}
	fmt.Println("In the MakeRequest fn")

	apps, err := app.FindAll()
	if err != nil {
		fmt.Println("Error fetching apps:", err)
		return
	}

	// Run the MakeRequest function every 5 minutes
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		makeAllRequests(apps)
	}
}

func makeAllRequests(apps []models.App) {
	var wg sync.WaitGroup

	for _, app := range apps {
		wg.Add(1)
		go func(app models.App) {
			defer wg.Done()
			makeRequest(app)
		}(app)
	}

	wg.Wait()
}

func makeRequest(app models.App) {
	response, err := services.MakeHTTPRequest(app.URL)
	if err != nil {
		log.Println("Request error:", err)
		return
	}

	request := models.Request{
		AppID:      app.ID,
		StatusCode: response.StatusCode,
		Duration:   response.RequestTimeMS,
	}

	_, err = request.Create(request)
	if err != nil {
		log.Println("Error saving request:", err)
	}
}
