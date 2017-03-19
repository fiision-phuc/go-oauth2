package oauth2

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/phuc0302/go-mongo"
	"github.com/phuc0302/go-oauth2/oauth_table"
	"github.com/phuc0302/go-server"
	"github.com/phuc0302/go-server/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// TestEnv describes an implementation for OAuth2 unit test.
type TestEnv struct {
	Session  *mgo.Session
	Database *mgo.Database

	User1  User
	User2  User
	Client Client

	ExpiredAccessToken  Token
	ExpiredRefreshToken Token

	Username     string
	Password     string
	UserID       bson.ObjectId
	ClientID     bson.ObjectId
	ClientSecret bson.ObjectId
	CreatedTime  time.Time
}

// Setup initializes environment.
func (u *TestEnv) Setup() {
	InitializeWithMongoDB(true, false)

	u.Session, u.Database = mongo.GetMonotonicSession()
	u.Username = "admin"
	u.Password = "Password"
	u.UserID = bson.NewObjectId()
	u.ClientID = bson.NewObjectId()
	u.ClientSecret = bson.NewObjectId()
	u.CreatedTime, _ = time.Parse(time.RFC822, "02 Jan 06 15:04 MST")

	// Generate test data
	u.Client = &MongoDBClient{
		ID:     u.ClientID,
		Secret: u.ClientSecret,
		Grants: []string{AuthorizationCodeGrant, PasswordGrant, RefreshTokenGrant},

		Redirects: []string{"http://www.sample01.com", "http://www.sample02.com"},
	}

	password1, _ := util.EncryptPassword(u.Password)
	u.User1 = &MongoDBUser{
		ID:    u.UserID,
		User:  u.Username,
		Pass:  password1,
		Roles: []string{"r_user", "r_admin"},
	}

	password2, _ := util.EncryptPassword(u.ClientSecret.Hex())
	u.User2 = &MongoDBUser{
		ID:    u.ClientID,
		User:  u.ClientID.Hex(),
		Pass:  password2,
		Roles: []string{"r_device"},
	}

	u.ExpiredAccessToken = &MongoDBToken{
		ID:      bson.NewObjectId(),
		User:    u.UserID,
		Client:  u.ClientID,
		Created: u.CreatedTime,
		Expired: u.CreatedTime.Add(Cfg.AccessTokenDuration),
	}
	u.ExpiredRefreshToken = &MongoDBToken{
		ID:      bson.NewObjectId(),
		User:    u.UserID,
		Client:  u.ClientID,
		Created: u.CreatedTime,
		Expired: u.CreatedTime.Add(Cfg.RefreshTokenDuration),
	}

	u.Database.C(oauthTable.User).Insert(u.User1, u.User2)
	u.Database.C(oauthTable.Client).Insert(u.Client)
	u.Database.C(TableAccessToken).Insert(u.ExpiredAccessToken)
	u.Database.C(TableRefreshToken).Insert(u.ExpiredRefreshToken)

	// Generate test resources
	util.CreateDir("resources", (os.ModeDir | os.ModePerm))
	output, _ := os.Create("resources/LICENSE")
	input, _ := os.Open("LICENSE")
	io.Copy(output, input)
	output.Close()
	input.Close()
}

// Teardown cleans up environment.
func (u *TestEnv) Teardown() {
	os.Remove(mongo.ConfigFile)
	os.Remove(server.Debug)

	os.RemoveAll("resources")

	u.Database.DropDatabase()
	u.Session.Close()
}

// parseResult parses OAuth response.
func parseResult(response *http.Response) *OAuthResponse {
	data, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	token := OAuthResponse{}
	json.Unmarshal(data, &token)

	return &token
}
