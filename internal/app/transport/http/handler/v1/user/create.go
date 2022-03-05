package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	pb "go-scaffold/internal/app/api/v1/user"
	"go-scaffold/internal/app/pkg/responsex"
	"go-scaffold/internal/app/transport/http/pkg/bindx"
)

type CreateReq struct {
	Name  string `json:"name" binding:"required"`        // 名称
	Age   int8   `json:"age" binding:"min=1"`            // 年龄
	Phone string `json:"phone" binding:"required,phone"` // 电话
}

func (CreateReq) ErrorMessage() map[string]string {
	return map[string]string{
		"Name.required":  "名称不能为空",
		"Age.min":        "年龄必须大于 0",
		"Phone.required": "手机号码不能为空",
		"Phone.phone":    "手机号码格式错误",
	}
}

// Create 创建用户
// @Router       /v1/user [post]
// @Summary      创建用户
// @Description  创建用户
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        user_info  body      CreateReq                 true  "用户信息"  format(string)
// @Success      200        {object}  example.Success           "成功响应"
// @Failure      500        {object}  example.ServerError       "服务器出错"
// @Failure      400        {object}  example.ClientError       "客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
// @Failure      401        {object}  example.Unauthorized      "登陆失效"
// @Failure      403        {object}  example.PermissionDenied  "没有权限"
// @Failure      404        {object}  example.ResourceNotFound  "资源不存在"
// @Failure      429        {object}  example.TooManyRequest    "请求过于频繁"
func (h *handler) Create(ctx *gin.Context) {
	req := new(CreateReq)
	if err := bindx.ShouldBindJSON(ctx, req); err != nil {
		h.logger.Error(err.Error())
		return
	}

	param := new(pb.CreateRequest)
	if err := copier.Copy(param, req); err != nil {
		h.logger.Error(err.Error())
		responsex.ServerError(ctx)
		return
	}

	_, err := h.service.Create(ctx.Request.Context(), param)
	if err != nil {
		responsex.ServerError(ctx, responsex.WithMsg(err.Error()))
		return
	}

	responsex.Success(ctx)
	return
}
