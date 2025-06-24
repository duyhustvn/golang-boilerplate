#include "http/http_client.h"
#include "signup/signup.h"

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>  

int main() {
    // Initialize client (call once at program start)
    http_client_init();

    char err[1024];
    
    signup_request_body body = {
        .application = "application_bench_psql_cluster",
        .organization = "organization_bench_psql_cluster",
        .username = "user_0.000.000.004",
        .password = "randompasswd"
    };

    bool ok = call_signup_api(body, err);
    if (!ok) {
        printf("%s \n", err);
        return EXIT_FAILURE;
    }
     
    // Cleanup client (call once at program end)
    http_client_cleanup();
    return 0;
}
