package celeritas

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/evanw/esbuild/pkg/api"
)

const (
	randomString = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321_+"
)
func (c *Celeritas) SetAppDebug(debug bool){
  c.Debug = debug
  c.Render.SetDebug(debug)
}

//RandomString  generates a random string length n from values in the const randomString
func (c *Celeritas) RandomString(n int) string {
	s, r := make([]rune, n), []rune(randomString)

	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}
	return string(s)
}

func (c *Celeritas) CreateDirIfNotExist(path string) error {
	const mode = 0755
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, mode)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Celeritas) CreateFileIfNotExists(path string) error {
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}

		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	}

	return nil
}

type Encryption struct {
	Key []byte
}

func (e *Encryption) Encrypt(text string) (string, error) {
	plainText := []byte(text)

	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", nil
	}

	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plainText)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (e *Encryption) Decrypt(cryptoText string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("invalid block size")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

//BuildJScript take a javascript/typescript/react module and compile to esModule
//Write esModule to file in the static folder
func (c *Celeritas) BuildJSCSSscript(module, src string) error {
	if !c.BundleJS {
		return nil
	}
	buildPath := c.RootPath + "/views/" + module
	log.Println("build path:", buildPath+"/"+src)
	result := api.Build(api.BuildOptions{
		NodePaths:   []string{buildPath + "/node_modules/"},
		EntryPoints: []string{buildPath + "/" + src},
		Bundle:      true,
		Loader: map[string]api.Loader{
			".css": api.LoaderCSS,
		},
		// disable minification for development
		//MinifyWhitespace:  true,
		//MinifyIdentifiers: true,
		//MinifySyntax:      true,
		Metafile:    true,
		Sourcemap:   api.SourceMapLinked,
		Write:       true,
		Outfile:     c.RootPath + "/public/views/" + module + "/bundle.js",
		TreeShaking: api.TreeShakingTrue,
	})
	if len(result.Errors) > 0 {
		for i, m := range result.Errors {
			log.Printf("Esbuild error %d: %s:%d %s\n", i, m.Location.File, m.Location.Line, m.Text)
		}
		return errors.New("error generating js and/or css bundle(s)")
	}
	ioutil.WriteFile("public/views/"+module+"/meta.json", []byte(result.Metafile), 0644)
	return nil
}

//BuildJScript take a javascript/typescript/react module and compile to esModule
//Write esModule to file in the static folder
func (c *Celeritas) BuildWithNpmScript(module string) error {
	if !c.BundleJS {
		return nil
	}
	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = c.RootPath + "/views/" + module

	stderr, err := cmd.StderrPipe()
	if err != nil {
		c.ErrorLog.Println("Command error 1 compiling Javascript \u2192", err)
		return err
	}
	err = cmd.Start()
	if err != nil {
		c.ErrorLog.Println("Command error 2 compiling Javascript \u2192", err)
		return err
	}

	data, err := ioutil.ReadAll(stderr)
	if err != nil {
		c.ErrorLog.Println("Command error 3 compiling Javascript \u2192", err)
		return err
	}
	err = cmd.Wait()
	if err != nil {
		c.ErrorLog.Println("Command error 4 compiling Javascript \u2192", err)
		return err
	}

	if data != nil && len(data) > 0 {
		//log.Printf("%+v\n", data)
		errStr := fmt.Sprintf("Compiling '/views/%s': Error/Warning messages \u2192 \n%s\n", module, string(data))
		err = errors.New(errStr)
		c.ErrorLog.Println(err)
		return err
	}

	return nil
}


