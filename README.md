# Payment Gateway

This project is a robust payment gateway system built with Go, providing secure and efficient payment processing capabilities.

## Features

- User authentication and authorization
- Merchant management
- Payment processing
- Refund handling
- Prometheus metrics
- Swagger API documentation
- PostgreSQL database
- Docker support

## Prerequisites

- Docker and Docker Compose
- Go 1.22 or later
- Make

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/popeskul/payment-gateway.git
   cd payment-gateway
   go mod tydy
   ```

2. Set up environment variables:
   Create a `.env` file in the root directory and add the following variables:
   ```
   DB_PASSWORD=your_database_password
   ACCESS_TOKEN_SECRET=your_access_token_secret
   REFRESH_TOKEN_SECRET=your_refresh_token_secret
   GRAFANA_ADMIN_PASSWORD=your_grafana_password
   ```

3. Generate Swagger documentation:
   ```
   make generate-swagger
   ```

4. Set up the database migrations:
   ```
   make migrate-setup DB_USER=paymentuser DB_PASSWORD=your_database_password DB_HOST=localhost DB_PORT=5432 DB_NAME=paymentdb DB_SSLMODE=disable
   ```

## Running the Application

1. Start the application and all required services:
   ```
   docker-compose up --build
   ```

2. Once the application is running, the following services will be available:

   | Service | Port | Description |
      |---------|------|-------------|
   | Payment Gateway API | 8080 | The main application API |
   | Swagger UI | 8081 | API documentation interface |
   | PostgreSQL Database | 5432 | Database for the application |
   | Prometheus | 9090 | Metrics and monitoring |
   | Grafana | 3000 | Visualization for Prometheus metrics |

   You can access these services as follows:
   - Main API: `http://localhost:8080`
   - Swagger UI: `http://localhost:8081`
   - Prometheus: `http://localhost:9090`
   - Grafana: `http://localhost:3000`
   - Metrics endpoint: `http://localhost:8080/metrics`

   Note: The PostgreSQL database is also exposed on port 5432, but it's typically accessed by the application internally and not meant for direct external access unless needed for development or debugging purposes.

3. To stop the application and all services, use:
   ```
   docker-compose down
   ```

   If you want to remove all data volumes as well, use:
   ```
   docker-compose down -v
   ```

## Areas for Improvement

1. **Test Coverage**: Increase test coverage, especially for edge cases and error handling scenarios.

2. **API Documentation**: Enhance the Swagger documentation with more detailed descriptions and examples.

3. **Logging**: Implement more comprehensive logging throughout the application for better debugging and monitoring.

4. **Error Handling**: Implement a more robust error handling system with custom error types and consistent error responses.

5. **Security Enhancements**:
    - Implement rate limiting to prevent abuse
    - Add support for HTTPS
    - Implement more robust input validation and sanitization

6. **Performance Optimization**:
    - Implement caching mechanisms for frequently accessed data
    - Optimize database queries and indexes

7. **Scalability**:
    - Implement horizontal scaling capabilities
    - Consider using message queues for asynchronous processing of payments and refunds

8. **Monitoring and Alerting**:
    - Set up more detailed Prometheus metrics
    - Configure Grafana dashboards and alerts

9. **Containerization**:
    - Optimize Docker images for smaller size and faster builds
    - Implement Docker health checks

10. **CI/CD**: Set up a comprehensive CI/CD pipeline for automated testing, building, and deployment.

11. **Documentation**: Improve in-code documentation and add architecture diagrams to explain the system design.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.