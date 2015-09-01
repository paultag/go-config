package config

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path"
	"reflect"
	"unsafe"

	"pault.ag/go/debian/control"
)

//
func Load(name string, data interface{}) error {
	localUser, err := user.Current()
	if err != nil {
		return err
	}
	rcPath := path.Join(localUser.HomeDir, fmt.Sprintf(".%src", name))
	fd, err := os.Open(rcPath)
	if err != nil {
		return err
	}
	defer fd.Close()
	err = control.Unmarshal(data, fd)
	return err
}

func flagPointer(incoming reflect.Value, data *flag.FlagSet) error {
	if incoming.Type().Kind() == reflect.Ptr {
		return flagPointer(incoming.Elem(), data)
	}

	for i := 0; i < incoming.NumField(); i++ {
		field := incoming.Field(i)
		fieldType := incoming.Type().Field(i)

		if it := fieldType.Tag.Get("flag"); it != "" {
			/* Register the flag */
			switch field.Type().Kind() {
			case reflect.Int:
				data.IntVar(
					(*int)(unsafe.Pointer(field.Addr().Pointer())),
					it,
					int(field.Int()),
					fieldType.Tag.Get("description"),
				)
				continue
			case reflect.String:
				data.StringVar(
					(*string)(unsafe.Pointer(field.Addr().Pointer())),
					it,
					field.String(),
					fieldType.Tag.Get("description"),
				)
				continue
			default:
				return fmt.Errorf("Unknown type: %s", field.Type().Kind())
			}
		}
	}
	return nil
}

func Flag(data interface{}) (*flag.FlagSet, error) {
	dataValue := reflect.ValueOf(data)
	flagSet := flag.NewFlagSet(dataValue.Type().Name(), flag.ExitOnError)
	err := flagPointer(dataValue, flagSet)
	if err != nil {
		return nil, err
	}
	return flagSet, nil
}
