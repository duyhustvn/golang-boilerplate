#include "unity.h"
#include "circle_linkedlist.h"
#include <stdlib.h>

void setUp(void) {}
void tearDown(void) {}

void test_init_cll(void) {
    char errstr[1024];
    CircleLinkedList *list = init_cll(errstr, sizeof(errstr));

    TEST_ASSERT_NOT_NULL(list);
    TEST_ASSERT_EQUAL_INT(list->length, 0);
    TEST_ASSERT_NULL(list->head);

    free_cll(list);
}

void test_insert_only_one_into_cll(void) {
    char errstr[1024];
    CircleLinkedList *list = init_cll(errstr, sizeof(errstr));

    int *i;
    int *data;

    i = malloc(sizeof(int));
    *i = 1;
    insert_cll(list, i, errstr, sizeof(errstr));
    TEST_ASSERT_EQUAL_INT(list->length, 1);
    data = (int *) (list->head->data);
    TEST_ASSERT_EQUAL_INT(*data, 1);
    data = (int *) (list->head->next->data);
    TEST_ASSERT_EQUAL_INT(*data, 1);

    free_cll(list);
}

void test_insert_multiple_into_cll(void) {
    char errstr[1024];
    CircleLinkedList *list = init_cll(errstr, sizeof(errstr));

    int *i;
    int *data;

    i = malloc(sizeof(int));
    *i = 1;
    insert_cll(list, i, errstr, sizeof(errstr));
    TEST_ASSERT_EQUAL_INT(list->length, 1);

    // List has 2 elements:
    // 1 -> 2 (head)
    i = malloc(sizeof(int));
    *i = 2;
    insert_cll(list, i, errstr, sizeof(errstr));
    TEST_ASSERT_EQUAL_INT(list->length, 2);

    data = (int *)(list->head->data);
    TEST_ASSERT_EQUAL_INT(*data, 2);
    data = (int *)(list->head->next->data);
    TEST_ASSERT_EQUAL_INT(*data, 1);
    data = (int *)(list->head->next->next->data);
    TEST_ASSERT_EQUAL_INT(*data, 2);
    data = (int *)(list->head->next->next->next->data);
    TEST_ASSERT_EQUAL_INT(*data, 1);

    // List has 3 elements:
    // 1 -> 2 -> 3 (head)
    i = malloc(sizeof(int));
    *i = 3;
    insert_cll(list, i, errstr, sizeof(errstr));
    TEST_ASSERT_EQUAL_INT(list->length, 3);

    data = (int *)(list->head->data);
    TEST_ASSERT_EQUAL_INT(*data, 3);
    data = (int *)(list->head->next->data);
    TEST_ASSERT_EQUAL_INT(*data, 1);
    data = (int *)(list->head->next->next->data);
    TEST_ASSERT_EQUAL_INT(*data, 2);
    data = (int *)(list->head->next->next->next->data);
    TEST_ASSERT_EQUAL_INT(*data, 3);

    // List has 4 elements:
    // 1 -> 2 -> 3 -> 4 (head)
    i = malloc(sizeof(int));
    *i = 4;
    insert_cll(list, i, errstr, sizeof(errstr));
    TEST_ASSERT_EQUAL_INT(list->length, 4);

    data = (int *)(list->head->data);
    TEST_ASSERT_EQUAL_INT(*data, 4);
    data = (int *)(list->head->next->data);
    TEST_ASSERT_EQUAL_INT(*data, 1);
    data = (int *)(list->head->next->next->data);
    TEST_ASSERT_EQUAL_INT(*data, 2);
    data = (int *)(list->head->next->next->next->data);
    TEST_ASSERT_EQUAL_INT(*data, 3);
    data = (int *)(list->head->next->next->next->next->data);
    TEST_ASSERT_EQUAL_INT(*data, 4);

    free_cll(list);
}
