package main

import (
	"StockFinance/4/config"
	"StockFinance/4/entity"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"time"
)

var render *template.Template

func main() {
	config.Init()
	render = template.Must(template.ParseGlob("templates/*.html"))
	http.HandleFunc("/", index)
	//calculateProfit2()
	calculateProfit()
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

	rows, err := config.Db.Query("SELECT * FROM stocks ORDER BY date ASC;")
	if err != nil {
		http.Error(w, "General Error", 500)
	}
	defer rows.Close()

	response := entity.Response{}
	for rows.Next() {
		stk := entity.Stock{}
		err := rows.Scan(&stk.Id, &stk.Date, &stk.Open, &stk.High, &stk.Low, &stk.Close, &stk.Volume, &stk.Action)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		response.List = append(response.List, stk)
	}
	//fmt.Println(response.List)
	//fmt.Println(calculateProfit(response.List))
	response.Summary = calculateProfit()

	failed := render.ExecuteTemplate(w, "index.html", response)

	if failed != nil {
		panic(failed)
	}

}

func getProfit() []entity.Profit {
	rows, err := config.Db.Query("SELECT date,close_price FROM stocks")
	if err != nil {
		panic(err)
	}

	priceList := []entity.Profit{}

	for rows.Next() {
		profit := entity.Profit{}
		err := rows.Scan(&profit.Date, &profit.ClosePrice)
		if err != nil {
			panic(err)
		}

		priceList = append(priceList, profit)
	}
	sort.Slice(priceList, func(i, j int) bool {
		return priceList[i].ClosePrice < priceList[j].ClosePrice
	})
	return priceList
}

func createForm(w http.ResponseWriter, r *http.Request) {
	render.ExecuteTemplate(w, "form.html", nil)
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
	stk.Action = setAction(stk.Open, stk.Close)

	_, err = config.Db.Exec("INSERT INTO stocks (date, open_price, high_price, low_price, close_price, volume,action) VALUES ($1, $2, $3, $4, $5, $6, $7)", stk.Date, stk.Open, stk.High, stk.Low, stk.Close, stk.Volume, stk.Action)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	render.ExecuteTemplate(w, "confirmation.html", stk)

}

func setAction(openPrice, closePrice int) string {
	if openPrice < closePrice {
		return "sell"
	} else if openPrice > closePrice {
		return "buy"
	} else {
		return "hold"
	}

}

func calculateProfit() string {

	listStock := getProfit()

	result := "profit is :" + strconv.Itoa(listStock[len(listStock)-1].ClosePrice-listStock[0].ClosePrice) + "between" + listStock[len(listStock)-1].Date.String() + "and" + listStock[0].Date.String()
	return result
}
