package data

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/camelcase"
)

// Decoder decodes a Map into a struct.
type Decoder struct {
	config *DecoderConfig

	// This Decoder is strongly inspired by github.com/mitchellh/mapstructure.
	// The decoder may be going to support "json" compatible tags in structs.
}

// DecoderConfig is used to configure the behavior of Decoder.
type DecoderConfig struct {
	// ErrorUnused, if set to true, reports an error when a Map contains keys
	// that are not defined in a struct.
	ErrorUnused bool

	// Metadata has meta information of decode. If this is nil, meat information
	// will not be tracked.
	Metadata *DecoderMetadata

	TagName string
}

// DecoderMetadata tracks field names that are used or not used for decoding.
type DecoderMetadata struct {
	// Keys contains keys in a Map that are processed.
	Keys []string

	// Unsed contains keys in a Map that are not defined in a struct.
	Unused []string
}

// NewDecoder creates a new Decoder with the given config.
func NewDecoder(c *DecoderConfig) *Decoder {
	if c == nil {
		c = &DecoderConfig{}
	}
	if c.TagName == "" {
		c.TagName = "bql"
	}
	return &Decoder{
		config: c,
	}
}

// Decode decodes a Map into a struct. The argument must be a pointer to a
// struct.
func (d *Decoder) Decode(m Map, v interface{}) error {
	p := reflect.ValueOf(v)
	if p.Kind() != reflect.Ptr {
		return errors.New("result must be a pointer to a struct")
	}
	s := p.Elem()
	if s.Kind() != reflect.Struct {
		return errors.New("result must be pointer to a struct")
	}
	return d.decodeStruct(m, s)
}

func (d *Decoder) decode(src Value, dst reflect.Value, weaklyTyped bool) error {
	switch dst.Kind() {
	case reflect.Bool:
		return d.decodeBool(src, dst, weaklyTyped)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// TODO: support detecting overflows
		return d.decodeInt(src, dst, weaklyTyped)

	case reflect.Float32, reflect.Float64:
		return d.decodeFloat(src, dst, weaklyTyped)

	case reflect.String:
		return d.decodeString(src, dst, weaklyTyped)

	case reflect.Interface: // Only interface{} and Value is supported
		if !reflect.TypeOf(func(interface{}) {}).In(0).AssignableTo(dst.Type()) {
			return errors.New("only empty interface{} is supported")
		}
		return d.decodeInterface(src, dst)

	case reflect.Map:
		return d.decodeMap(src, dst, weaklyTyped)

	case reflect.Slice:
		return d.decodeSlice(src, dst, weaklyTyped)

	case reflect.Struct:
		return d.decodeStruct(src, dst)
	}
	return fmt.Errorf("decoder doesn't support the type: %v", dst.Kind())
}

func (d *Decoder) decodeBool(src Value, dst reflect.Value, weaklyTyped bool) error {
	var (
		b   bool
		err error
	)
	if weaklyTyped {
		b, err = ToBool(src)
	} else {
		b, err = AsBool(src)
	}
	if err != nil {
		return err
	}
	dst.SetBool(b)
	return nil
}

func (d *Decoder) decodeInt(src Value, dst reflect.Value, weaklyTyped bool) error {
	var (
		i   int64
		err error
	)

	// TODO: support time.Duration

	if weaklyTyped {
		i, err = ToInt(src)
	} else {
		i, err = AsInt(src)

		// special handling for floats having integer values
		if err != nil && src.Type() == TypeFloat {
			f, _ := AsFloat(src)
			i = int64(f)
			if f == float64(i) {
				err = nil
			}
		}
	}
	if err != nil {
		return err
	}
	dst.SetInt(i)
	return nil
}

func (d *Decoder) decodeFloat(src Value, dst reflect.Value, weaklyTyped bool) error {
	var (
		f   float64
		err error
	)
	if weaklyTyped {
		f, err = ToFloat(src)
	} else {
		f, err = AsFloat(src)
	}
	if err != nil {
		return err
	}
	dst.SetFloat(f)
	return nil
}

