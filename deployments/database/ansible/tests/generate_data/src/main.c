#include "common.h"

#include "http/http_client.h"
#include "api_caller/signup_api.h"
#include "queue/queue.h"
#include <stdio.h>

#define NUM_THREADS 8


int main() {
    // Initialize client (call once at program start)
    http_client_init();

    char errstr[1024];
    Queue *q = init_queue(errstr);
    if (!q) {
        printf("%s\n", errstr);
    }

    int *v;
    for (int i = 0; i < 100; i++) {
        v = malloc(sizeof(int));
        *v = i;
        bool ok = enqueue(q, v, errstr);
        if (!ok) {
            printf("%s\n", errstr);
            continue;
        }
    }

    if (is_empty(q)) {
        return 0;
    }

    sign_up_thread_vars thread_vars = {
        .q = q,
    };

    pthread_t threads[NUM_THREADS];
    for (int i = 0; i < NUM_THREADS; i++) {
        pthread_create(&threads[i], NULL, call_signup_api_thread, (void *)&thread_vars);
    }
     
    // Cleanup client (call once at program end)
    http_client_cleanup();
    return 0;
}
