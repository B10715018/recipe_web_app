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
