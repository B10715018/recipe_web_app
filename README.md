# recipe_web_app

# To start

- clone this repository:

```
git clone https://github.com/B10715018/recipe_web_app.git
```

## To install Go Dependency

- execute following command:

```
go mod download
```

## Important Notes

- Alternatives for testing without Postman by using cURL:

```
curl --location --request POST 'http://localhost:8080/recipes' \
--header 'Content-Type: application/json' \
--data-raw '{
   "name": "Homemade Pizza",
   "tags" : ["italian", "pizza", "dinner"],
   "ingredients": [
       "1 1/2 cups (355 ml) warm water (105°F-115°F)",
       "1 package (2 1/4 teaspoons) of active dry yeast",
       "3 3/4 cups (490 g) bread flour",
       "feta cheese, firm mozzarella cheese, grated"
   ],
   "instructions": [
       "Step 1.",
       "Step 2.",
       "Step 3."
   ]
}' | jq -r
```

Extra Note:
The jq utility https://stedolan.github.io/jq/ is used to format the response body in JSON format. It's a powerful command-line JSON processor.

- For testing the GET command for Recipes Web Application

```
curl -s --location --request GET 'http://localhost:8080/recipes' \
--header 'Content-Type: application/json'
```

- The `|` symbol in the terminal commands mean pipe output from first command to input of the second command

- To count the number of recipes returned by the request:

```
curl -s -X GET 'http://localhost:8080/recipes' | jq length
```

-

- To test the DELETE endpoint without postman, use terminal command:

```
curl -v -sX DELETE http://localhost:8080/recipes/{recipeId} | jq -r
```

## API Documentation

- Use OpenAPI Specification to generate the describe the specification of our API

- We need to create the comment in the `main.go` file that mnatches the swagger metadata syntax likewise:

```
// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin.
//
//  Schemes: http
//  Host: localhost:8080
//  BasePath: /
//  Version: 1.0.0
//  Contact: Brandon<brandon.wymer@rocketmail.com>
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
// swagger:meta
```

- These comments include things such as the API's description, version, base URL, and so on. There are more fields that you can include (a full list is available at https://goswagger.io/use/spec/meta.html).

- To generate the swagger json file execute:

```
swagger generate spec -o swagger.json
```

- To deploy the swagger ui execute:

```
swagger serve ./swagger.json
```

- If you're a fan of the Swagger UI, you can set the flavor flag to swagger with the following command:

```
swagger serve -F swagger ./swagger.json
```

- After `swagger: metadata` we need to define the `swagger:operation`. You can find all the properties at https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#operationObject.

## Data Persistency

- Recommended creating DB using docker:

```
docker run -d --name mongodb -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -p 27017:27017 mongo:4.4.3
```

- To remove container:

```
docker rm -f {CONTAINER_NAME} || true
```

- To display logs of the container, execute:

```
docker logs –f {CONTAINER_ID}
```

- To interact with MongoDB is recommended to use MongoDB compass, but we can also communicate directly with the MongoDB server

- To connect with MongoDB Compass

```
 mongodb://admin:password@localhost:27017/test.
```

- To Stop Docker Container from running:

```
docker stop {CONTAINER_ID}
```

- Check if application can successfully connect:

```
MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" go run main.go
```

- Complete doc for mongoDB driver:
  You can view the full documentation for the MongoDB Go driver on the GoDoc website (https://godoc.org/go.mongodb.org/mongo-driver).

- Populate the mongoDB without writing code in Golang

```
mongoimport --username admin --password password --authenticationDatabase admin --db demo --collection recipes --file recipes.json --jsonArray
```

- Don't forget to initialize the mongoDB collection and also the database to prevent from unwanted error !!

- Refactor Go Project Layout into handler folder, model folder and main file, to execute, as you can see it will run all the Go files:

```
MONGO_URI="mongodb://admin:password@localhost:27017/admin?authSource=admin" MONGO_DATABASE=demo go run *.go
```
