#ifndef SIGNUP_API_H
#define SIGNUP_API_H

#include "../common.h"
#include "../queue/queue.h"

typedef struct sign_up_thread_vars_ {
    Queue *q;   
} sign_up_thread_vars;

typedef struct signup_request_body_ {
   char* application; 
   char* organization;
   char* username;
   char* password;
} signup_request_body;

char *signup_request_body_to_json_string(signup_request_body body);
bool call_signup_api(signup_request_body body, char* err);
void *call_signup_api_thread(void *thread_args);

#endif // SIGNUP_API_H
