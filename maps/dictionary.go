package maps

type Dictionary map[string]string

const (
	ErrNotFound   = DictionaryErr("could not find the word you are looking for")
	ErrWordExists = DictionaryErr("the key already exists in the map with a value assigned")
	ErrWordDoesntExist = DictionaryErr("the key has no definition assigned")
)

type DictionaryErr string

func (e DictionaryErr) Error() string {
	return string(e)
}

func (d Dictionary) Search(key string) (string, error) {
	if val, ok := d[key]; ok {
		return val, nil
	}
	return "", ErrNotFound
}

func (d Dictionary) Add(key, value string) error {
	_, err := d.Search(key)
	switch err {
	case ErrNotFound:
		d[key] = value
	case nil:
		return ErrWordExists
	default:
		return err
	}
	return nil
}

func (d Dictionary) Update(key, value string) error {
	_, err := d.Search(key)
	switch err{
	case ErrNotFound:
		return ErrWordDoesntExist
	case nil:
		d[key] = value
	default:
		return err
	}
	return nil
}

func (d Dictionary) Delete(key string){
	delete(d, key)
}
