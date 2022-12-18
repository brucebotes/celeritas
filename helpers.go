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
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/pusher/pusher-http-go"
)

const (
	randomString = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321_+"
)

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
func (c *Celeritas) BuildJSCSSscript(view, src string) error {
	if !c.Debug {
		return nil
	}
	result := api.Build(api.BuildOptions{
		NodePaths:   []string{view + "/node_modules/"},
		EntryPoints: []string{view + "/" + src},
		Bundle:      true,
		Loader: map[string]api.Loader{
			".css": api.LoaderCSS,
		},
		// disable minification for development
		//MinifyWhitespace:  true,
		//MinifyIdentifiers: true,
		//MinifySyntax:      true,
		Metafile:  true,
		Sourcemap: api.SourceMapLinked,
		Write:     true,
		Outdir:    "public/" + view,
	})
	if len(result.Errors) > 0 {
		log.Println("Esbuild errors \u2192", result.Errors)
		return errors.New("error compiling jsx and css")
	}
	ioutil.WriteFile("public/"+view+"/meta.json", []byte(result.Metafile), 0644)
	return nil
}

//BuildJScript take a javascript/typescript/react module and compile to esModule
//Write esModule to file in the static folder
func (c *Celeritas) BuildWithNpmScript(mod string) error {
	if !c.Debug {
		return nil
	}
	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = c.RootPath + "/views/" + mod

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println("Command error 1 compiling Javascript \u2192", err)
		return err
	}
	err = cmd.Start()
	if err != nil {
		log.Println("Command error 2 compiling Javascript \u2192", err)
		return err
	}

	data, err := ioutil.ReadAll(stderr)
	if err != nil {
		log.Println("Command error 3 compiling Javascript \u2192", err)
		return err
	}
	err = cmd.Wait()
	if err != nil {
		log.Println("Command error 4 compiling Javascript \u2192", err)
		return err
	}

	if data != nil && len(data) > 0 {
		c.ErrorLog.Println(fmt.Printf("Compilng '/views/%s': Error/Warning messages \u2192 \n%s\n", mod, string(data)))
	}

	return nil
}

func (c *Celeritas) AuthenticateWebsocket(userID int, userName string, params []byte) ([]byte, error) {

	presenceData := pusher.MemberData{
		UserID: strconv.Itoa(userID),
		UserInfo: map[string]string{
			"name": "FirstName",
			"id":   strconv.Itoa(userID),
		},
	}

	response, err := c.wsClient.AuthenticatePresenceChannel(params, presenceData)
	return response, err
}

func (c *Celeritas) ListenForWebsocketEvents(r *http.Request, buff []byte) (*pusher.Webhook, error) {
	webhook, err := c.wsClient.Webhook(r.Header, buff)
	if err != nil {
		return nil, err
	}
	return webhook, err
}

func (c *Celeritas) BroadcastWebsocketMessage(channel, messageType string, data interface{}) error {
	err := c.wsClient.Trigger(channel, messageType, data)
	if err != nil {
		c.ErrorLog.Println("Websocket broadcast error: " + err.Error())
		return err
	}

	return nil
}
