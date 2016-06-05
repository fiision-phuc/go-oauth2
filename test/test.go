package test

// Template error message
const (
	ExpectedInvalidParameter     = "Expected \"Invalid %s parameter.\" but found \"%s\"."
	ExpectedNil                  = "Expected nil but found not nil."
	ExpectedNotNil               = "Expected not nil but found nil."
	ExpectedPanic                = "Expected panic but never received."
	ExpectedBoolButFoundBool     = "Expected \"%t\" but found \"%t\"."
	ExpectedNumberButFoundNumber = "Expected \"%d\" but found \"%d\"."
	ExpectedStringButFoundString = "Expected \"%s\" but found \"%s\"."

//	TemplateNil          = "Expected nil but found not nil."
//	TemplateNotNil       = "Expected not nil but found nil."
//	TemplateError        = "Invalid %s parameter."
//	TemplateErrorCode    = "Expected error code %d but found %d."
//	TemplateErrorMessage = "Expected \"Invalid %s parameter.\" but found \"%s\"."
//	TemplateErrorStatus  = "Expected \"%d\" but found \"%d\"."
//	TemplateErrorValue   = "Expected \"%s\" but found \"%s\"."
)

//// SetupOAuth2Server returns an oauth2 server with default settings for testing purpose.
//func SetupOAuth2Server(function func(s *oauth2.Server, tokenStore oauth2.TokenStore, client oauth2.AuthClient, admin oauth2.AuthUser, token *oauth2.TokenResponse)) {
//	defer os.Remove(oauth2.ConfigFile)
//	defer os.Remove(mongo.ConfigFile)
//	mongo.ConnectMongo()

//	// Clean up date after used.
//	session, database := mongo.GetMonotonicSession()
//	defer session.Close()
//	defer database.DropDatabase()

//	// Create single client for testing
//	clientCollection := database.C(oauth2.TableClient)
//	appClient := oauth2.AuthClientDefault{
//		ClientID:     bson.NewObjectId(),
//		ClientSecret: bson.NewObjectId(),
//		GrantTypes:   []string{oauth2.PasswordGrant, oauth2.RefreshTokenGrant},
//		RedirectURIs: []string{"http://www.sample01.com"},
//	}
//	clientCollection.Insert(appClient)

//	// Create single user for testing
//	userCollection := database.C(oauth2.TableUser)
//	password, _ := utils.EncryptPassword("P@ssw0rd")
//	admin := &oauth2.AuthUserDefault{
//		UserID:   bson.NewObjectId(),
//		Username: "admin",
//		Password: password,
//		Roles:    []string{"r_user", "r_admin"},
//	}
//	userCollection.Insert(admin)

//	// Setup server
//	tokenStore := &oauth2.MongoDBTokenStore{}
//	server := oauth2.DefaultServerWithTokenStore(tokenStore)

//	// Login to get access token
//	request, _ := http.NewRequest("POST", "http://localhost:8080/token", strings.NewReader(fmt.Sprintf(
//		"grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s",
//		oauth2.PasswordGrant,
//		appClient.GetClientID(),
//		appClient.GetClientSecret(),
//		"admin",
//		"P@ssw0rd",
//	)))
//	request.Header.Set("content-type", "application/x-www-form-urlencoded")
//	response := httptest.NewRecorder()
//	server.ServeHTTP(response, request)

//	data, _ := ioutil.ReadAll(response.Body)
//	token := oauth2.TokenResponse{}
//	json.Unmarshal(data, &token)

//	function(server, tokenStore, &appClient, admin, &token)
//}
