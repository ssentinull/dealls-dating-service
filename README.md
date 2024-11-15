<a id="readme-top"></a>

<br />

<div align="center">
  <h3 align="center">dealls-dating-service</h3>
  <p align="center">
    The backend service that powers Dealls! dating app.
  </p>
</div>

<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
        <li><a href="#codebase-structure">Codebase Structure</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#run-the-service">Run The Service</a></li>
        <li><a href="#test-the-service">Test The Service</a></li>
      </ul>
    </li>
  </ol>
</details>


## About The Project

This is the repository for the backend service of Dealls! dating app. The dating app itself is similar to Tinder and Bumble:
* You signup and login using your email and password.
* You fill out the characteristics that you're looking for in a future partner.
* You search through potential candidates using a *stacked-card* mechanism:
    * Swipe right on the card to show interest.
    * Swipe left on the card to ignore.
* If the person that you showed interest are also interested in you, then you found a match!


### Built With

These are the technologies that I used to build this service:

* Golang
* PostgreSQL
* Redis
* REST API
* Swagger API Documentation


### Codebase Structure

This is how the codebase is structured:

```
dealls-dating-service/
├─ .modd/
├─ cmd/
├─ docs/
│  ├─ api/
│  ├─ sql/
├─ internal/
│  ├─ business/
│  │  ├─ domain/
│  │  ├─ model/
│  │  ├─ usecase/
│  ├─ handler/
│  ├─ types/
│  ├─ mocks/
├─ pkg/
│  ├─ build_util/
│  ├─ common/
│  ├─ constants/
│  ├─ stdlib/
├─ main.go
├─ config.yaml
├─ go.mod
├─ go.sum
├─ Makefile
├─ README.md
├─ .gitignore
├─ swagger-cli
```

These are the details of each directories in the structure:
* `/.modd/`: configuration for [modd](https://github.com/cortesi/modd) hot-reload daemon tool.
* `/cmd/`: entry point for executables (e.g.: database migration executable, web server executable).
* `/docs/api/`: swagger documentation definitions that we later use to generate concrete golang types using [go-swagger](https://github.com/go-swagger/go-swagger).
* `/docs/sql/`: directory for database migration files.
* `/internal/business/domain/`: repository layer for grouped objects based on their domain.
* `/internal/business/model/`: data model definitions.
* `/internal/business/usecase/`: business logic layer.
* `/internal/handler/`: layer to encode and decode data models passed on from the usecase layer based on the handler user (e.g.: rest, graphql, pubusb).
* `/internal/types/`: the target directory of files generated by go-swagger.
* `/internal/mocks/`: the target directory of files generated by [mock generator](https://github.com/uber-go/mock).
* `/pkg/build_util/`: build utility definitions.
* `/pkg/common/`: helper functions for common usecases.
* `/pkg/constans/`: constant variable definitions for codebase-wide usecases.
* `/pkg/stdlib/`: interface wrapping for external golang packages.

## Getting Started

If you want to run this service locally, you need to have the Prerequisites first before trying to execute the Installation phase.

### Prerequisites

You need to install these dependencies before running the service:
* Golang version 1.22 or above
* [Make](https://www.gnu.org/software/make/)
* [modd](https://github.com/cortesi/modd)
  ```sh
  $ go install github.com/cortesi/modd/cmd/modd@latest
  ```
* [Mock Generator](https://github.com/uber-go/mock)
  ```sh
  $ go install go.uber.org/mock/mockgen@latest
  ```
* [go-swagger cli](https://github.com/go-swagger/go-swagger) (this is done in the root directory after you've cloned the repo)
  ```sh
  $ make swagger-install
  ```
* [golang-migrate](https://github.com/golang-migrate/migrate)
* [Redis server](https://redis.io/docs/latest/operate/oss_and_stack/install/install-redis/)
* [PostgreSQL server](https://www.postgresql.org/download/)

### Run The Service

Once you've install the prerequisite, follow these steps to start the service:

1. Rename `example.config.yaml` to `config.yaml` in the root directory and fill in the variables. You don't need to fill all of it, just the `Redis.Address` and `SQL.Leader.DSN` variables.

2. Run the migration executable to create the database tables.
   ```sh
   $ make migrate-up
   ```
3. Generate the concrete golang types from the swagger definitions. If this doesn't work, you can skip this step because the concrete types have been pushed to the repo.
   ```sh
   $ make swagger
   ```
4. Run the service in normal mode.
   ```sh
   $ make run
   ```
   Or daemon mode.
    ```sh
   $ make run-dev
   ```

### Test The Service

* The service has unit tests in place. Run the following command to execute the unit tests.
    ```sh
    $ make test
    ```

* Once the service is running, you can open [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) to checkout the API documentation through Swagger UI. 

* If you want to test the APIs, you can do so using the Swagger UI, or through [Postman](https://www.postman.com/) by importing the [collection file](https://github.com/ssentinull/dealls-dating-service/blob/master/postman-collection.json) located in the root directory.

<p align="right">(<a href="#readme-top">back to top</a>)</p>
