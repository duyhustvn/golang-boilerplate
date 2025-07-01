#include "login_api.h"
#include "../http/http_client.h"
#include "../cJSON/cJSON.h"

char *login_request_body_to_json_string(login_request_body body) {
    cJSON *data = cJSON_CreateObject();

    cJSON_AddStringToObject(data, "grant_type", body.grant_type);
    cJSON_AddStringToObject(data, "client_id", body.client_id);
    cJSON_AddStringToObject(data, "username", body.username);
    cJSON_AddStringToObject(data, "password", body.password);

    char *json_string = cJSON_Print(data);
    cJSON_Delete(data);

    return json_string;
}

bool call_login_api(login_request_body body, char* errstr, size_t errstr_size) {
    char *body_string = login_request_body_to_json_string(body);
    if (!body_string) {
        snprintf(errstr, errstr_size, "login_request_body_to_json_string failed");
        return false;
    }
    printf("body_string: %s\n", body_string);    
    const char* headers[] = {
        "Accept: application/json",
        NULL, // terminator
    };

    int timeout_sec = 30; // second

    HttpResponse res = http_request(
        "http://localhost:8000/api/login/oauth/access_token",
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

void *call_login_api_thread(void *thread_args) {
    login_thread_vars *thread_vars = (login_thread_vars *) thread_args; 
    Queue *q = thread_vars->q;

    while (true) {
        char errstr[256];
        Node *node = dequeue(q, errstr, sizeof(errstr));
        if (!node) {
            printf("[call_login_api_thread] error %s \n", errstr);
            break;
        }

        int *i = (int *) node->data;

        char username[20];     
        snprintf(username, sizeof(username), "user_%013d", *i); // %0 will automatically append 0 to the beginning of the string if it is shorter than 13
        printf("threadID: %d, username: %s\n", thread_vars->threadId, username);

        
        login_request_body body = {
            .grant_type = "password",
            .client_id = "6c875db314d6760e3b69",
            .username = username,
            .password = "randompasswd"
        };

        bool ok = call_login_api(body, errstr, sizeof(errstr));
        if (!ok) {
            printf("[call_login_api_thread] %s", errstr);
            pthread_exit(NULL);
        }
    }

    pthread_exit(NULL);
} 
