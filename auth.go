package invoice

// GetAuthorization 获取授权
func (c *Client) GetAuthorization(nsrsbh string, accountType ...string) (*AuthResponse, error) {
	params := map[string]string{
		"nsrsbh": nsrsbh,
	}

	// accountType 5 管理
	if len(accountType) > 0 && accountType[0] != "" {
		params["type"] = accountType[0]
	}

	resp, err := c.doRequest("POST", "/v5/enterprise/authorization", params)
	if err != nil {
		return nil, err
	}

	var authResp AuthResponse
	if err := parseResponseData(resp, &authResp); err != nil {
		return nil, err
	}

	// 自动设置token
	c.SetToken(authResp.Token)
	return &authResp, nil
}

// LoginDppt 登录数电发票平台
func (c *Client) LoginDppt(nsrsbh, username, password, sms string, options ...map[string]string) (*Response, error) {
	params := map[string]string{
		"nsrsbh":   nsrsbh,
		"username": username,
		"password": password,
	}

	if sms != "" {
		params["sms"] = sms
	}

	// 添加可选参数
	if len(options) > 0 {
		for k, v := range options[0] {
			params[k] = v
		}
	}

	return c.doRequest("POST", "/v5/enterprise/loginDppt", params)
}

// GetFaceImg 获取人脸二维码
func (c *Client) GetFaceImg(nsrsbh string, options ...map[string]string) (*FaceQRCodeResponse, error) {
	params := map[string]string{
		"nsrsbh": nsrsbh,
	}

	// 添加可选参数
	if len(options) > 0 {
		for k, v := range options[0] {
			params[k] = v
		}
	}

	resp, err := c.doRequest("GET", "/v5/enterprise/getFaceImg", params)
	if err != nil {
		return nil, err
	}

	var faceResp FaceQRCodeResponse
	if err := parseResponseData(resp, &faceResp); err != nil {
		return nil, err
	}

	return &faceResp, nil
}

// GetFaceState 获取人脸二维码认证状态
func (c *Client) GetFaceState(nsrsbh, rzid string, options ...map[string]string) (*FaceStateResponse, error) {
	params := map[string]string{
		"nsrsbh": nsrsbh,
		"rzid":   rzid,
	}

	// 添加可选参数
	if len(options) > 0 {
		for k, v := range options[0] {
			params[k] = v
		}
	}

	resp, err := c.doRequest("GET", "/v5/enterprise/getFaceState", params)
	if err != nil {
		return nil, err
	}

	var stateResp FaceStateResponse
	if err := parseResponseData(resp, &stateResp); err != nil {
		return nil, err
	}

	return &stateResp, nil
}

// QueryFaceAuthState 获取认证状态
func (c *Client) QueryFaceAuthState(nsrsbh string, options ...map[string]string) (*Response, error) {
	params := map[string]string{
		"nsrsbh": nsrsbh,
	}

	// 添加可选参数
	if len(options) > 0 {
		for k, v := range options[0] {
			params[k] = v
		}
	}

	return c.doRequest("POST", "/v5/enterprise/queryFaceAuthState", params)
}
