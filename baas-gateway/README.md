# baas-gateway
统一Api网关管理，调用入口,前端所有调用的中心  

步骤
* 通过 mysql.sql 初始化 mysql,对应修改dbconfig.yaml
                             
* 修改配置文件 gwconfig.yaml  
  配置文件的加载目录/etc/baas和当前目录  
  
运行
* go run main.go
* 管理员帐号密码 
  * admin 123456