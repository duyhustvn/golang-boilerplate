#include "circle_linkedlist.h"
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>


CllNode *create_cll_node(void *data, char *errstr, size_t errstr_size) {
    CllNode *node = (CllNode *) malloc(sizeof(CllNode));
    if (!node) {
        snprintf(errstr, errstr_size, "Failed to init circle linked list node");
        return NULL;
    }
    node->data = data;
    node->next = NULL;
    return node;
}

void free_cll_node(CllNode *node) {
    if (!node) {
        return;
    }

    free(node->data);
    free(node);
}



CircleLinkedList *init_cll(char *errstr, size_t errstr_size) {
    CircleLinkedList *list = (CircleLinkedList *) malloc(sizeof(CircleLinkedList));
    if (!list) {
        snprintf(errstr, errstr_size, "Failed to init circle linked list");
        return NULL;
    }

    list->length = 0;
    pthread_mutex_init(&list->mutex, NULL);
    return list;
};

bool insert_cll(CircleLinkedList *list, void *data, char *errstr, size_t errstr_size) {
    pthread_mutex_lock(&list->mutex);

    CllNode *node = create_cll_node(data, errstr, errstr_size);
    if (!node) {
        return NULL;
    }

    if (!list->head) {
        // list is empty
        list->head = node;
        list->head->next = node;
    } else {
        // list has element
        node->next = list->head->next;
        list->head->next = node;
        list->head = node;
    }

    list->length++;

    pthread_mutex_unlock(&list->mutex);
    return NULL;
};


void free_cll(CircleLinkedList *list) {
    if (!list || list->length == 0) {
        return;
    }

    free(list);
};
