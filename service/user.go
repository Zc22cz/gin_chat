package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetIndex
// @Summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code", "msg"}
// @Router /user/getUserList [post]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "查找成功！",
		"data": data,
	})
}

// GetIndex
// @Summary 用户注册
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code", "msg"}
// @Router /user/creatUser [post]
func CreatUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("Identity")

	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.UserBasic{}
	utils.DB.Where("name = ?", user.Name).First(&data)
	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "用户名或密码不能为空！",
			"data": user,
		})
		return
	}
	if data.Name != "" {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "用户名已注册！",
			"data": user,
		})
		return
	}
	if password != repassword {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "两次密码不一致！",
			"data": user,
		})
		return
	}
	user.Password = utils.MakePassword(password, salt)
	user.Salt = salt
	user.LoginTime = time.Now()
	user.HeartBeatTime = time.Now()
	user.LoginOutTime = time.Now()
	utils.DB.Create(&user)
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "用户注册成功！",
		"data": user,
	})
}

// GetIndex
// @Summary 用户登录
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code", "msg"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	user := models.UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)

	//token加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	tmp := utils.MD5Encode(str)
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identity", tmp)

	if user.Name == "" {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}

	if !utils.ValidPassword(password, user.Salt, user.Password) {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "密码错误",
			"data": nil,
		})
		return
	}

	pwd := utils.MakePassword(password, user.Salt)
	data := models.UserBasic{}
	utils.DB.Where("name = ? and password = ?", name, pwd).First(&data)

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "登录成功",
		"data": data,
	})
}

// GetIndex
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code", "msg"}
// @Router /user/deleteUser [post]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	utils.DB.Delete(&user)
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "删除用户成功！",
		"data": user,
	})
}

// GetIndex
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code", "msg"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "输入参数格式错误",
			"data": user,
		})
		return
	}
	utils.DB.Model(&user).Updates(models.UserBasic{
		Name:     user.Name,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
	})
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "修改用户成功！",
		"data": user,
	})
}

// websocket 升级并跨域
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(c, ws)
}

func RedisMsg(c *gin.Context) {
	userIdA, _ := strconv.Atoi(c.PostForm("userIdA"))
	userIdB, _ := strconv.Atoi(c.PostForm("userIdB"))
	start, _ := strconv.Atoi(c.PostForm("start"))
	end, _ := strconv.Atoi(c.PostForm("end"))
	isRev, _ := strconv.ParseBool(c.PostForm("isRev"))
	res := models.RedisMsg(int64(userIdA), int64(userIdB), int64(start), int64(end), isRev)
	utils.RespOkList(c.Writer, "ok", res)
}

func MsgHandler(c *gin.Context, ws *websocket.Conn) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(" MsgHandler 发送失败", err)
		}

		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

func SearchFriends(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("userId"))
	user := models.SearchFriend(uint(id))
	//c.JSON(200, gin.H{
	//	"code": 0,
	//	"msg":  "修改用户成功！",
	//	"data": user,
	//})
	utils.RespOkList(c.Writer, user, len(user))
}

func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	targetName := c.Request.FormValue("targetName")
	code, msg := models.AddFriend(uint(userId), targetName)
	if code == 0 {
		utils.RespOk(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// 新建群
func CreateCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	name := c.Request.FormValue("name")
	icon := c.Request.FormValue("icon")
	desc := c.Request.FormValue("desc")
	community := models.Community{
		OwnerId: uint(ownerId),
		Name:    name,
		Img:     icon,
		Desc:    desc,
	}
	code, msg := models.CreateCommunity(community)
	if code == 0 {
		utils.RespOk(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// 加载群列表
func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	data, msg := models.LoadCommunity(uint(ownerId))
	if len(data) != 0 {
		utils.RespList(c.Writer, 0, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// 加入群
func JoinGroups(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	comId := c.Request.FormValue("comId")
	data, msg := models.JoinGroup(uint(userId), comId)
	if data == 0 {
		utils.RespOk(c.Writer, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

func FindByID(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	data := models.UserBasic{}
	utils.DB.Where("id = ?", userId).First(&data)
	utils.RespOk(c.Writer, data, "OK")
}
