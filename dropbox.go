package dropbox
import(
	"net/http"
	"strings"
	"net/url"
	"github.com/Esseh/goauth"
)

var Config struct {
	Redirect string
	ClientID string
	SecretID string
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	UID         string `json:"uid"`
	State		string
}

type DropboxAccountInfo struct {
	UID int `json:"uid"`
	DisplayName string `json:"display_name"`
	NameDetails struct {
		FamiliarName string `json:"familiar_name"`
		GivenName string `json:"given_name"`
		Surname string `json:"surname"`
	} `json:"name_details"`
	ReferralLink string `json:"referral_link"`
	Country string `json:"country"`
	Locale string `json:"locale"`
	Email string `json:"email"`
	EmailVerified bool `json:"email_verified"`
	IsPaired bool `json:"is_paired"`
	Team struct {
		Name string `json:"name"`
		TeamID string `json:"team_id"`
	} `json:"team"`
	QuotaInfo struct {
		Shared int64 `json:"shared"`
		Quota int64 `json:"quota"`
		Normal int64 `json:"normal"`
	} `json:"quota_info"`
}

func (d Token)AccountInfo(req *http.Request)(DropboxAccountInfo , error){
	ai := DropboxAccountInfo{}
	values := make(url.Values)
	values.Add("access_token",d.AccessToken)
	err := goauth.CallAPI(req,"GET", "https://api.dropboxapi.com/1/account/info", values, &ai)	
	return ai,err
}

//////////////////////////////////////////////////////////////////////////////////
// Send for Dropbox OAuth
//////////////////////////////////////////////////////////////////////////////////
func Send(res http.ResponseWriter, req *http.Request){
	values := goauth.RequiredSend(res,req,Config.Redirect,Config.ClientID)
	http.Redirect(res, req, "https://www.dropbox.com/1/oauth2/authorize?"+values.Encode(), http.StatusSeeOther)
}

//////////////////////////////////////////////////////////////////////////////////
// Recieve for Dropbox OAuth
//////////////////////////////////////////////////////////////////////////////////
func Recieve(res http.ResponseWriter,req *http.Request) Token {
	token := Token{}
	resp, err := goauth.RequiredRecieve(res,req,Config.ClientID,Config.SecretID,Config.Redirect,"https://api.dropbox.com/1/oauth2/token") 
	if err != nil { return Token{} }
	
	err = goauth.ExtractValue(resp,&token)
	if err != nil { return Token{} }
	
	token.State = strings.Split(req.FormValue("state"),"](|)[")[1]
	return token
}	