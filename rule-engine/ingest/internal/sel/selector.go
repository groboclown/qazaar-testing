// Under the Apache-2.0 License
package sel

import "fmt"

type SelectHandler func(val map[string]any) error

type SelectHandlerMap map[string]SelectHandler

func TypeSelector(val any, key string, handlers SelectHandlerMap) error {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case map[string]any:
		if tg, ok := v[key]; ok {
			if t, ok := tg.(string); ok {
				if h, ok := handlers[t]; ok {
					return h(v)
				}
				return fmt.Errorf("unknown %s: '%s'", key, t)
			}
			return fmt.Errorf("no string for %s (%v)", key, tg)
		}
		return fmt.Errorf("no %s", key)
	}
	return fmt.Errorf("invalid type (%v)", val)
}
