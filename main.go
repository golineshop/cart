package main

import (
	"github.com/golineshop/cart/domain/repository"
	"github.com/golineshop/cart/domain/service"
	"github.com/golineshop/cart/handler"
	proto "github.com/golineshop/cart/proto"
	"github.com/golineshop/common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/common/log"
)

var QPS = 100

func main() {
	// 设置配置中心
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		log.Error(err)
	}
	// 设置注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 设置链路追踪
	t, io, err := common.NewTracer("go.micro.service.cart", "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 服务参数设置
	serv := micro.NewService(
		micro.Name("go.micro.service.cart"),
		micro.Version("latest"),
		//暴露的服务地址
		micro.Address("0.0.0.0:8087"),
		//注册中心
		micro.Registry(consulRegistry),
		//链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		//添加consul 作为注册中心
		micro.Registry(consulRegistry),
		//添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
	)

	// 获取mysql配置，路径中不带前缀
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	// 连接数据库
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Error(err)
	}
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {
			println(err)
		}
	}(db)
	// 禁止复数表
	db.SingularTable(true)

	// 只执行一次
	//rp := repository.NewCartRepository(db)
	//rp.InitTable()

	// 初始化服务
	serv.Init()
	cartService := service.NewCartService(repository.NewCartRepository(db))
	err = proto.RegisterCartHandler(serv.Server(), &handler.CartController{CartService: cartService})
	if err != nil {
		log.Error(err)
	}
	// 运行服务
	if err := serv.Run(); err != nil {
		log.Fatal(err)
	}
}
