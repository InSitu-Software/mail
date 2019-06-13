package mail

import (
	"bytes"

	"gitlab.insitu.de/golang/pgp"
)

type PGPMessageEncryptor struct {
	PlainMessage bytes.Buffer
	KeyProvider  pgp.KeyProvider
	To           []string
}

func (mw *PGPMessageEncryptor) Write(b []byte) (n int, err error) {
	mw.PlainMessage.Grow(len(b))
	n, err = mw.PlainMessage.Write(b)
	return n, err
}

func (mw *PGPMessageEncryptor) GetEncryptedString() (string, error) {

	encWriterTo, err := pgp.Encrypt(&mw.PlainMessage, mw.To, mw.KeyProvider)
	if err != nil {
		return "", err
	}

	var output bytes.Buffer
	if _, err := encWriterTo.WriteTo(&output); err != nil {
		return "", err
	}

	return output.String(), nil
}
