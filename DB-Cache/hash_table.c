#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <pthread.h>

#include "hash_table.h"

// Make sure this value is a power of two since we're using & instead of %
#define INDEX_COUNT 1048576

struct Q_Node {
	struct Q_Node* next;
	int64_t data;
};

pthread_mutex_t lock;

struct Q_Node* table[INDEX_COUNT] = { 0 };

int64_t Q_traverse_insert(struct Q_Node** node, int64_t data)
{
	while (1) { // This loops if another thread steals its insertion spot
		while (*node) {
			if ((*node)->data == data)
				return 0;
			node = &(*node)->next;
		}
		pthread_mutex_lock(&lock);
		if (*node == NULL) { // make sure another thread didn't steal its spot
			*node = calloc(sizeof(struct Q_Node), 1);
			(*node)->data = data;
			pthread_mutex_unlock(&lock);
			return 1;
		}
	pthread_mutex_unlock(&lock);
	}
}

int64_t table_insert(int64_t data)
{
	return Q_traverse_insert(&table[(INDEX_COUNT - 1) & data], data);
}
