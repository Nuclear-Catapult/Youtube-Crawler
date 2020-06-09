#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <pthread.h>

struct BST_Node {
	struct BST_Node* left;
	struct BST_Node* right;
	uint64_t data;
};

struct BST_Node* root = NULL;

pthread_mutex_t key;

// returns 1 if success
// returns 0 if duplicate value caused failure
uint64_t BST_insert(uint64_t data)
{
	struct BST_Node** node = &root;
	pthread_mutex_lock(&key);
	while (*node) {
		if ((*node)->data == data) {
			pthread_mutex_unlock(&key);
			return 0;
		}
		node = ((*node)->data > data) ? &(*node)->left : &(*node)->right;
	}
	*node = calloc(sizeof(struct BST_Node), 1);
	(*node)->data = data;

	pthread_mutex_unlock(&key);

	return 1;
}
