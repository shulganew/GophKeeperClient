package client

/*

func UserLogin(login, pw string) {

	var user model.User

	reqBodyDel := bytes.NewBuffer([]byte{})

	err = json.NewEncoder(reqBodyDel).Encode(&users[nUser])
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8088/api/user/login", reqBodyDel)
	if err != nil {
		fmt.Println(err)
	}

	// add jwt
	request.Header.Add("Authorization", users[nUser].JWT.String)

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

	//jwt := res.Header.Get("Authorization")[len("Bearer "):]
	jwt := res.Header.Get("Authorization")

	_, err = conn.Exec(ctx, "UPDATE users SET jwt = $1", jwt)
	if err != nil {
		panic(err)
	}

}
*/
