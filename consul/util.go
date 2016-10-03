// Package consul provides utility functions for the consul key/value store.
package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

// internalGet returns nil if key does not exist or if
// there is an error during lookup. In both cases a
// message is logged.
func internalGet(kv *api.KV, key string) []byte {

	kvp, _, err := kv.Get(key, nil)

	if err != nil {
		fmt.Printf("Error while reading key >%s<. Returning nil as key value.\n", err)
		return nil
	} else if kvp == nil {
		fmt.Printf("Key >%s< does not exist in registry. Returning nil as key value.", key)
		return nil
	}

	return kvp.Value
}
