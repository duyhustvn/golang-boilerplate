#include "queue.h"
#include <pthread.h>
#include <stdio.h>


Node *create_node(void *data, char *errstr) {
    Node *node = (Node *) malloc(sizeof(Node));
    if (!node) {
        errstr = "Failed to malloc data for node";
        return NULL;
    }

    node->data = data;
    node->next = NULL;
    return node;
};

Queue *init_queue(char *errstr) {
    Queue *q = (Queue *) malloc(sizeof(Queue));
    if (!q) {
        errstr = "Failed to malloc data for queue";
        return NULL;
    }

    q->length = 0;

    pthread_mutex_init(&q->mutex, NULL);
    return q;
};


bool is_empty(Queue *q) {
    return q->front == NULL; 
};


bool enqueue(Queue *q, void *data, char *errstr) {
    Node *node = create_node(data, errstr);
    if (!node) {
        return false;
    }

    pthread_mutex_lock(&q->mutex);
    if (is_empty(q)) {
        q->front = node;         
        q->rear = node;
    } else {
        q->rear->next = node;
        q->rear = node;
    }

    q->length++;
    pthread_mutex_unlock(&q->mutex);
    return true;
};


void *dequeue(Queue *q, char *errstr) {
    if (is_empty(q)) {
        return NULL;
    }
    
    pthread_mutex_lock(&q->mutex);
    Node *node = q->front;
    if (q->length == 1) {
        q->front = NULL;
        q->rear = NULL;
    } else {
        q->front = q->front->next;
    }

    q->length--;
    pthread_mutex_unlock(&q->mutex);
    return node;
};
