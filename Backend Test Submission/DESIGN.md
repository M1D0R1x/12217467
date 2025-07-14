# URL Shortener Microservice - Design Document

## Architecture Overview

The URL Shortener Microservice is designed as a modular, maintainable, and scalable service that follows clean architecture principles. The application is structured into distinct layers with clear responsibilities:

1. **API Layer**: Handles HTTP requests/responses and input validation
2. **Service Layer**: Contains the business logic
3. **Storage Layer**: Manages data persistence
4. **Middleware Layer**: Provides cross-cutting concerns like logging
5. **Models Layer**: Defines the data structures
6. **Utilities Layer**: Contains helper functions

## Key Design Decisions

### 1. Modular Architecture

The codebase is organized into separate packages with clear responsibilities:

- `api`: Contains HTTP handlers and request/response processing
- `middleware`: Provides logging functionality
- `models`: Defines data structures used throughout the application
- `storage`: Manages data persistence
- `utils`: Contains utility functions for shortcode generation and validation

This modular approach allows for:
- Clear separation of concerns
- Easier testing of individual components
- Flexibility to replace implementations (e.g., switching storage backends)

### 2. Interface-Based Design

Key components are defined as interfaces, allowing for:
- Loose coupling between components
- Easier mocking for tests
- Flexibility to change implementations without affecting other parts of the system

Examples include:
- `Logger` interface for logging
- `URLStore` interface for storage

### 3. In-Memory Storage

For this implementation, an in-memory storage solution was chosen for simplicity. The storage layer is designed with an interface that would allow for easy replacement with a persistent database in a production environment.

Benefits of the current approach:
- Simplicity for demonstration purposes
- Fast performance for testing
- No external dependencies required

In a production environment, this would be replaced with a database like Redis, MongoDB, or PostgreSQL, depending on scaling requirements.

### 4. Concurrency Handling

The implementation includes proper concurrency handling:
- Mutex locks in the storage layer to prevent race conditions
- Asynchronous processing of click events to ensure fast redirects

### 5. Error Handling

Comprehensive error handling is implemented throughout the application:
- Specific error types for common scenarios (not found, already exists, expired)
- Consistent error response format
- Appropriate HTTP status codes for different error scenarios
- Detailed logging of errors

### 6. Shortcode Generation and Validation

- Secure random generation of shortcodes using crypto/rand
- Validation of custom shortcodes using regular expressions
- Configurable shortcode length

## Data Model

### ShortURL
- **ID**: Unique identifier (shortcode)
- **OriginalURL**: The original long URL
- **CreatedAt**: Creation timestamp
- **ExpiresAt**: Expiration timestamp
- **Clicks**: Number of times the URL has been accessed
- **ClickData**: Detailed information about each click

### Click
- **Timestamp**: When the click occurred
- **Referrer**: Where the click came from
- **Location**: Approximate geographical location
- **UserAgent**: User agent of the client

## API Design

The API follows RESTful principles with clear endpoint definitions:

1. **POST /shorturls**: Create a new shortened URL
2. **GET /shorturls/:shortcode**: Get statistics for a shortened URL
3. **GET /:shortcode**: Redirect to the original URL

## Scalability Considerations

While the current implementation uses in-memory storage, the design allows for scaling in several ways:

1. **Horizontal Scaling**: The stateless nature of the API handlers allows for deploying multiple instances behind a load balancer.

2. **Database Scaling**: The storage interface can be implemented with a distributed database for horizontal scaling of data storage.

3. **Caching**: Frequently accessed URLs could be cached to reduce database load.

4. **Rate Limiting**: Could be added to prevent abuse of the service.

## Future Enhancements

1. **Persistent Storage**: Replace in-memory storage with a database
2. **Authentication and Authorization**: Add user accounts and API keys
3. **Analytics Dashboard**: Provide visual analytics for URL usage
4. **Custom Domains**: Allow users to use custom domains for short URLs
5. **QR Code Generation**: Generate QR codes for shortened URLs
6. **URL Expiration Cleanup**: Implement a background job to clean up expired URLs