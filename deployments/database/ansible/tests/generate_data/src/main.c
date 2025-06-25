#include "http/http_client.h"
#include "signup/signup.h"

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>  
#include <pthread.h>

#define NUM_THREADS 8

void *call_api(void* threadArg) {

}

int main() {
    // Initialize client (call once at program start)
    http_client_init();

    char err[1024];
    pthread_t threads[NUM_THREADS];

    // bool ok = call_signup_api(body, err);
    // if (!ok) {
    //     printf("%s \n", err);
    //     return EXIT_FAILURE;
    // }
    
    

    for (int i = 0; i < NUM_THREADS; i++) {
        char username[15];

        sprintf(username, "%04d", i);
        printf("username: %s \n", username);

        signup_request_body body = {
              .application = "application_bench_psql_cluster",
              .organization = "organization_bench_psql_cluster",
              .username = "user_0.000.000.004",
              .password = "randompasswd"
          };

        pthread_create(threads[i], NULL, );
    }
     
    // Cleanup client (call once at program end)
    http_client_cleanup();
    return 0;
}
