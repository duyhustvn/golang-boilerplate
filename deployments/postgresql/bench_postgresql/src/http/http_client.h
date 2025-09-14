#ifndef HTTP_CLIENT_H
#define HTTP_CLIENT_H

#include <curl/curl.h>

// HTTP Methods
typedef enum {
    HTTP_GET,
    HTTP_POST,
    HTTP_PUT,
    HTTP_DELETE,
    HTTP_PATCH
} HttpMethod;

// HTTP Response structure
typedef struct {
    char* data;         // Response body (NULL-terminated)
    size_t size;        // Size of response body
    long status_code;   // HTTP status code
    CURLcode curl_code; // libcurl error code
} HttpResponse;

// Initialize HTTP client - call once at program start
void http_client_init();

// Cleanup HTTP client - call once at program end
void http_client_cleanup();

// Execute HTTP request
HttpResponse http_request(
    const char* url,
    HttpMethod method,
    const char* body,           // Request body (NULL for no body)
    const char* content_type,   // Content-Type header (NULL for default)
    const char* headers[],      // Additional headers (NULL-terminated array)
    long timeout_sec            // Request timeout in seconds
);

// Free resources in HttpResponse
void http_response_free(HttpResponse* response);

#endif // HTTP_CLIENT_H
