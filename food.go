package main

import (
	"fmt"
	"log"
	"net/http"
)

type Food struct {
	Id        uint64 `json:"id"`
	Order_num int32  `json:"order_num"`
	Picture   string `json:"picture"`
	Price     uint32 `json:"price"`
	Name      string `json:"name"`
	Remark    string `json:"remark"`
}

type GetFoodListRequest struct {
	Num      uint32
	Passback string
}

type GetFoodListResponse struct {
	Has_more  bool   `json:"has_more"`
	Pass_back string `json:"pass_back"`
	Food_list []Food `json:"food_list"`
}

func getFoodListHandler(w http.ResponseWriter, r *http.Request) {
	res := &JsonResponse{}
	defer log.Printf("%v %v %v %+v\n", r.RemoteAddr, r.Method, r.URL, res)
	ctx := NewContext(w, r)
	err := r.ParseForm()
	if err != nil {
		log.Printf("parseForm err:%v\n", err)
		res.Code = 403
		res.Message = "invalid request"
		ctx.SetJsonResponse(res)
		return
	}
	query := fmt.Sprintf("select id,order_num,picture,price,name,remark from %v.%v ", conf.Database.Database, conf.Database.FoodTable)
	query += fmt.Sprintf(" order by order_num limit %v", 10)
	err = db.Ping()
	if err != nil {
		log.Printf("db err:%v\n", err)
		db, err = NewDB(&conf.Database.DBConf)
		if err != nil {
			res.Code = 500
			res.Message = "lost connection to database"
			ctx.SetJsonResponse(res)
			return
		}
	}
	food := &Food{}
	err = db.QueryRow(query).Scan(&food.Id, &food.Order_num, &food.Picture, &food.Price, &food.Name, &food.Remark)
	if err != nil {
		log.Printf("db err:%v\n", err)
		res.Code = 500
		res.Message = "query failed"
		ctx.SetJsonResponse(res)
		return
	}
	data := GetFoodListResponse{
		Has_more: false,
		Food_list: []Food{
			*food,
		},
	}
	res.Data = data
	ctx.SetJsonResponse(res)
}
