/* {{{ Copyright (c) Paul R. Tagliamonte <paultag@debian.org>, 2015
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE. }}} */
package config

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path"
	"reflect"
	"strings"
	"unsafe"

	"pault.ag/go/debian/control"
)

func LoadFlags(name string, data interface{}) (*flag.FlagSet, error) {
	if err := Load(name, data); err != nil {
		return nil, err
	}
	return Flag(data)
}

//
func Load(name string, data interface{}) error {
	localUser, err := user.Current()
	if err != nil {
		return nil
	}
	rcPath := path.Join(localUser.HomeDir, fmt.Sprintf(".%src", name))
	fd, err := os.Open(rcPath)
	if err != nil {
		return nil
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
	name := strings.ToLower(dataValue.Elem().Type().Name())
	flagSet := flag.NewFlagSet(name, flag.ExitOnError)
	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", name)
		flagSet.PrintDefaults()
	}

	err := flagPointer(dataValue, flagSet)
	if err != nil {
		return nil, err
	}
	return flagSet, nil
}

// vim: foldmethod=marker
