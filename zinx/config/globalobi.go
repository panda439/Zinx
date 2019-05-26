package config

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type GlobalObj struct {


	Host string //当前监听的IP
	Port int
	Name string//当前zinxserver名称

	Version string //当前框架的版本号
	MaxPackageSize uint32//每次Read一次的最大长度

}

//定义一个全局的对外的配置的对象
var GlobalObject *GlobalObj

func (g GlobalObj)LoadConfig()  {
	data,err := ioutil.ReadFile("./conf/zinx.json")
	if err != nil {
		panic(err)
	}

	//将zinx.json 的数据
	err = json.Unmarshal(data,&GlobalObject)
	if err != nil {
		panic(err)
	}
	//查看配置文件是否加载

	if GlobalObject.MaxPackageSize == 4096 {
		fmt.Println("配置文件加载成功")
	}else {
		fmt.Println("配置文件未加载成功")

	}

}

//只要import 当前模块 就会执行init方法 加载配置文件
func init() {
	//配置文件读取操作
	GlobalObject = &GlobalObj{
		//默认值
		Name : "Zinx",
		Host :  "0.0.0.0",
		Port : 8999,
		Version: "V1.0",
		MaxPackageSize:512,
	}

	//GlobalObject.LoadConfig()

}
