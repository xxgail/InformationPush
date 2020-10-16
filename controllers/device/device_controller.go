package device

import (
	"InformationPush/common"
	"InformationPush/controllers"
	"InformationPush/helpers"
	"InformationPush/lib/mysqllib"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"time"
)

type InitParam struct {
	DeviceToken string `form:"device_token" json:"device_token" binding:"required"`
	GroupId     string `form:"group_id" json:"group_id" binding:"required"`
	Uid         string `form:"uid" json:"uid"`
}

type Device struct {
	id          string
	channel     string
	deviceToken string
	appId       string
	groupId     string
}

/**
 * @Date 2020/10/16
 * @Desc 开启APP时注册设备
 * @Param channel、device_token、app_id、group_id
 * @return
 **/
func Register(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")
	fmt.Println("Content-Type:", contentType)

	var initRegisterParam InitParam
	var err error

	switch contentType {
	case "application/json":
		err = c.ShouldBindJSON(&initRegisterParam)
	case "application/x-www-form-urlencoded":
		err = c.ShouldBindWith(&initRegisterParam, binding.Form)
	}
	if err != nil {
		fmt.Println(err)
	}

	channel := c.Request.Header.Get("Channel")
	deviceToken := initRegisterParam.DeviceToken
	appId := c.Request.Header.Get("appid")
	groupId := initRegisterParam.GroupId

	var device Device
	mysqlClient := mysqllib.GetMysqlConn()
	query := "SELECT id FROM device WHERE channel = '" + channel + "' AND device_token = '" + deviceToken + "' AND app_id = '" + appId + "' AND group_id = '" + groupId + "'"
	fmt.Println(query)
	err = mysqlClient.QueryRow(query).Scan(&device.id)
	if err != nil {
		fmt.Println("Register 查询数据库单条用户信息 出错：", err)
	}

	fmt.Println("deviceInfo", device)

	ip := helpers.GetServerIp()
	currentTime := time.Now()
	//开启事务
	tx, err := mysqlClient.Begin()
	if err != nil {
		fmt.Println("tx fail")
	}
	if device.id == "" {
		queryBid := "INSERT INTO device (`channel`,`device_token`,`app_id`,`group_id`,`ip`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?)"
		stemBid, err := tx.Prepare(queryBid)
		if err != nil {
			fmt.Println("Insert-device----Prepare fail")
		}
		_, err = stemBid.Exec(channel, deviceToken, appId, groupId, ip, currentTime, currentTime)
		if err != nil {
			fmt.Println("Insert-device----Exec fail")
		}
	} else {
		//准备sql语句
		stmt, err := tx.Prepare("UPDATE device SET ip = ?, updated_at = ? WHERE id = ?")
		if err != nil {
			fmt.Println("Prepare fail")
		}
		//设置参数以及执行sql语句
		res, err := stmt.Exec(ip, currentTime, device.id)
		if err != nil {
			fmt.Println("Exec fail")
		}
		//提交事务
		tx.Commit()
		val, _ := res.RowsAffected()
		fmt.Println(val, "----- 用户ID：", "更新完毕")
	}
	// 定义接口返回data
	data := make(map[string]interface{})
	controllers.Response(c, common.HTTPOK, "注册成功！", data)
}

/**
 * @Date 2020/10/16
 * @Desc 登录
 * @Param channel、device_token、app_id、group_id、uid
 * @return
 **/
func Login(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")
	fmt.Println("Content-Type:", contentType)

	var initRegisterParam InitParam
	var err error

	switch contentType {
	case "application/json":
		err = c.ShouldBindJSON(&initRegisterParam)
	case "application/x-www-form-urlencoded":
		err = c.ShouldBindWith(&initRegisterParam, binding.Form)
	}
	if err != nil {
		fmt.Println(err)
	}

	channel := c.Request.Header.Get("Channel")
	deviceToken := initRegisterParam.DeviceToken
	appId := c.Request.Header.Get("appid")
	groupId := initRegisterParam.GroupId
	uid := initRegisterParam.Uid
	currentTime := time.Now()

	mysqlClient := mysqllib.GetMysqlConn()
	//开启事务
	tx, err := mysqlClient.Begin()
	if err != nil {
		fmt.Println("tx fail")
	}
	// ① 下线其他uid登录的设备
	// 准备sql语句
	stmt, err := tx.Prepare("UPDATE device SET uid = null, updated_at = ? WHERE uid = ?")
	if err != nil {
		fmt.Println("Prepare fail")
	}
	//设置参数以及执行sql语句
	res, err := stmt.Exec(currentTime, uid)
	if err != nil {
		fmt.Println("Exec fail")
	}
	// ② 更新本条UID
	// 准备sql语句
	stmt1, err := tx.Prepare("UPDATE device SET uid = ?, updated_at = ? WHERE channel = ? AND device_token = ? AND app_id = ? AND group_id = ?")
	if err != nil {
		fmt.Println("Prepare fail")
	}
	//设置参数以及执行sql语句
	res1, err := stmt1.Exec(uid, currentTime, channel, deviceToken, appId, groupId)
	if err != nil {
		fmt.Println("Exec fail")
	}
	//提交事务
	tx.Commit()

	val, _ := res.RowsAffected()
	fmt.Println(val, "----- 用户清除UID：", "更新完毕")

	val1, _ := res1.RowsAffected()
	fmt.Println(val1, "----- 用户更新UID：", "更新完毕")

	// 定义接口返回data
	data := make(map[string]interface{})
	controllers.Response(c, common.HTTPOK, "登录成功！", data)
}

/**
 * @Date 2020/10/16
 * @Desc 退出登录
 * @Param uid
 * @return
 **/
func Logout(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")
	fmt.Println("Content-Type:", contentType)

	var initRegisterParam InitParam
	var err error

	switch contentType {
	case "application/json":
		err = c.ShouldBindJSON(&initRegisterParam)
	case "application/x-www-form-urlencoded":
		err = c.ShouldBindWith(&initRegisterParam, binding.Form)
	}
	if err != nil {
		fmt.Println(err)
	}

	uid := initRegisterParam.Uid
	currentTime := time.Now()

	mysqlClient := mysqllib.GetMysqlConn()
	//开启事务
	tx, err := mysqlClient.Begin()
	if err != nil {
		fmt.Println("tx fail")
	}
	// ① 下线其他uid登录的设备
	// 准备sql语句
	stmt, err := tx.Prepare("UPDATE device SET uid = null, updated_at = ? WHERE uid = ?")
	if err != nil {
		fmt.Println("Prepare fail")
	}
	//设置参数以及执行sql语句
	res, err := stmt.Exec(currentTime, uid)
	if err != nil {
		fmt.Println("Exec fail")
	}
	//提交事务
	tx.Commit()

	val, _ := res.RowsAffected()
	fmt.Println(val, "----- 用户清除UID：", "更新完毕")

	// 定义接口返回data
	data := make(map[string]interface{})
	controllers.Response(c, common.HTTPOK, "退出登录！", data)
}
