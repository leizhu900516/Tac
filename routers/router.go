package routers

import (
	"Tac/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/login", &controllers.LoginController{})
    beego.Router("/quit", &controllers.QuitController{})
    beego.Router("/auth", &controllers.AuthController{})
    beego.Router("/api/v1/backendtask", &controllers.BackgroundtaskController{})//get background list
    beego.Router("/api/v1/backendtaskmanage", &controllers.BackgroundtaskManageGetController{})  //filter taskid for detail
    beego.Router("/api/v1/addbackendtaskmanage", &controllers.BackgroundtaskManagePostController{})//add background task
    beego.Router("/api/v1/delbackendtaskmanage", &controllers.DelgroundtaskManageGetController{})//add background task
}
