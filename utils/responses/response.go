package responses

type MapResponse struct {
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func WebResponse(message string, data interface{}) MapResponse {
    return MapResponse{
        Message: message,
        Data:    data,
    }
}

type MapResponsePagi struct {
    Message   string      `json:"message"`
    TotalPage int         `json:"total_page"`
    Data      interface{} `json:"data,omitempty"`
}

func WebResponsePagi(message string, data interface{}, totalPage int) MapResponsePagi {
    return MapResponsePagi{
        Message:   message,
        TotalPage: totalPage,
        Data:      data,
    }
}