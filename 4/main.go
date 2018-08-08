package main

import (
	"StockFinance/4/config"
	"StockFinance/4/entity"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func main() {
	config.Init()
	http.HandleFunc("/", index)
	http.HandleFunc("/stocks", getAllStocks)
	http.HandleFunc("/stocks/form", createForm)
	http.HandleFunc("/stocks/create", createStocks)
	http.ListenAndServe(":8000", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/stocks", http.StatusSeeOther)
}

func getAllStocks(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := config.Db.Query("SELECT * FROM stocks;")
	if err != nil {
		http.Error(w, "IULALALA", 500)
	}
	defer rows.Close()

	stocks := make([]entity.Stock, 0)
	for rows.Next() {
		stk := entity.Stock{}
		err := rows.Scan(&stk.Id, &stk.Date, &stk.Open, &stk.High, &stk.Low, &stk.Close, &stk.Volume, &stk.Action)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		stocks = append(stocks, stk)
		for _, value := range stocks {
			value.Date.Format("2006-01-02")
		}
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	config.Tpl.ExecuteTemplate(w, "stocks.gohtml", stocks)

}

func createForm(w http.ResponseWriter, r *http.Request) {
	config.Tpl.ExecuteTemplate(w, "form.gohtml", nil)
}

func createStocks(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method != "POST" {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	stk := entity.Stock{}
	stk.Date, err = time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		panic(err)
	}
	stk.Open, err = strconv.Atoi(r.FormValue("open"))
	if err != nil {
		panic(err)
	}
	stk.High, err = strconv.Atoi(r.FormValue("high"))
	if err != nil {
		panic(err)
	}
	stk.Low, err = strconv.Atoi(r.FormValue("low"))
	if err != nil {
		panic(err)
	}
	stk.Close, err = strconv.Atoi(r.FormValue("close"))
	if err != nil {
		panic(err)
	}
	stk.Volume, err = strconv.Atoi(r.FormValue("volume"))
	if err != nil {
		panic(err)
	}
	stk.Action = "buy"

	_, err = config.Db.Exec("INSERT INTO stocks (date, open_price, high_price, low_price, close_price, volume,action) VALUES ($1, $2, $3, $4, $5, $6, $7)", stk.Date, stk.Open, stk.High, stk.Low, stk.Close, stk.Volume, stk.Action)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.Tpl.ExecuteTemplate(w, "confirmation.gohtml", stk)

}

func calculateProfit() {
	priceList := []int{100, 120, 100, 130}
	priceMin := priceList[0]
	priceMax := priceList[0]
	different := 0
	// buyValue := 0
	// sellValue := 0
	profit := 0

	for i := 0; i < len(priceList)-1; i++ {
		if priceMax < priceList[i+1] {
			priceMax = priceList[i+1]
		} else {
			different = priceMax - priceMin
			if different > profit {
				profit = different
				//profit, buyValue, sellValue = different, priceMin, priceMax
			}
			priceMin, priceMax = priceList[i+1], priceList[i+1]
		}
		//fmt.Println(priceMin,priceMax,buyValue, sellValue, profit)

	}
	different = priceMax - priceMin
	if different > profit {
		profit = different
		//profit, buyValue, sellValue = different, priceMin, priceMax
	}

	fmt.Println(priceMin, priceMax, profit)
}
