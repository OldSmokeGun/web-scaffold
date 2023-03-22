package v1

import (
	"net/http"

	"go-scaffold/internal/app/controller"
	berr "go-scaffold/internal/app/pkg/errors"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type GreetHandler struct {
	controller *controller.GreetController
}

func NewGreetHandler(controller *controller.GreetController) *GreetHandler {
	return &GreetHandler{controller}
}

type HelloResponse struct {
	Msg string `json:"msg"`
}

// Hello 示例方法
//
//	@Router			/v1/greet [get]
//	@Summary		示例接口
//	@Description	示例接口
//	@Tags			示例
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			name	query		string								true	"名称"	format(string)	default(Tom)
//	@Success		200		{object}	example.Success{data=HelloResponse}	"成功响应"
//	@Failure		500		{object}	example.ServerError					"服务器出错"
//	@Failure		400		{object}	example.ClientError					"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized				"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied			"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound			"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest				"请求过于频繁"
//	@Security		Authorization
func (h *GreetHandler) Hello(ctx echo.Context) error {
	req := new(controller.HelloRequest)
	if err := ctx.Bind(req); err != nil {
		return errors.Wrap(err, berr.ErrBadRequest.Error())
	}

	ret, err := h.controller.Hello(ctx.Request().Context(), *req)
	if err != nil {
		return err
	}

	data := HelloResponse{Msg: ret.Msg}

	return ctx.JSON(http.StatusOK, data)
}
