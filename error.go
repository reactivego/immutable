package immutable

type MapError string

func (e MapError) Error() string {
	return string(e)
}

const UnhashableKeyType = MapError("Unhashable Key Type")
