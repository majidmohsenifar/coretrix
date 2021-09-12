# About the project
this is a simple backend for searching product, adding them to cart and delete them.
for search the redisSearch module has been used to index the data related to products.

# How to run project
- install docker and docker-compose 
- clone the project
- run by `docker-compose up -d` which would starts the http server on port 8000

# APIs
for testing API use postman collection provided in project.

# Structure
- the cmd directory contains codes that could be compiled to executable binaries 
- business logic is stored in internal directory
- product package is responsible for product related features including search
- order package is responsible for managing cart (get,update,delete)
- platform package contain the clients and common utilities used in every project like logger,configs and different clients
- mocks contains generated mocks using mockery 

# Points
- In addition to unit tests, integration tests are also written for cart and search functionality
- When the project is run it tries to index some products using a indexProducts command
- Cart is stored in redis for better performance.

#Extra libraries used
- gin for routing
- viper for configs
- zap for logging
- testify for tests

