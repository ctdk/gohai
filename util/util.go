package util

// MergeMap merges two map[string]interface{} structures together. It's useful
// for adding collected data to the overall map of data to be converted to JSON.
func MergeMap(dst, src map[string]interface{}) error {
	for k, v := range src {
		if _, ok := dst[k]; !ok {
			// easiest situation; the destination doesn't have this
			// key.
			dst[k] = v
			continue
		} else {
			// otherwise it's more complicated
			err := merginate(dst[k], v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func merginate(dst, src interface{}) error {
	switch dst := dst.(type) {
		case map[string]interface{}:
			// but what is src?
			switch src := src.(type) {
				case map[string]interface{}:
					err := MergeMap(dst, src)
					if err != nil {
						return err
					}
				default:
					; // others later though?
			}
		case nil:
			dst = src
	}
	return nil
}