func (d *Decoder) decodeString(src Value, dst reflect.Value, weaklyTyped bool) error {
	var (
		s   string
		err error
	)
	if weaklyTyped {
		s, err = ToString(src)
	} else {
		s, err = AsString(src)
	}
	if err != nil {
		return err
	}
	dst.SetString(s)
	return nil
}

func (d *Decoder) decodeInterface(src Value, dst reflect.Value) error {
	return errors.New("not implemented yet")
}

func (d *Decoder) decodeMap(src Value, dst reflect.Value, weaklyTyped bool) error {
	if src.Type() != TypeMap {
		return fmt.Errorf("cannot decode to a map: %v", src.Type())
	}
	m, _ := AsMap(src)

	t := dst.Type()
	if k := t.Key().Kind(); k != reflect.String {
		return fmt.Errorf("key must be string: %v", k)
	}
	valueType := t.Elem()

	res := reflect.MakeMap(t)
	for k, e := range m {
		v := reflect.Indirect(reflect.New(valueType))
		if err := d.decode(e, v, weaklyTyped); err != nil {
			// TODO: this should probably be multierror, too.
			return err
		}
		res.SetMapIndex(reflect.ValueOf(k), v)
	}
	dst.Set(res)
	return nil
}

func (d *Decoder) decodeSlice(src Value, dst reflect.Value, weaklyTyped bool) error {
	if src.Type() != TypeArray {
		return fmt.Errorf("cannot decode to an array: %v", src.Type())
	}
	a, _ := AsArray(src)

	res := reflect.MakeSlice(dst.Type(), len(a), len(a))
	for i, e := range a {
		v := res.Index(i)
		if err := d.decode(e, v, weaklyTyped); err != nil {
			// TODO: this should probably be multierror, too.
			return err
		}
	}
	dst.Set(res)
	return nil
}

func (d *Decoder) decodeStruct(src Value, dst reflect.Value) error {
	// TODO: support time.Time

	m, err := AsMap(src)
	if err != nil {
		return errors.New("struct can only be decoded from a map")
	}

	// Aggregates all error informations to help users debug BQL.
	// TODO: replace this with github.com/hashicorp/go-multierror
	var errs []error

	t := dst.Type()
	for i, n := 0, t.NumField(); i < n; i++ {
		f := t.Field(i)
		opts := strings.Split(f.Tag.Get(d.config.TagName), ",")

		// parse options
		var (
			required    bool
			weaklyTyped bool
		)
		for ti, opt := range opts {
			if ti == 0 { // skip name
				continue
			}
			switch opt {
			case "required":
				required = true
			case "weaklytyped":
				weaklyTyped = true
			default:
				errs = append(errs, fmt.Errorf("%v field has an undefined option: %v", f.Name, opt))
			}
		}

		name := strings.TrimSpace(opts[0])
		if name == "" {
			name = toSnakeCase(f.Name)
		}
		src, ok := m[name]
		if !ok {
			if required {
				errs = append(errs, fmt.Errorf("%v is required but missing", name))
			}
			continue
		}

		if err := d.decode(src, dst.Field(i), weaklyTyped); err != nil {
			errs = append(errs, err)
		}
	}
	if errs != nil {
		// TODO: flatten errors
		return errs[0]
	}
	return nil
}

func toSnakeCase(name string) string {
	words := camelcase.Split(name)
	buf := bytes.NewBuffer(nil)
	for i, w := range words {
		if i > 0 {
			buf.WriteString("_")
		}
		buf.WriteString(strings.ToLower(w))
	}
	return string(buf.Bytes())
}

// Decode decodes a Map into a struct. The argument must be a pointer to a
// struct.
func Decode(m Map, v interface{}) error {
	return NewDecoder(nil).Decode(m, v)
}
