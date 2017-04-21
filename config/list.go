package config

// List represents single or multi strings.
type List struct {
	Values []string
}

// UnmarshalYAML unmarshal single or multi strings.
func (l *List) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		l.Values = make([]string, 1)
		l.Values[0] = single
	} else {
		l.Values = multi
	}
	return nil
}
