package handler
import(
	"go-gcs/model"
	"go-gcs/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *service.UserService
}

// Create godoc
// @Summary 创建用户
// @Description 创建一个新用户，包含用户名、邮箱、密码和头像链接
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param user body model.User true "用户信息"
// @Success 200 {object} model.User "创建成功返回用户信息"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /users/create [post]
func (r *UserHandler) CreateUser(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	user, err := r.Service.CreateUser(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user created",
		"user": user,
	})
}

func (r *UserHandler) UpdateUser(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx :=  c.Request.Context()
	user, err := r.Service.UpdateUser(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user updated",
		"user": user,
	})
}
