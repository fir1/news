## How to run program locally 
To run the program locally, use the following command:

````
go run cmd/*.go 
````
## Run tests
To run the tests, use the command:
`go test ./...`

## API Endpoints
The application currently provides the following endpoints:
1. ``GET /health`` - This endpoint checks the health of the server.

2. ``GET /news``: This endpoint returns a list of news articles from a public news feed. It allows filtering news articles by category, such as general and technology news. By default, news articles are returned in the order in which they are published. Optionally, you can sort the articles by providing the `sort_by_publish_date` field with values DESC or ASC. Additionally, it allows selecting different sources of news by category and provider (sky, bbc). You can also provide a custom news_source_url as an RSS feed ending with .xml to source news from other providers.


3. ``GET /article``: This endpoint displays a single news article on the screen using an HTML display. You should provide the url query parameter to get a single article converted to HTML display.


## SWAGGER Documentation
The application also has SWAGGER documentation that provides detailed information about the API endpoints. To access the documentation, run the server using the command `go run cmd/*.go` and visit http://localhost:8080/swagger/index.html in your browser.

## Additional features
- The API server utilizes caching to improve response times. Currently, memory caching from the github.com/allegro/bigcache/v3 library is used, but it can be replaced with other caching solutions, such as redis, by implementing the `CacheClientInterface` in `pkg/cache/cache.go`. The use of interfaces allows for easy swapping of caching implementations without changing the application details.

- The project benefits from automatic dependency injection tools, such as "uber/fx", to manage dependencies and facilitate modular and testable code.

- Swagger documentation is implemented for the API, providing better API visibility and documentation

- Test coverage is provided wherever feasible to ensure code quality.

- The application is dockerized, making it easy to deploy and run in containers.

## Docker

In this section, we will explain how to create a Docker image, run the web service, and publish the Docker image to Docker Hub.

**Creating an image from the Dockerfile**
````
docker build -t 2112fir/news -f build/Dockerfile .
````

**Publishing to Docker Hub**
To publish the image to Docker Hub, you need to log in first using either of the following methods:

Login to Docker Hub via CLI
1) Direct usage of password in CLI is not recommended 
````
docker login -p {Passowrd} -u {Username}
````

2) Creating an access token from https://hub.docker.com/settings/security (preferred way):
````
docker login -u 2112fir
````

At the password prompt, enter the personal access token.

You can push a new image to this repository using the CLI
````
docker push 2112fir/news
````

**Running server from the above created docker image**
````
docker run -p 8080:8080 2112fir/news
````

Later on we will use publicly pushed image inside Kubernetes manifest.

Alternatively, you can use the Docker Compose file to build the image locally and test it from your local Docker environment:
````
cd build

docker-compose up
````


# What can be improved ?
- Currently, the BDD (Behavior-Driven Development) integration tests are missing from the project due to time constraints for this project. In real projects, the use of BDD tests with the help of the Cucumber framework would be useful to check the behavior of the business logic, as it allows making client REST API calls and verifying the expected behavior.

- mplement a third-party API provider to leverage any additional functionality related to the news. This would enhance the capabilities of the application by fetching news from multiple sources and providing a richer user experience.