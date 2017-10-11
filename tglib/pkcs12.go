package tglib

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"os/exec"
)

// GeneratePKCS12FromFiles generates a full PKCS certificate based on the input keys.
func GeneratePKCS12FromFiles(out, certPath, keyPath, caPath, passphrase string) error {

	args := []string{
		"pkcs12",
		"-export",
		"-out", out,
		"-inkey", keyPath,
		"-in", certPath,
		"-certfile", caPath,
		"-passout", "pass:" + passphrase,
	}

	return exec.Command("openssl", args...).Run()
}

// GeneratePKCS12 generates a pkcs12
func GeneratePKCS12(cert []byte, key []byte, ca []byte, passphrase string) ([]byte, error) {
	// cert
	tmpcert, err := ioutil.TempFile("", "tmpcert")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpcert.Name()) // nolint: errcheck
	defer tmpcert.Close()           // nolint: errcheck
	if _, err = tmpcert.Write(cert); err != nil {
		return nil, err
	}

	// key
	tmpkey, err := ioutil.TempFile("", "tmpkey")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpkey.Name()) // nolint: errcheck
	defer tmpkey.Close()           // nolint: errcheck
	if _, err = tmpkey.Write(key); err != nil {
		return nil, err
	}

	// ca
	tmpca, err := ioutil.TempFile("", "tmpca")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpca.Name()) // nolint: errcheck
	defer tmpca.Close()           // nolint: errcheck
	if _, err = tmpca.Write(ca); err != nil {
		return nil, err
	}

	// p12
	tmpp12, err := ioutil.TempFile("", "tmpp12")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpp12.Name()) // nolint: errcheck
	defer tmpp12.Close()           // nolint: errcheck

	if err = GeneratePKCS12FromFiles(tmpp12.Name(), tmpcert.Name(), tmpkey.Name(), tmpca.Name(), passphrase); err != nil {
		return nil, err
	}

	p12data, err := ioutil.ReadAll(tmpp12)
	if err != nil {
		return nil, err
	}
	return p12data, nil
}

// GenerateBase64PKCS12 generates a full PKCS certificate based on the input keys.
func GenerateBase64PKCS12(cert []byte, key []byte, ca []byte, passphrase string) (string, error) {

	p12data, err := GeneratePKCS12(cert, key, ca, passphrase)

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(p12data), nil
}
