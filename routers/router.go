package routers

import (
	"Tac/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/api/v1/backendtask", &controllers.BackgroundtaskController{})
    beego.Router("/api/v1/backendtaskmanage", &controllers.BackgroundtaskManageGetController{})
}
