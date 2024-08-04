package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync"
)

var sourceurls []string

func createUrls(){
	for i:= 1; i <= 10; i++ {
		sourceurls = append(sourceurls, "https://reqres.in/api/users/"+strconv.Itoa(i))
	}
}

type Response struct {
    Data struct {
        ID        int    `json:"id"`
        Email     string `json:"email"`
        FirstName string `json:"first_name"`
        LastName  string `json:"last_name"`
        Avatar    string `json:"avatar"`
    } `json:"data"`
    Support struct {
       URL  string `json:"url"`
       Text string `json:"text"`
    } `json:"support"`
}

var users = []Response{}

func fetchJSON(wg *sync.WaitGroup, url string){
	response,err := http.Get(url)
	if(err != nil){
		log.Fatal(err)	
	}
	resStruct := Response{}
	value,reserr := io.ReadAll(response.Body)
	err = json.Unmarshal(value,&resStruct)
	if(err != nil || reserr != nil){
		return
	}
	users = append(users, resStruct)	
	response.Body.Close()
	wg.Done()
}

func structToMap(res interface{}) map[string]interface{}{
	result := make(map[string]interface{})
	val := reflect.ValueOf(res)

	if val.Kind() == reflect.Ptr{
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
        fieldName := typ.Field(i).Name
        fieldValueKind := val.Field(i).Kind()
        var fieldValue interface{}

        if fieldValueKind == reflect.Struct {
            fieldValue = structToMap(val.Field(i).Interface())
        } else {
            fieldValue = val.Field(i).Interface()
        }

        result[fieldName] = fieldValue
    }
	result["custom_editedJSON"] = true

    return result

}


func main(){
	createUrls()
	var wg sync.WaitGroup
	wg.Add(len(sourceurls))
	file,fileerr := os.OpenFile("user.txt",os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if(fileerr != nil){
		return
	}
	for _,url := range(sourceurls){
		go fetchJSON(&wg,url)
	}
	defer file.Close()
	wg.Wait()
	var mapusers = []map[string]interface{}{}
	for _,userval := range(users){
		mapusers = append(mapusers, structToMap(userval))
	}
	jsonvalue,err :=  json.Marshal(mapusers)
	if(err != nil){
		log.Fatal(err)
	}else{
		file.Write(jsonvalue)
	}
}
