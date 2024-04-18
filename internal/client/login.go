package client

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/shulganew/GophKeeperClient/internal/entities"
	"go.uber.org/zap"
)

func UserReg() {
	i := strconv.Itoa(rand.Intn(1000))
	user := entities.User{Login: "Igor" + i, Password: "MySecret"}

	reqBodyDel := bytes.NewBuffer([]byte{})

	err := json.NewEncoder(reqBodyDel).Encode(&user)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8088/api/user/register", reqBodyDel)
	if err != nil {
		fmt.Println(err)
	}

	//reqest
	request.Header.Add("Content-Type", "application/json")
	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	//=============================Response===================
	for k, v := range res.Header {

		fmt.Printf("%s: %v\r\n", k, v[0])

	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Body: ", string(body))
	fmt.Printf("Status Code: %d\r\n", res.StatusCode)

	jwt := res.Header.Get("Authorization")[len("Bearer "):]
	user.JWT = sql.NullString{String: jwt, Valid: true}
	
	zap.S().Infoln(user.JWT, user.Login, user.Password)
}
