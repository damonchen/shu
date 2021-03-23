package client

type AuthClient struct {
}

func (c AuthClient) Login(userLogin UserLogin) (*UserLoginResp, error) {
	resp, err := bundle.Get("auth").Call("login", userLogin)
	if err != nil {
		return nil, err
	}
	return resp.(*UserLoginResp), nil
}

func GetAuthClient() AuthClient {
	return AuthClient{}
}
