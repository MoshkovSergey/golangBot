package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gosimple/slug"
)

// main is the entry point of the program.
//
// It takes a URL as an argument, and optionally a file name as a second argument.
// If no file name is provided, it generates a file name from the URL.
// The generated file name is printed to the console.
func main() {
	// Define the variables
	var fileName string // File name for the screenshot
	var imageBuf []byte // Buffer to store the screenshot

	// Create a new context with cancel function
	ctx, cancel := chromedp.NewContext(context.Background())

	// Make sure to cancel the context at the end
	defer cancel()

	// Get the URL from the command line arguments
	url := os.Args[1]

	// If a file name is provided as a second argument, use it
	if len(os.Args) == 3 {
		fileName = os.Args[2] // Use the provided file name
	} else {
		// Otherwise, generate a file name from the URL
		fileName = slug.Make(url) + ".jpg" // Generate a file name from the URL
	}

	// Run the Chrome DevTools Protocol (CDP) tasks
	// The CDP tasks are defined in the screenShotTask function
	// The tasks capture a screenshot of the URL and save it to the buffer
	// The buffer is passed as a pointer to the function
	if err := chromedp.Run(
		ctx, // Context: The context to run the tasks in
		screenShotTask(url, fileName, &imageBuf), // Task: The task to run
	); err != nil {
		fmt.Println(err) // Panic if an error occurs during the execution of the tasks
	}

	if err := os.WriteFile(fileName, imageBuf, 0644); err != nil {
		fmt.Println(err) // Panic if an error occurs during the execution of the tasks
	}
}

// screenShotTask is a function that captures a screenshot of a given URL and
// saves it to the provided buffer. It returns a list of chromedp.Tasks.
//
// Parameters:
// - url (string): The URL to capture a screenshot of.
// - fileName (string): The file name to save the screenshot as.
// - imageBuf (*[]byte): A pointer to the buffer where the screenshot will be saved.
//
// Returns:
// - chromedp.Tasks: A list of chromedp.Tasks that perform the screenshot capture.
func screenShotTask(url string, _ string, imageBuf *[]byte) chromedp.Tasks {
	// ActionFunc captures a screenshot and saves it to the buffer
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(2 * time.Second),
		chromedp.ActionFunc(func(ctx context.Context) (err error) {
			// Capture a screenshot of the URL and save it to the buffer
			*imageBuf, err = page.CaptureScreenshot().Do(ctx)
			return err
		}),
	}
}
