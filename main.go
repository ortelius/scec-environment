// Ortelius v11 Environment Microservice that handles creating and retrieving Environments
package main

import (
	"context"
	"encoding/json"

	_ "cli/docs"

	driver "github.com/arangodb/go-driver/v2/arangodb"
	"github.com/arangodb/go-driver/v2/arangodb/shared"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/ortelius/scec-commons/database"
	"github.com/ortelius/scec-commons/model"
)

var logger = database.InitLogger()
var dbconn = database.InitializeDB("evidence")

// GetEnvironments godoc
// @Summary Get a List of Environments
// @Description Get a list of environments for the user.
// @Tags environment
// @Accept */*
// @Produce json
// @Success 200
// @Router /msapi/environment [get]
func GetEnvironments(c *fiber.Ctx) error {

	var cursor driver.Cursor       // db cursor for rows
	var err error                  // for error handling
	var ctx = context.Background() // use default database context

	// query all the environments in the collection
	aql := `FOR environment in evidence
			FILTER (Environment.objtype == 'Environment')
			RETURN environment`

	// execute the query with no parameters
	if cursor, err = dbconn.Database.Query(ctx, aql, nil); err != nil {
		logger.Sugar().Errorf("Failed to run query: %v", err) // log error
	}

	defer cursor.Close() // close the cursor when returning from this function

	var environments []*model.Environment // define a list of environments to be returned

	for cursor.HasMore() { // loop thru all of the documents

		environment := model.NewEnvironment() // fetched environment
		var meta driver.DocumentMeta          // data about the fetch

		// fetch a document from the cursor
		if meta, err = cursor.ReadDocument(ctx, environment); err != nil {
			logger.Sugar().Errorf("Failed to read document: %v", err)
		}
		environments = append(environments, environment)                     // add the environment to the list
		logger.Sugar().Infof("Got doc with key '%s' from query\n", meta.Key) // log the key
	}

	return c.JSON(environments) // return the list of environments in JSON format
}

// GetEnvironment godoc
// @Summary Get a Environment
// @Description Get a environment based on the _key or name.
// @Tags environment
// @Accept */*
// @Produce json
// @Success 200
// @Router /msapi/environment/:key [get]
func GetEnvironment(c *fiber.Ctx) error {

	var cursor driver.Cursor       // db cursor for rows
	var err error                  // for error handling
	var ctx = context.Background() // use default database context

	key := c.Params("key")                // key from URL
	parameters := map[string]interface{}{ // parameters
		"key": key,
	}

	// query the environments that match the key or name
	aql := `FOR environment in evidence
			FILTER (environment.name == @key or environment._key == @key)
			RETURN environment`

	// run the query with patameters
	if cursor, err = dbconn.Database.Query(ctx, aql, &driver.QueryOptions{BindVars: parameters}); err != nil {
		logger.Sugar().Errorf("Failed to run query: %v", err)
	}

	defer cursor.Close() // close the cursor when returning from this function

	environment := model.NewEnvironment() // define a environment to be returned

	if cursor.HasMore() { // environment found
		var meta driver.DocumentMeta // data about the fetch

		if meta, err = cursor.ReadDocument(ctx, environment); err != nil { // fetch the document into the object
			logger.Sugar().Errorf("Failed to read document: %v", err)
		}
		logger.Sugar().Infof("Got doc with key '%s' from query\n", meta.Key)

	} else { // not found so get from NFT Storage
		if jsonStr, exists := database.MakeJSON(key); exists {
			if err := json.Unmarshal([]byte(jsonStr), environment); err != nil { // convert the JSON string from LTF into the object
				logger.Sugar().Errorf("Failed to unmarshal from LTS: %v", err)
			}
		}
	}

	return c.JSON(environment) // return the environment in JSON format
}

// NewEnvironment godoc
// @Summary Create a Environment
// @Description Create a new Environment and persist it
// @Tags Environment
// @Accept application/json
// @Produce json
// @Success 200
// @Router /msapi/environment [post]
func NewEnvironment(c *fiber.Ctx) error {

	var err error                         // for error handling
	var meta driver.DocumentMeta          // data about the document
	var ctx = context.Background()        // use default database context
	environment := new(model.Environment) // define a environment to be returned

	if err = c.BodyParser(environment); err != nil { // parse the JSON into the environment object
		return c.Status(503).Send([]byte(err.Error()))
	}

	cid, dbStr := database.MakeNFT(environment) // normalize the object into NFTs and JSON string for db persistence

	logger.Sugar().Infof("%s=%s\n", cid, dbStr) // log the new nft

	var resp driver.CollectionDocumentCreateResponse
	// add the environment to the database.  Ignore if it already exists since it will be identical
	if resp, err = dbconn.Collection.CreateDocument(ctx, environment); err != nil && !shared.IsConflict(err) {
		logger.Sugar().Errorf("Failed to create document: %v", err)
	}
	meta = resp.DocumentMeta
	logger.Sugar().Infof("Created document in collection '%s' in db '%s' key='%s'\n", dbconn.Collection.Name(), dbconn.Database.Name(), meta.Key)

	return c.JSON(environment) // return the environment object in JSON format.  This includes the new _key
}

// setupRoutes defines maps the routes to the functions
func setupRoutes(app *fiber.App) {

	app.Get("/swagger/*", swagger.HandlerDefault)      // handle displaying the swagger
	app.Get("/msapi/environment", GetEnvironment)      // list of environments
	app.Get("/msapi/environment/:key", GetEnvironment) // single environment based on name or key
	app.Post("/msapi/environment", NewEnvironment)     // save a single environment
}

// @title Ortelius v11 Environment Microservice
// @version 11.0.0
// @description RestAPI for the Environment Object
// @description ![Release](https://img.shields.io/github/v/release/ortelius/scec-environment?sort=semver)
// @description ![license](https://img.shields.io/github/license/ortelius/.github)
// @description
// @description ![Build](https://img.shields.io/github/actions/workflow/status/ortelius/scec-environment/build-push-chart.yml)
// @description [![MegaLinter](https://github.com/ortelius/scec-environment/workflows/MegaLinter/badge.svg?branch=main)](https://github.com/ortelius/scec-environment/actions?query=workflow%3AMegaLinter+branch%3Amain)
// @description ![CodeQL](https://github.com/ortelius/scec-environment/workflows/CodeQL/badge.svg)
// @description [![OpenSSF-Scorecard](https://api.securityscorecards.dev/projects/github.com/ortelius/scec-environment/badge)](https://api.securityscorecards.dev/projects/github.com/ortelius/scec-environment)
// @description
// @description ![Discord](https://img.shields.io/discord/722468819091849316)

// @termsOfService http://swagger.io/terms/
// @contact.name Ortelius Google Group
// @contact.email ortelius-dev@googlegroups.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /msapi/environment
func main() {
	port := ":" + database.GetEnvDefault("MS_PORT", "8080") // database port
	app := fiber.New()                                      // create a new fiber application
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowOrigins: "*",
	}))

	setupRoutes(app) // define the routes for this microservice

	if err := app.Listen(port); err != nil { // start listening for incoming connections
		logger.Sugar().Fatalf("Failed get the microservice running: %v", err)
	}
}
