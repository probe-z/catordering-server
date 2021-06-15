package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/schema"
)

type Food struct {
	Id       uint64 `json:"id"`
	OrderNum int32  `json:"orderNum"`
	Picture  string `json:"picture"`
	Price    uint32 `json:"price"`
	Name     string `json:"name"`
	Remark   string `json:"remark"`
}

type GetFoodListRequest struct {
	Cursor string `schema:"cursor"`
	Num    uint32 `schema:"num"`
}

type GetFoodListResponse struct {
	Cursor   string `json:"cursor"`
	HasMore  bool   `json:"hasMore"`
	FoodList []Food `json:"foodList"`
}

func getFoodListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r)
	err := r.ParseForm()
	if err != nil {
		ctx.SetJsonResponse(&JsonResponse{Code: 403, Message: "invalid request"})
		return
	}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	req := &GetFoodListRequest{}
	err = decoder.Decode(req, ctx.Request.Form)
	if err != nil {
		ctx.SetJsonResponse(&JsonResponse{Code: 403, Message: "invalid request"})
		return
	}
	cursor := uint64(0)
	if req.Cursor != "" {
		cursor, err = strconv.ParseUint(req.Cursor, 10, 64)
		if err != nil {
			ctx.SetJsonResponse(&JsonResponse{Code: 403, Message: "invalid request"})
			return
		}
	}
	if req.Num == 0 || req.Num > 30 {
		ctx.SetJsonResponse(&JsonResponse{Code: 403, Message: "invalid request"})
		return
	}
	query := fmt.Sprintf("select id,order_num,picture,price,name,remark from %v.%v", conf.Database.Database, conf.Database.FoodTable)
	if cursor > 0 {
		query += fmt.Sprintf(" where order_num>%v", cursor)
	}
	query += fmt.Sprintf(" order by order_num limit %v", req.Num+1)
	log.Printf("sql:%v\n", query)
	err = db.Ping()
	if err != nil {
		log.Printf("db err:%v\n", err)
		ctx.SetJsonResponse(&JsonResponse{Code: 500, Message: "database error"})
		return
	}
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("db err:%v\n", err)
		ctx.SetJsonResponse(&JsonResponse{Code: 500, Message: "database error"})
		return
	}
	if rows == nil {
		ctx.SetJsonResponse(&JsonResponse{Code: 0, Message: "ok"})
		return
	}
	data := GetFoodListResponse{}
	for rows.Next() {
		food := &Food{}
		err = rows.Scan(&food.Id, &food.OrderNum, &food.Picture, &food.Price, &food.Name, &food.Remark)
		if err != nil {
			log.Printf("db err:%v\n", err)
			ctx.SetJsonResponse(&JsonResponse{Code: 500, Message: "database error"})
			return
		}
		data.FoodList = append(data.FoodList, *food)
		if len(data.FoodList) >= int(req.Num) {
			break
		}
	}
	data.HasMore = rows.Next()
	if len(data.FoodList) > 0 {
		data.Cursor = fmt.Sprintf("%v", data.FoodList[len(data.FoodList)-1].OrderNum)
	}
	ctx.SetJsonResponse(&JsonResponse{Code: 0, Message: "ok", Data: data})
}
