package oauth2

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/phuc0302/go-oauth2/mongo"
	"github.com/phuc0302/go-oauth2/util"
)

// TestUnit describes an implementation for OAuth2 unit test.
type TestUnit struct {
	Session  *mgo.Session
	Database *mgo.Database

	Client IClient
	User1  IUser
	User2  IUser
	//	ExpiredAccessToken  IToken
	//	ExpiredRefreshToken IToken

	Username     string
	Password     string
	UserID       bson.ObjectId
	ClientID     bson.ObjectId
	ClientSecret bson.ObjectId
	CreatedTime  time.Time
}

// Setup initializes environment.
func (u *TestUnit) Setup() {
	mongo.ConnectMongo()

	u.Session, u.Database = mongo.GetMonotonicSession()
	u.Username = "admin"
	u.Password = "admin"
	u.UserID = bson.NewObjectId()
	u.ClientID = bson.NewObjectId()
	u.ClientSecret = bson.NewObjectId()
	u.CreatedTime, _ = time.Parse(time.RFC822, "02 Jan 06 15:04 MST")

	// Define global variables
	Cfg = loadConfig(debug)
	objectFactory = &DefaultFactory{}
	TokenStore = objectFactory.CreateStore()

	// Generate test data
	u.Client = &DefaultClient{
		ID:     u.ClientID,
		Secret: u.ClientSecret,
		Grants: []string{AuthorizationCodeGrant, PasswordGrant, RefreshTokenGrant},

		Redirects: []string{"http://www.sample01.com", "http://www.sample02.com"},
	}

	password1, _ := util.EncryptPassword(u.Password)
	u.User1 = &DefaultUser{
		ID:    u.UserID,
		User:  u.Username,
		Pass:  password1,
		Roles: []string{"r_user", "r_admin"},
	}

	password2, _ := util.EncryptPassword(u.ClientSecret.Hex())
	u.User2 = &DefaultUser{
		ID:    u.ClientID,
		User:  u.ClientID.Hex(),
		Pass:  password2,
		Roles: []string{"r_device"},
	}

	//	u.ExpiredAccessToken = &DefaultToken{
	//		ID:      bson.NewObjectId(),
	//		User:    u.UserID,
	//		Client:  u.ClientID,
	//		Created: u.CreatedTime,
	//		Expired: u.CreatedTime.Add(Cfg.AccessTokenDuration),
	//	}
	//	u.ExpiredRefreshToken = &DefaultToken{
	//		ID:      bson.NewObjectId(),
	//		User:    u.UserID,
	//		Client:  u.ClientID,
	//		Created: u.CreatedTime,
	//		Expired: u.CreatedTime.Add(Cfg.RefreshTokenDuration),
	//	}

	u.Database.C(TableUser).Insert(u.User1, u.User2)
	u.Database.C(TableClient).Insert(u.Client)
	//	u.Database.C(TableAccessToken).Insert(u.ExpiredAccessToken)
	//	u.Database.C(TableRefreshToken).Insert(u.ExpiredRefreshToken)

	// Generate test resources
	util.CreateDir("resources", (os.ModeDir | os.ModePerm))
	output, _ := os.Create("resources/LICENSE")
	input, _ := os.Open("LICENSE")
	io.Copy(output, input)
	output.Close()
	input.Close()
}

// Teardown cleans up environment.
func (u *TestUnit) Teardown() {
	os.Remove(mongo.ConfigFile)
	os.Remove(debug)

	os.RemoveAll("resources")

	u.Database.DropDatabase()
	u.Session.Close()
}

func parseResult(response *http.Response) *TokenResponse {
	data, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	token := TokenResponse{}
	json.Unmarshal(data, &token)

	return &token
}
