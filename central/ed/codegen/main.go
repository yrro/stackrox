// Package main generates the code for (scrambled) keys.
// This is ugly, but that's ok. It is not intended to be general/reusable.
package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
)

func writeToFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	_, err = f.Write(data)
	return err
}

func writeValueGetter(name string, value []byte, out io.Writer) error {
	_, err := fmt.Fprintf(out,
		`func %s() []byte {
	var result [%d]byte
`, name, len(value))
	if err != nil {
		return err
	}

	for i, val := range value {
		var err error
		if i == 0 {
			_, err = fmt.Fprintf(out,
				"	result[%d] = %#02x ^ %#02x\n", i, ^byte(i), val^^byte(i))
		} else {
			_, err = fmt.Fprintf(out, "	result[%d] = result[%d] ^ %#02x ^ %#02x\n", i, i-1, ^byte(i), val^^byte(i)^value[i-1])
		}
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(out, `	return result[:]
}`)
	return err
}

func writeGeneratedCode(key, iv []byte, out io.Writer) error {
	_, err := fmt.Fprintln(out, `// Code generate by ed/codegen. DO NOT EDIT.

package ed`)
	if err != nil {
		return err
	}

	if err := writeValueGetter("k", key[:], out); err != nil {
		return err
	}
	if err := writeValueGetter("i", iv[:], out); err != nil {
		return err
	}

	return nil
}

func main() {
	key := make([]byte, 32)
	iv := make([]byte, 16)

	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err)
	}
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	if err := writeToFile("../../image/keys/data-key", key); err != nil {
		panic(err)
	}
	if err := writeToFile("../../image/keys/data-iv", iv); err != nil {
		panic(err)
	}

	outFile, err := os.OpenFile("gen-keys.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = outFile.Close()
	}()
	if err := writeGeneratedCode(key, iv, outFile); err != nil {
		panic(err)
	}
}
