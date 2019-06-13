### apimanager
统一 Api 接口管理，调用入口,前端所有调用的中心

* 通过 mysql.sql 初始化 mysql,对应修改config.yaml
* 添加环境变量 
  * BaasFabricEngine=http://localhost:4991
* go run main.go
* 管理员帐号密码 
  * admin 123456