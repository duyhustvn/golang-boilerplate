#include "signup_api.h"
#include "../http/http_client.h"
#include "../cJSON/cJSON.h"
#include <pthread.h>
#include <stdio.h>

char *signup_request_body_to_json_string(signup_request_body body) {
    cJSON *data = cJSON_CreateObject();

    cJSON_AddStringToObject(data, "application", body.application);
    cJSON_AddStringToObject(data, "organization", body.organization);
    cJSON_AddStringToObject(data, "username", body.username);
    cJSON_AddStringToObject(data, "password", body.password);

    char *json_string = cJSON_Print(data);
    cJSON_Delete(data);

    return json_string;
}

bool call_signup_api(signup_request_body body, char* err) {
    char *body_string = signup_request_body_to_json_string(body);
    if (!body_string) {
        err = "signup_request_body_to_json_string failed";
        return false;
    }
    printf("body_string: %s\n", body_string);    
    const char* headers[] = {
        "Accept: application/json",
        NULL, // terminator
    };

    int timeout_sec = 30; // second

    HttpResponse res = http_request(
        "http://localhost:8000/api/signup",
        HTTP_POST,
        body_string,
        "application/json",
        headers,
        timeout_sec
    );

    printf("Response:\n");
    printf("Status: %ld\n", res.status_code);
    printf("Body: %s\n", res.data);
    printf("cURL code: %d\n\n", res.curl_code);

    http_response_free(&res);
    return true;
};

void *call_signup_api_thread(void *thread_args) {
    sign_up_thread_vars *thread_vars = (sign_up_thread_vars *) thread_args;
    Queue *q = thread_vars->q;
    char errstr[256];
    Node *node = dequeue(q, errstr);
    if (!node) {
        printf("[call_signup_api_thread] %s", errstr);
        pthread_exit(errstr);
    }

    int *i = (int *) node->data;

    char username[15];     
    sprintf(username, "%04d", *i);
    printf("username: %s\n", username);

    signup_request_body body = {
        .application = "application_bench_psql_cluster",
        .organization = "organization_bench_psql_cluster",
        .username = "user_0.000.000.004",
        .password = "randompasswd"
    };

    bool ok = call_signup_api(body, errstr);
    if (!ok) {
        printf("[call_signup_api_thread] %s", errstr);
        pthread_exit(errstr);
    }

    pthread_exit(NULL);
}
