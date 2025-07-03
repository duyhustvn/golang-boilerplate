#ifndef LOGIN_API_H
#define LOGIN_API_H

#include "../common.h"
#include "../queue/queue.h"

/* 
Endpoint: /api/login/oauth/access_token 
Headers:
    "content-type": "application/json"
Body:  
    {
        "grant_type": "password",
        "client_id": "6c875db314d6760e3b69",
        "username": "user_0.000.000.002",
        "password": "randompasswd"
    }
*/

typedef struct login_thread_vars_ {
    int threadId;
    Queue *q;   
} login_thread_vars;

typedef struct login_request_body_ {
   char* grant_type; 
   char* client_id;
   char* username;
   char* password;
   char *country_code;
} login_request_body;


char *login_request_body_to_json_string(login_request_body body);
bool call_login_api(login_request_body body, char* errstr, size_t errstr_size);
void *call_login_api_thread(void *thread_args);

#endif // LOGIN_API_H
