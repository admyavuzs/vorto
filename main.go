package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "github.com/lib/pq"
	"time"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
			for {
					opsProcessed.Inc()
					time.Sleep(2 * time.Second)
			}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
			Name: "myapp_processed_ops_total",
			Help: "The total number of processed events",
	})
)


type Delivery struct {
	id   int `json:"id"`
}

const (
	host     = "postgre.postgresql.svc.cluster.local"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "coffee"
)

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	rows, err := db.Query("select id from delivery where delivery.id not in (select de.id from delivery de, carrier_bean_type cbt  ,driver  dr ,supplier_bean_type  sbt where  de.supplier_id=sbt.supplier_id and de.driver_id=dr.id and dr.carrier_id = cbt.carrier_id and cbt.bean_type_id =sbt.bean_type_id)")
	//rows, err := db.Query("select delivery.id from delivery;")
	if err != nil {
		log.Fatal(err)
	}

	var reqdel []Delivery

	for rows.Next() {
		var delivery Delivery
		rows.Scan(&delivery.id)
		reqdel = append(reqdel, delivery)
	}

	reqdelBytes, _ := json.MarshalIndent(reqdel, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(reqdelBytes)

	defer rows.Close()
	defer db.Close()
}

func main() {

	recordMetrics()

    http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", GETHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
} 