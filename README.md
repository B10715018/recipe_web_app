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

## Caching API with Redis

- Need to use Redis as a caching layer to improve user experience by giving the data needed faster in term of speed performance

- Easiest way to start Redis server with Docker container by:

```
docker run -d --name redis -p 6379:6379 redis:6.0
```

This command does the following two main things:

The –d flag runs the Redis container as a daemon.
The –p flag maps port 6379 of the container to port 6379 of the host. Port 6379 is the port where the Redis server is exposed.

- In production we need to set eviction policy for the Redis, for example using Least Recently Used Algorithm (LRU) by writing this in `redis.conf` file:

```
maxmemory-policy allkeys-lru
maxmemory 512mb
```

- We can pass the redis configuration file during the runtime of the container

```
docker run -d -v $PWD/conf:/usr/local/etc/redis --name redis -p 6379:6379 redis:6.0
```

Here, `$PWD/conf` is the folder containing the redis.conf file.

- We can verify that data is being cached in Redis by running the Redis CLI from the container. Run the following commands:

```
docker ps
docker exec –it CONTAINER_ID bash
```

- Now that we're attached to the Redis container, we can use the Redis command line:

```
redis-cli
```

- From there, we can use the EXISTS command to check if the recipes key exists:

```
EXISTS recipes
```

This command will return 1 (if the key exists) or 0 (if the key doesn't exist). In our case, the list of recipes has been cached in Redis:

- For GUI fans, you can use Redis Insights (https://redislabs.com/fr/redis-enterprise/redis-insight/). It provides an intuitive interface to explore Redis and interact with its data. Similar to the Redis server, you can deploy Redis Insights with Docker:

```
docker run -d --name redisinsight --link redis -p 8001:8001 redislabs/redisinsight
```

- This command will run a container based on the Redis Insight official image and expose the interface on port 8001.

- Navigate with your browser to http://localhost:8081. The Redis Insights home page should appear. Click on I already have a database and then on the Connect to Redis database button

- Set the Host to redis, port to 6379, and name the database. The settings are as follows

- Next, click on ADD REDIS DATABASE. The local database will be saved; click on it:

- Some Problem when using caching system is there will be data inconsistency, how to fix is by setting TTL (Time To Live) and delete the recipes in the redis every time there is update or add recipe

## Performance Benchmark

- We can take this further and see how the API will behave under a huge volume of requests. We can simulate multiple requests with Apache Benchmark (https://httpd.apache.org/docs/2.4/programs/ab.html).

- First, let's test the API without the caching layer. You can run 2,000 GET requests in total on the /recipes endpoint with 100 concurrent requests with the following command:

```
ab -n 2000 -c 100 -g without-cache.data http://localhost:8080/recipes
```

- Next, we will issue the same requests but this time on the API with caching (Redis):

```
ab -n 2000 -c 100 -g with-cache.data http://localhost:8080/recipes
```

- To compare both results, we can use the gnuplot utility to plot a chart based on the without-cache.data and with-cache.data files. But first, create an apache-benchmark.p file to render data into a graph:

```
set terminal png
set output "benchmark.png"
set title "Cache benchmark"
set size 1,0.7
set grid y
set xlabel "request"
set ylabel "response time (ms)"
plot "with-cache.data" using 9 smooth sbezier with lines title "with cache", "without-cache.data" using 9 smooth sbezier with lines title "without cache"
```

- These commands will draw two plots on the same graph based on the .data files and save the output as a PNG image. Next, run the gnuplot command to create the image:

```
gnuplot apache-benchmark.p
```

## Authentication

### Authentication with API key

- Run the application after running the MongoDB and Redis containers, but this time set the X-API-KEY environment variable as follows:

```
X_API_KEY=eUbP9shywUygMx7u  MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" MONGO_DATABASE=demo go run *.go
```

- Set the API key in the Postman OR use cURL:

```
curl --location --request POST 'http://localhost:8080/recipes' \
--header 'X-API-KEY: eUbP9shywUygMx7u' \
--header 'Content-Type: application/json' \
--data-raw '{
   "name": "Homemade Pizza",
   "ingredients": ["..."],
   "instructions": ["..."],
   "tags": ["dinner", "fastfood"]
}'
```

- Using API key only is vulnerable to the man in the middle attack and need to encrypt the data by using JWT that consist of 3 part: header, payload and signature

- Header where we specify the algorithm for the signature, payload is where the data lives and the signature is the result of hashing header and payload parts with a secrets key

- When the user is signed in, it will get JWT token that will be used for authorization whenever it try to made a request to API

- To test our app with JWT execute:

```
JWT_SECRET=eUbP9shywUygMx7u MONGO_URI="mongodb://admin:password@localhost:27017/admin?authSource=admin" MONGO_DATABASE=demo go run *.go
```

- To decode the JWT token:
  The token consists of three parts separated by a dot. You can decode the token by going to https://jwt.io/ to return the following output (your results might look different)

## Populate the MongoDB with user

- Create new project and `main.go` file like below:

```
package main

import (
	"context"
	"crypto/sha256"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	users := map[string]string{
		"admin":      "fCRmh4Q2J7Rseqkz",
		"packt":      "RE4zfHB35VPtTkbT",
		"businessboy": "L3nSFRcZzNQ67bcc",
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	h := sha256.New()

	for username, password := range users {
		collection.InsertOne(ctx, bson.M{
			"username": username,
			"password": string(h.Sum([]byte(password))),
		})
	}
}
```

- and run this command:

```
MONGO_URI="mongodb://admin:password@localhost:27017/admin?authSource=admin" MONGO_DATABASE=demo go run main.go
```
