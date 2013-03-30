package main

import (
  "encoding/gob"
  "io/ioutil"
	"fmt"
  "flag"
  "net/http"
  "net/http/httptest"
  "errors"
  "log"
  "net/url"
  "os"
  "os/exec"
  "path/filepath"
  "runtime"
  "strings"
  "time"

	"code.google.com/p/goauth2/oauth"
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

func main() {
  // Parse the flags first
  flag.Parse()

  // Set the Client ID and Secret
  oauthConfig.ClientId = *clientId;
  oauthConfig.ClientSecret = *secret;

  oauthClient := getOAuthClient(oauthConfig);

  fmt.Println(oauthClient)
}

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
                log.Printf("Using cached token %#v from %q", token, cacheFile)
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