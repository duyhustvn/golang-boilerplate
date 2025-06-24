#ifndef SIGNUP_H
#define SIGNUP_H

#include <stdbool.h>
#include "../cJSON/cJSON.h"

typedef struct signup_request_body_ {
   char* application; 
   char* organization;
   char* username;
   char* password;
} signup_request_body;

char *signup_request_body_to_json_string(signup_request_body body);
bool call_signup_api(signup_request_body body, char* err);

#endif // SIGNUP_H
