package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"code.google.com/p/goauth2/oauth"
	calendar "code.google.com/p/google-api-go-client/calendar/v3"
)

// Initialize OAuth configuration
var oauthConfig = &oauth.Config{
	ClientId:     "", // pass in --clientid
	ClientSecret: "", // pass in --secret
	Scope:        "https://www.googleapis.com/auth/calendar",
	AuthURL:      "https://accounts.google.com/o/oauth2/auth",
	TokenURL:     "https://accounts.google.com/o/oauth2/token",
}

// Define flags from command line
var (
	clientId        = flag.String("clientid", "", "OAuth Client ID")
	secret          = flag.String("secret", "", "OAuth Client Secret")
	cacheToken      = flag.Bool("cachetoken", true, "cache the OAuth token")
	debug           = flag.Bool("debug", false, "show HTTP traffic")
	singleCalendar  = flag.Bool("single", false, "only put events into a single calendar")
	primaryCalendar = flag.String("primary", "", "specify the default calendar for events to go in, by name")
)

// Make the OAuth Client and the Calendar API client available for all
var oauthClient *http.Client
var calendarApi *calendar.Service

// Initialize an empty calendar list for later caching of calendars
var calendarList = make(map[string]string)

// Initialize a string with the calendar name of the primary calendar
var primaryCalendarName string;

// Initialize a string with the calendar ID of the primary calendar
var primaryCalendarID string;

func main() {
	// Parse the flags first
	flag.Parse()

	fmt.Println("")
	fmt.Println("-- vestige 1.0 ----------------------------")

	// Display flags if necessary
	if *singleCalendar {
		fmt.Println(" * Loading in single calendar mode.")
	}

	fmt.Println(" * Authenticating to Google...")

	// Set the Client ID and Secret
	oauthConfig.ClientId = *clientId
	oauthConfig.ClientSecret = *secret

	// Create the OAuth Client
	oauthClient = getOAuthClient(oauthConfig)

	// Initialize the Calendar API
	calendarApi, _ = calendar.New(oauthClient)

	// Load all calendars and populate the calendarList variable.
	fmt.Println(" * Loading calendars...")
	loadCalendars()

	fmt.Println(" * Ready.")
	fmt.Println("")
	fmt.Println("")

	// Start the application loop.
	applicationLoop()
}

func applicationLoop() {
	for {
		var workItem string
		var scanDummy string

		fmt.Println("-- NEW WORK ITEM --------------------------")
		fmt.Println(" * What are you working on?")

		// Get the work item
		fmt.Print("   ")
		bufioReader := bufio.NewReader(os.Stdin)
		workItem, _ = bufioReader.ReadString('\n')

		// Save the current time
		startTime := time.Now()

		// Echo the starting time
		fmt.Println("")
		fmt.Println(" * Started at", startTime.Format(time.Kitchen))
		fmt.Print("   Hit Enter to finish work")

		// Wait for this scan to come in
		fmt.Scanln(&scanDummy)

		// All set now, record the end time
		endTime := time.Now()

		// Next, create the event
		fmt.Println("")
		fmt.Println(" * Sending to Google...")
		err := createEvent(workItem, startTime, endTime)

		if err == nil {
			fmt.Println(" * Sent.")
		} else {
			fmt.Println(" * An error occurred:")
			fmt.Println(err)
		}

		fmt.Println("-- END WORK ITEM --------------------------")
		fmt.Println("")
		fmt.Println("")
	}
}

func loadCalendars() {
	// List available calendars
	calendarListFromApi, err := calendarApi.CalendarList.List().MaxResults(50).MinAccessRole("writer").Do()

	if err != nil {
		fmt.Println(" * An error occurred:")
		fmt.Println(err)
		os.Exit(1)
	}

	// Calendars are now listed into calendarList
	// Let's get them out into the calendarList map
	for _, element := range calendarListFromApi.Items {
		// Assign the lowercased name to the ID of the calendar
		calendarList[strings.ToLower(element.Summary)] = element.Id

		// Check if it's the primary
		if element.Primary {
			primaryCalendarName = strings.ToLower(element.Summary)
			primaryCalendarID = strings.ToLower(element.Id)
		}
	}
}

// createEvent
// Creates an event with the parameters.
// 
// IN:  summary (string), startTime (Time), endTime (Time)
// OUT: error (error)

