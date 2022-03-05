package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	pb "go-scaffold/internal/app/api/v1/user"
	"go-scaffold/internal/app/pkg/responsex"
	"go-scaffold/internal/app/transport/http/pkg/bindx"
)

type DeleteReq struct {
	Id uint `uri:"id" binding:"required"` // 用户 ID
}

func (DeleteReq) ErrorMessage() map[string]string {
	return map[string]string{
		"ID.required": "用户 ID 不能为空",
	}
}

// Delete 删除用户
// @Router       /v1/user/{id} [delete]
// @Summary      删除用户
// @Description  删除用户
// @Tags         用户
// @Accept       plain
// @Produce      json
// @Param        id   path      integer                   true  "用户 ID"  format(uint)  minimum(1)
// @Success      200  {object}  example.Success           "成功响应"
// @Failure      500  {object}  example.ServerError       "服务器出错"
// @Failure      400  {object}  example.ClientError       "客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
// @Failure      401  {object}  example.Unauthorized      "登陆失效"
// @Failure      403  {object}  example.PermissionDenied  "没有权限"
// @Failure      404  {object}  example.ResourceNotFound  "资源不存在"
// @Failure      429  {object}  example.TooManyRequest    "请求过于频繁"
func (h *handler) Delete(ctx *gin.Context) {
	req := new(DeleteReq)
	if err := bindx.ShouldBindUri(ctx, req); err != nil {
		h.logger.Error(err.Error())
		return
	}

	param := new(pb.DeleteRequest)
	if err := copier.Copy(param, req); err != nil {
		h.logger.Error(err.Error())
		responsex.ServerError(ctx)
		return
	}

	_, err := h.service.Delete(ctx.Request.Context(), param)
	if err != nil {
		responsex.ServerError(ctx, responsex.WithMsg(err.Error()))
		return
	}

	responsex.Success(ctx)
	return
}
