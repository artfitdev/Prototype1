package main


import (
    "fmt"
    "net/http"
    //"encoding/json"
    //"github.com/google/go-github/github"
    "io/ioutil"
    "golang.org/x/oauth2"
    //githuboauth "golang.org/x/oauth2/github"
    "github.com/gorilla/sessions"
)

var (
    // You must register the app at https://github.com/settings/applications
    // Set callback to http://127.0.0.1:7000/github_oauth_cb
    // Set ClientId and ClientSecret to
    FitBitEndpoint = oauth2.Endpoint{
        AuthURL:  "https://www.fitbit.com/oauth2/authorize",
        TokenURL: "https://api.fitbit.com/oauth2/token",
    }
    oauthConf = &oauth2.Config{
        ClientID:     "229W98",
        ClientSecret: "8eec31295b71c172c01260bff308fc46",
        // select level of access you want https://developer.github.com/v3/oauth/#scopes
        Scopes:       []string{"activity", "heartrate"},
        Endpoint:     FitBitEndpoint,
    }
    
    token *oauth2.Token

    authenticated = false    
    // random string for oauth2 API calls to protect against CSRF
    oauthStateString = "thisshouldberandom"

    //SessionStore = sessions.NewCookieStore([]byte("FitBitToken"))
)
    const htmlIndex = `<html><body>
    Logged in with <a href="/login">FitBit</a>
    </body></html>
    `


  


    // /
    func handleMain(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(htmlIndex))

        if (token!= nil) {

            fmt.Printf("Token '%s'\n", token.AccessToken)
            
            oauthClient := oauthConf.Client(oauth2.NoContext, token)

            resp, err := oauthClient.Get("https://api.fitbit.com/1/user/-/activities/heart/date/today/1d.json")
            if err != nil {
                fmt.Printf("http.Conf failed '%s'\n", err)
                http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
                return
            }

            body, err := ioutil.ReadAll(resp.Body)
            resp.Body.Close()
            if err != nil {
               fmt.Printf("Error '%s'\n", err)
            }
            w.Write(body)
        }

        
        
    }

    func handleFitBitLogin(w http.ResponseWriter, r *http.Request) {
        url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
        http.Redirect(w, r, url, http.StatusTemporaryRedirect)
    }

    // /github_oauth_cb. Called by fitbit after authorization is granted
    func handleGitFitBitCallback(w http.ResponseWriter, r *http.Request) {
        state := r.FormValue("state")
        if state != oauthStateString {
            fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
            http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
            return
        }

        code := r.FormValue("code")
        token, err := oauthConf.Exchange(oauth2.NoContext, code)
        if err != nil {
            fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
            http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
            return
        }

        fmt.Printf("accessToken >> %s\n", token.AccessToken )
        



        if err != nil {
            fmt.Printf("client.Users.Get() failed with '%s'\n", err)
            http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
            return
        }
        fmt.Printf("Authorized")     



        authenticated = true
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
    }


func main() {
    http.HandleFunc("/", handleMain)
    http.HandleFunc("/login", handleFitBitLogin)
    http.HandleFunc("/oauth/fitbit/callback", handleGitFitBitCallback)

   //Request URL: http://www.fotolite.net/oauth/fitbit/callback?state=thisshouldberandom&code=608c2a7a99105b4bf6ac3ab13a9d99cb6b9e00a6
    fmt.Print("Started running on http://127.0.0.1:80\n")
    fmt.Println(http.ListenAndServe(":80", nil))


}
