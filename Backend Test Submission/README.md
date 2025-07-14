# URL Shortener Microservice

A robust HTTP URL Shortener Microservice that provides core URL shortening functionality along with basic analytical capabilities for the shortened links.

## Architecture

The microservice is built using a modular architecture with the following components:

1. **API Layer** - Handles HTTP requests and responses
2. **Storage Layer** - Manages the persistence of shortened URLs and click data
3. **Middleware** - Provides cross-cutting concerns like logging
4. **Utilities** - Contains helper functions for shortcode generation and validation
5. **Models** - Defines the data structures used throughout the application

## Features

- Create shortened URLs with optional custom shortcodes
- Set custom validity periods for shortened URLs (default: 30 minutes)
- Redirect to original URLs via shortened links
- Track and retrieve statistics for shortened URLs
- Extensive logging of all operations

## API Endpoints

### Create Short URL

- **Method**: POST
- **Route**: `/shorturls`
- **Request Body**:
  ```json
  {
    "url": "https://example.com/very-long-url",
    "validity": 30,
    "shortcode": "custom"
  }
  ```
  - `url` (string, required): The original long URL to be shortened
  - `validity` (integer, optional): The duration in minutes for which the short link remains valid (defaults to 30 minutes)
  - `shortcode` (string, optional): A desired custom shortcode (if omitted, a unique shortcode will be generated)

- **Response** (Status Code: 201):
  ```json
  {
    "shortLink": "http://hostname:port/custom",
    "expiry": "2023-05-01T12:30:00Z"
  }
  ```

### Retrieve Short URL Statistics

- **Method**: GET
- **Route**: `/shorturls/:shortcode`
- **Response**:
  ```json
  {
    "shortcode": "custom",
    "originalUrl": "https://example.com/very-long-url",
    "createdAt": "2023-05-01T12:00:00Z",
    "expiresAt": "2023-05-01T12:30:00Z",
    "clicks": 5,
    "clickData": [
      {
        "timestamp": "2023-05-01T12:05:00Z",
        "referrer": "https://referrer.com",
        "location": "Location from 192.168.1.1",
        "userAgent": "Mozilla/5.0 ..."
      }
    ]
  }
  ```

### Redirect to Original URL

- **Method**: GET
- **Route**: `/:shortcode`
- **Behavior**: Redirects to the original URL associated with the shortcode

## Error Handling

The API returns appropriate HTTP status codes and descriptive JSON responses for various error scenarios:

- **400 Bad Request**: Invalid input parameters
- **404 Not Found**: Shortcode not found
- **409 Conflict**: Shortcode already exists
- **410 Gone**: Shortcode has expired
- **500 Internal Server Error**: Server-side errors

## Running the Service

1. Clone the repository
2. Build the application:
   ```
   go build -o url-shortener
   ```
3. Run the service:
   ```
   ./url-shortener
   ```
   
The service will start on port 8080 by default. You can change the port by setting the `PORT` environment variable.

## Design Considerations

- **In-Memory Storage**: The current implementation uses in-memory storage for simplicity. In a production environment, this would be replaced with a persistent database.
- **Concurrency**: The service is designed to be thread-safe with proper mutex locking in the storage layer.
- **Asynchronous Click Recording**: Click events are recorded asynchronously to ensure fast redirections.
- **Shortcode Generation**: Random shortcodes are generated using cryptographically secure random number generation.
- **Logging**: Extensive logging is implemented throughout the application to track operations and errors.