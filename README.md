## Getting started

Make sure that you're in the root of the project directory, fetch the dependencies with `go mod tidy`, then run the
application using `go run main.go serve`:

```
$ go mod tidy
$ go run main.go serve
```

## Project structure

Everything in the codebase is designed to be editable. Feel free to change and adapt it to meet your needs.

|     |     |
| --- | --- |
| **`assets`** | Contains the non-code assets for the application. |
| `↳ assets/migrations/` | Contains SQL migrations. |
| `↳ assets/efs.go` | Declares an embedded filesystem containing all the assets. |

|     |     |
| --- | --- |
| **`cmd/api`** | Application-specific code (handlers, routing, middleware) for dealing with HTTP requests and responses. |

|     |                                                                                                                                  |
| --- |----------------------------------------------------------------------------------------------------------------------------------|
| **`internal`** | Contains private application and library code. |
| `↳ internal/database/` | Contains database-related code. |
| `↳ internal/server/` | Contains a helper function for starting and gracefully shutting down the server. |
| `↳ internal/ecommerce/` | Contains the ecommerce domain specific code. |

|                          |     |
|--------------------------| --- |
| **`pkg`**                | Contains various helper packages used by the application. |
| `↳ internal/leveledlog/` | Contains a leveled logger implementation. |
| `↳ internal/request/`    | Contains helper functions for decoding JSON requests. |
| `↳ internal/response/`   | Contains helper functions for sending JSON responses. |
| `↳ internal/validator/`  | Contains validation helpers. |
| `↳ internal/version/`    | Contains the application version number definition. |

## Configuration settings

Configuration settings are managed via yaml file or os environment.
os enviroment > config yaml > default config in `cmd/api/config.go`

## How to start at local environment
Make sure that you're in the root of the project directory, start docker instance using:
`docker-compose up --build`

If you make a request to the `GET /_/status` endpoint using `curl` you should get a response like this:
```
$ curl -i localhost:3000/_/status
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 29 Nov 2022 17:46:40 GMT
Content-Length: 20

{
	"Status": "OK"
}
```
Application contains 2 endpoins:
##### get product facets:
GET `/v1/products/facets`
```
$ curl -i localhost:3000/_/status
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 29 Nov 2022 17:46:40 GMT
Content-Length: 20

{
	"Status": "OK"
}
```
##### get product by search term & filter, ordering:
GET `/v1/products?term={keyword}&{filter attribute}={value}:{operator}&sort={field}:{asc|desc}&limit={number}&offset={number}`
For example: I want to get product with: search term is `woman` and collection is `winter` and price `greater than 99` with curl:
```
curl -i 'localhost:3000/products?term=woman&collection=winter&price=99:gt'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 29 Nov 2022 17:58:33 GMT
Content-Length: 466

[
	{
		"id": 1,
		"name": "Clover cut-out knitted top",
		"description": "Taking form in a body-skimming knit with bold cut-out and cherry red for the festive season woman.",
		"price": 100.99,
		"attributes": [
			"colors:red",
			"colors:pink",
			"colors:orange",
			"size:xl",
			"size:x",
			"size:s",
			"categories:woman-top",
			"brand:amadeus",
			"collection:winter"
		],
		"created_at": "0001-01-01T00:00:00Z",
		"updated_at": "0001-01-01T00:00:00Z"
	}
]
```
## Explaint
The reasons to use only one table `products` instead of use multiple models to store data like `products`,`categories`,`attributes` then make that's relationship with pk is:
I was resolve the issue of search, sorted ... then just index into only one table for do everythings with `tsvector`, for example: I want to filter the product with categories woman
```
SELECT * FROM "products" AS "product" WHERE (tsv @@ 'categories\:woman'::tsquery);
```
And can be collect facets data
```
SELECT
  split_part(word, ':', 1) AS attr,
  split_part(word, ':', 2) AS value,
  ndoc AS count
FROM ts_stat($$ SELECT tsv FROM products $$)
ORDER BY word;
```
|   attr  |   value  |   count  |
| --- | --- |--- |
| brand | amadeus  | 3 |
| brand | dearjose | 2 |
| categories | woman-top | 5 |
| collection | spring | 2 |
| ... | ... | ... |

How I order product with default sort condition:
I used `ts_rank` to set weight of `name` stronger than `description` so you can check the result with search term `clover` using command:
```
curl -i 'localhost:3000/products?term=clover'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 29 Nov 2022 18:07:54 GMT
Content-Length: 969

[
	{
		"id": 1,
		"name": "Clover cut-out knitted top",
		"description": "Taking form in a body-skimming knit with bold cut-out and cherry red for the festive season woman.",
		"price": 100.99,
		"attributes": [
			"colors:red",
			"colors:pink",
			"colors:orange",
			"size:xl",
			"size:x",
			"size:s",
			"categories:woman-top",
			"brand:amadeus",
			"collection:winter"
		],
		"created_at": "0001-01-01T00:00:00Z",
		"updated_at": "0001-01-01T00:00:00Z"
	},
	{
		"id": 4,
		"name": "Bunny off-shoulder poplin top",
		"description": "Comfort abounds in an off-shoulder crop top with sleeve accents and elegant laps. Clover Bunny Top from the Holiday 2023 Collection woman.",
		"price": 75.25,
		"attributes": [
			"colors:red",
			"colors:pink",
			"colors:green",
			"size:s",
			"size:x",
			"size:m",
			"categories:woman-top",
			"brand:amadeus",
			"collection:spring"
		],
		"created_at": "0001-01-01T00:00:00Z",
		"updated_at": "0001-01-01T00:00:00Z"
	}
]
```

by default the keyword `clover` in product name will on top in description.