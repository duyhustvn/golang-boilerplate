#include "queue.h"

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

void free_node(Node *node) {
    if (!node) {
        return;
    }

    free(node->data);
    free(node);
};

Queue *init_queue(char *errstr) {
    Queue *q = (Queue *) malloc(sizeof(Queue));
    if (!q) {
        errstr = "Failed to malloc data for queue";
        return NULL;
    }

    q->length = 0;
    q->front = NULL;
    q->rear = NULL;

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

Node *dequeue(Queue *q, char *errstr, size_t err_msg_size) {
    pthread_mutex_lock(&q->mutex);

    if (is_empty(q)) {
        snprintf(errstr, err_msg_size, "Queue is empty");
        pthread_mutex_unlock(&q->mutex);
        return NULL;
    }

    Node *node = q->front;
    if (q->length == 1) {
        q->front = NULL;
        q->rear = NULL;
    } else {
        q->front = q->front->next;
    }
    node->next = NULL;
    q->length--;
    pthread_mutex_unlock(&q->mutex);
    return node;
};


void free_queue(Queue *q) {
    if (!q) {
        return;
    }

    char errstr[1024];
    while (q->length > 0) {
        Node *node = dequeue(q, errstr, sizeof(errstr));
        free_node(node);
    }
    free(q);
};
