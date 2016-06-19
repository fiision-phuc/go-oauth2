package oauth2

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/phuc0302/go-oauth2/mongo"
	"github.com/phuc0302/go-oauth2/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	session  *mgo.Session
	database *mgo.Database

	client IClient
	user1  IUser
	user2  IUser

	expiredAccessToken  IToken
	expiredRefreshToken IToken

	username       = "admin"
	password       = "admin"
	userID         = bson.NewObjectId()
	clientID       = bson.NewObjectId()
	clientSecret   = bson.NewObjectId()
	createdTime, _ = time.Parse(time.RFC822, "02 Jan 06 15:04 MST")
)

func setup() {
	mongo.ConnectMongo()
	session, database = mongo.GetMonotonicSession()

	// Define global variables
	Cfg = loadConfig(debug)
	objectFactory = &DefaultFactory{}
	TokenStore = objectFactory.CreateStore()

	// Generate test data
	client = &DefaultClient{
		ID:     clientID,
		Secret: clientSecret,
		Grants: []string{AuthorizationCodeGrant, PasswordGrant, RefreshTokenGrant},

		Redirects: []string{"http://www.sample01.com", "http://www.sample02.com"},
	}

	password1, _ := util.EncryptPassword("admin")
	user1 = &DefaultUser{
		ID:    userID,
		User:  "admin",
		Pass:  password1,
		Roles: []string{"r_user", "r_admin"},
	}

	password2, _ := util.EncryptPassword(clientSecret.Hex())
	user2 = &DefaultUser{
		ID:    clientID,
		User:  clientID.Hex(),
		Pass:  password2,
		Roles: []string{"r_device"},
	}

	//	expiredAccessToken = &DefaultToken{
	//		ID:      bson.NewObjectId(),
	//		User:    userID,
	//		Client:  clientID,
	//		Created: createdTime,
	//		Expired: createdTime.Add(Cfg.AccessTokenDuration),
	//	}
	//	expiredRefreshToken = &DefaultToken{
	//		ID:      bson.NewObjectId(),
	//		User:    userID,
	//		Client:  clientID,
	//		Created: createdTime,
	//		Expired: createdTime.Add(Cfg.RefreshTokenDuration),
	//	}

	database.C(TableUser).Insert(user1, user2)
	database.C(TableClient).Insert(client)
	//	database.C(TableAccessToken).Insert(expiredAccessToken)
	//	database.C(TableRefreshToken).Insert(expiredRefreshToken)

	// Generate test resources
	util.CreateDir("resources", (os.ModeDir | os.ModePerm))
	output, _ := os.Create("resources/LICENSE")
	input, _ := os.Open("LICENSE")
	io.Copy(output, input)
	output.Close()
	input.Close()
}

func teardown() {
	os.Remove(mongo.ConfigFile)
	os.Remove(debug)

	os.RemoveAll("resources")

	database.DropDatabase()
	session.Close()
}

func parseError(response *http.Response) *util.Status {
	data, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	status := util.Status{}
	json.Unmarshal(data, &status)

	return &status
}

func parseResult(response *http.Response) *TokenResponse {
	data, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	token := TokenResponse{}
	json.Unmarshal(data, &token)

	return &token
}
