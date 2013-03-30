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
	clientId   = flag.String("clientid", "", "OAuth Client ID")
	secret     = flag.String("secret", "", "OAuth Client Secret")
	cacheToken = flag.Bool("cachetoken", true, "cache the OAuth token")
	debug      = flag.Bool("debug", false, "show HTTP traffic")
)

var oauthClient *http.Client
var calendarApi *calendar.Service

func main() {
	// Parse the flags first
	flag.Parse()

	fmt.Println("-- vestige 1.0 ----------------------------")
	fmt.Println(" * Authenticating to Google...")

	// Set the Client ID and Secret
	oauthConfig.ClientId = *clientId
	oauthConfig.ClientSecret = *secret

	// Create the OAuth Client
	oauthClient = getOAuthClient(oauthConfig)

	// Initialize the Calendar API
	calendarApi, _ = calendar.New(oauthClient)

	// TODO: List calendars here
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

	_, err := calendarApi.Events.Insert("data@markbao.com", &eventNew).Do()

	return err
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
