package db

import (
	"fmt"

	"gopkg.in/olivere/elastic.v5"
)

type EsDb struct {
	client *elastic.Client
}

func NewEsDb(para *DbConnPara) *EsDb {
	//addr := fmt.Sprintf("http://%s:%d", para.Host, para.Port)
	client, err := elastic.NewClient(elastic.SetURL("http://192.168.56.133:9200"))
	if err != nil {
		fmt.Println(err, para.Host, para.Port)
		return nil
	}
	return &EsDb{client: client}
}
