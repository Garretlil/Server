package main

//func IndexHandler(w http.ResponseWriter, r *http.Request) {
//
//	rows, err := database.Query("select * from products ")
//	if err != nil {
//		log.Println(err)
//	}
//	defer rows.Close()
//	products := []Product{}
//
//	for rows.Next() {
//		p := Product{}
//		err := rows.Scan(&p.Id, &p.Model, &p.Company, &p.Price)
//		if err != nil {
//			fmt.Println(err)
//			continue
//		}
//		products = append(products, p)
//	}
//
//	tmpl, _ := template.ParseFiles("templates/index.html")
//	tmpl.Execute(w, products)
//}
