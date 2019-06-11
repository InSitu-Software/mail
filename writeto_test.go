package mail

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"
	"unicode/utf8"
)

type rot13 struct {
	PlainMessage     *Message
	EncryptedMessage strings.Builder
}

func rotN(c int32, shift int32) rune {
	return rune(c + shift)
}

func rotStringN(s string, shift int32) string {
	var out string
	for _, c := range s {
		out += fmt.Sprintf("%c", rotN(c, shift))
	}

	return out
}

func (r *rot13) String() string {
	return r.EncryptedMessage.String()
}

func (r *rot13) Write(b []byte) (n int, err error) {
	r.EncryptedMessage.Grow(len(b))

	for len(b) > 0 {
		c, size := utf8.DecodeRune(b)
		c = rune(c + 13)
		rn, err := r.EncryptedMessage.WriteRune(c)
		if err != nil {
			return n, err
		}
		b = b[size:]

		n += rn
	}

	return n, nil
}

func (r *rot13) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer
	_, err := r.PlainMessage.WriteTo(&buf)
	if err != nil {
		log.Fatal(err)
	}

	var written int64

	for _, c := range buf.String() {
		// naive implementation, no real rotation, only shifting, ignoring UTF-8 / rune sizes
		shifted := rotN(c, 13)
		bytedRune := byte(shifted)
		i, err := w.Write([]byte{bytedRune})
		written += int64(i)
		if err != nil {
			return written, err
		}

	}

	return written, nil
}

func TestRot13(t *testing.T) {
	s := "testá\u0061\u0300"
	e := rotStringN(s, 13)
	d := rotStringN(e, -13)

	if d != s {
		fmt.Printf("%s | %s | %s \n", s, e, d)
		fmt.Printf("%+q | %+q | %+q \n", s, e, d)
		t.Fail()
	}
}

func TestWriteTo(t *testing.T) {
	to := []string{"test@test.de"}
	from := "sender@sender.de"

	var b bytes.Buffer
	b.WriteString("test")

	secretMail := NewMessage()
	secretMail.SetHeader("To", to...)
	secretMail.SetHeader("From", from)
	secretMail.SetHeader("Subject", "my secret mail subject")
	secretMail.SetBody("text/plain", "my darkest secret is *#12!4//(+")

	rot := rot13{
		PlainMessage: secretMail,
	}

	secretMail.SetEnryption("application/pgp-encrypted", &rot)

	var out bytes.Buffer
	if _, err := secretMail.WriteTo(&out); err != nil {
		log.Fatal(err)
	}

	var outString string
	outString = out.String()

	fmt.Println("-----------encrptedOutString---------")
	fmt.Printf("%s\n", outString)

	// envelopeMail.SetEnryption()
	//
	//
	//
	// envelopeMail.WriteTo(os.Stdout)
}
