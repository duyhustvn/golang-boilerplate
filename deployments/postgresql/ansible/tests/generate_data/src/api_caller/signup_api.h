#ifndef SIGNUP_API_H
#define SIGNUP_API_H

#include "../common.h"
#include "../queue/queue.h"

/* 
Endpoint: /api/signup
Headers:
    "content-type": "application/json"
Body:  
    {
        "application": "application_bench_psql_cluster",
        "organization": "organization_bench_psql_cluster",
        "username": "user_0.000.000.001",
        "password": "randompasswd"
    }
*/

typedef struct sign_up_thread_vars_ {
    int threadId;
    Queue *q;   
} sign_up_thread_vars;

typedef struct signup_request_body_ {
   char* application; 
   char* organization;
   char* username;
   char* password;
} signup_request_body;

char *signup_request_body_to_json_string(signup_request_body body);
bool call_signup_api(signup_request_body body, char* errstr, size_t errstr_size);
void *call_signup_api_thread(void *thread_args);

#endif // SIGNUP_API_H