func createEvent(summary string, startTime time.Time, endTime time.Time) error {
	// create the event and return an err
	eventStart := calendar.EventDateTime{
		DateTime: startTime.Format(time.RFC3339),
	}

	eventEnd := calendar.EventDateTime{
		DateTime: endTime.Format(time.RFC3339),
	}

	eventNew := calendar.Event{
		Summary: summary,
		Start:   &eventStart,
		End:     &eventEnd,
	}

	// Category logic
	// Let's see if this summary has a spaced hyphen ( - ) in the middle
	splitSummary := strings.Split(summary, " - ")

	// Prepare a calendarForEvent string for the calendar name
	calendarIdForEvent := ""

	if *singleCalendar {
		// We are in single calendar mode.
		// No matter what, put it in one calendar
		// This would be the primary calendar
		calendarIdForEvent = primaryCalendarID
	} else {
		// Check splitSummary to see if it matches the original
		if splitSummary[0] == summary {
			// No hyphen, so we want to put this in the primary
			calendarIdForEvent = primaryCalendarID
		} else {
			// So, there is a hyphen, and we know what's on the left side now
			eventCategoryName := strings.ToLower(splitSummary[0])

			// See if this exists in the map...
			if calendarList[eventCategoryName] == "" {
				// This does not exist, so we have to create it
				calendarIdForEvent = createCalendar(splitSummary[0])
			} else {
				// It does exist, so let's use the resulting calendar ID
				calendarIdForEvent = calendarList[eventCategoryName]
			}
		}
	}

	_, err := calendarApi.Events.Insert(calendarIdForEvent, &eventNew).Do()

	return err
}

// createCalendar
// Create a calendar given a name, return the calendarId.
//
// IN:  name (string)
// OUT: calendarId (string)

func createCalendar(name string) string {
	// Initialize a Calendar struct
	calendarNew := calendar.Calendar {
		Summary: name,
	}

	// Insert it into the Calendar API
	calendarData, err := calendarApi.Calendars.Insert(&calendarNew).Do()

	if err != nil {
		fmt.Println(" * An error occurred:")
		fmt.Println(err)
	} else {
		fmt.Println(" * Calendar created:", name)
	}

	// Append this to the calendarList
	calendarList[strings.ToLower(name)] = calendarData.Id

	return calendarData.Id
}



// Google API

func osUserCacheDir() string {
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Caches")
	case "linux", "freebsd":
		return filepath.Join(os.Getenv("HOME"), ".cache")
	}
	log.Printf("TODO: osUserCacheDir on GOOS %q", runtime.GOOS)
	return "."
}

func tokenCacheFile(config *oauth.Config) string {
	return filepath.Join(osUserCacheDir(), url.QueryEscape(
		fmt.Sprintf("go-api-demo-%s-%s-%s", config.ClientId, config.ClientSecret, config.Scope)))
}

func tokenFromFile(file string) (*oauth.Token, error) {
	if !*cacheToken {
		return nil, errors.New("--cachetoken is false")
	}
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := new(oauth.Token)
	err = gob.NewDecoder(f).Decode(t)
	return t, err
}

func saveToken(file string, token *oauth.Token) {
	f, err := os.Create(file)
	if err != nil {
		log.Printf("Warning: failed to cache oauth token: %v", err)
		return
	}
	defer f.Close()
	gob.NewEncoder(f).Encode(token)
}

func condDebugTransport(rt http.RoundTripper) http.RoundTripper {
	return rt
}

func getOAuthClient(config *oauth.Config) *http.Client {
	cacheFile := tokenCacheFile(config)
	token, err := tokenFromFile(cacheFile)
	if err != nil {
		token = tokenFromWeb(config)
		saveToken(cacheFile, token)
	} else {
		// log.Printf("Using cached token %#v from %q", token, cacheFile)
	}

	t := &oauth.Transport{
		Token:     token,
		Config:    config,
		Transport: condDebugTransport(http.DefaultTransport),
	}
	return t.Client()
}

func tokenFromWeb(config *oauth.Config) *oauth.Token {
	ch := make(chan string)
	randState := fmt.Sprintf("st%d", time.Now())
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/favicon.ico" {
			http.Error(rw, "", 404)
			return
		}
		if req.FormValue("state") != randState {
			log.Printf("State doesn't match: req = %#v", req)
			http.Error(rw, "", 500)
			return
		}
		if code := req.FormValue("code"); code != "" {
			fmt.Fprintf(rw, "<h1>Success</h1>Authorized.")
			rw.(http.Flusher).Flush()
			ch <- code
			return
		}
		log.Printf("no code")
		http.Error(rw, "", 500)
	}))
	defer ts.Close()

	config.RedirectURL = ts.URL
	authUrl := config.AuthCodeURL(randState)
	go openUrl(authUrl)
	log.Printf("Authorize this app at: %s", authUrl)
	code := <-ch
	log.Printf("Got code: %s", code)

	t := &oauth.Transport{
		Config:    config,
		Transport: condDebugTransport(http.DefaultTransport),
	}
	_, err := t.Exchange(code)
	if err != nil {
		log.Fatalf("Token exchange error: %v", err)
	}
	return t.Token
}

func openUrl(url string) {
	try := []string{"xdg-open", "google-chrome", "open"}
	for _, bin := range try {
		err := exec.Command(bin, url).Run()
		if err == nil {
			return
		}
	}
	log.Printf("Error opening URL in browser.")
}

func valueOrFileContents(value string, filename string) string {
	if value != "" {
		return value
	}
	slurp, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading %q: %v", filename, err)
	}
	return strings.TrimSpace(string(slurp))
}
