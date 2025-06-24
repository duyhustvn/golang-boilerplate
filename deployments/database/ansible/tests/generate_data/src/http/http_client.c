#include "http_client.h"
#include <string.h>
#include <stdlib.h>

// Internal: Write callback for libcurl
static size_t write_callback(void* contents, size_t size, size_t nmemb, void* userp) {
    size_t realsize = size * nmemb;
    HttpResponse* response = (HttpResponse*)userp;
    
    char* ptr = realloc(response->data, response->size + realsize + 1);
    if (!ptr) return 0;  // Out of memory
    
    response->data = ptr;
    memcpy(&(response->data[response->size]), contents, realsize);
    response->size += realsize;
    response->data[response->size] = '\0';
    
    return realsize;
}

// Internal: Convert method enum to string
static const char* method_to_string(HttpMethod method) {
    switch (method) {
        case HTTP_GET:    return "GET";
        case HTTP_POST:   return "POST";
        case HTTP_PUT:    return "PUT";
        case HTTP_DELETE: return "DELETE";
        case HTTP_PATCH:  return "PATCH";
        default:          return "GET";
    }
}

void http_client_init() {
    curl_global_init(CURL_GLOBAL_DEFAULT);
}

void http_client_cleanup() {
    curl_global_cleanup();
}

HttpResponse http_request(const char* url,
                          HttpMethod method,
                          const char* body,
                          const char* content_type,
                          const char* headers[],
                          long timeout_sec) {
    HttpResponse response = {0};
    CURL* curl = curl_easy_init();
    
    if (!curl) {
        response.curl_code = CURLE_FAILED_INIT;
        return response;
    }

    // Set basic options
    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_callback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response);
    curl_easy_setopt(curl, CURLOPT_USERAGENT, "http_client/1.0");
    curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, timeout_sec);

    // Configure method-specific options
    const char* method_str = method_to_string(method);
    if (method == HTTP_POST) {
        curl_easy_setopt(curl, CURLOPT_POST, 1L);
    } else if (method != HTTP_GET) {
        curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, method_str);
    }

    // Set request body if provided
    if (body) {
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, body);
    }

    // Build headers list
    struct curl_slist* header_list = NULL;
    
    // Set Content-Type if specified
    if (content_type) {
        char content_type_header[128];
        snprintf(content_type_header, sizeof(content_type_header), 
                 "Content-Type: %s", content_type);
        header_list = curl_slist_append(header_list, content_type_header);
    }
    
    // Add custom headers
    if (headers) {
        for (int i = 0; headers[i]; i++) {
            header_list = curl_slist_append(header_list, headers[i]);
        }
    }
    
    if (header_list) {
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, header_list);
    }

    // Execute request
    response.curl_code = curl_easy_perform(curl);
    
    // Get HTTP status code
    if (response.curl_code == CURLE_OK) {
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &response.status_code);
    }

    // Cleanup
    curl_slist_free_all(header_list);
    curl_easy_cleanup(curl);
    
    return response;
}

void http_response_free(HttpResponse* response) {
    if (response) {
        free(response->data);
        response->data = NULL;
        response->size = 0;
    }
}
