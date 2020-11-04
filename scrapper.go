package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	title    string
	location string
	salary   string
	summary  string
}

//얼마나 많은 페이지가 있는지 알아내는 프로그램

//Scrape Indeed by a term
func scrape(term string) {
	var baseURL string = "https://kr.indeed.com/jobs?q=" + term + "&limit=50"
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages(baseURL)
	//fmt.Println(totalPages)

	for i := 0; i < totalPages; i++ {
		//extractedJobs := getPage(i)
		go getPage(i, baseURL, c)
		//jobs = append(jobs, extractedJobs...)
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Link", "Title", "Location", "Salary", "Summary"} //Id에서 링크로 바뀜

	wErr := w.Write(headers) //Write함수는 에러를 리턴함
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func getPage(page int, url string, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := url + "&start=" + strconv.Itoa(page*50) //=page*50 	string이 아니라 숫자니까 go에 포함된 패키지를 이용해 형변환해서 넣는다 생각 (int to asci/Itoa)
	fmt.Println("Requestring", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")

	searchCards.Each(func(i int, card *goquery.Selection) {
		// id, _ := card.Attr("data-jk")
		// //fmt.Println(id)
		// title := cleanString(card.Find(".title>a").Text())
		// //fmt.Println(title)
		// location := cleanString(card.Find(".sjcl").Text())
		// fmt.Println(id, title, location)
		go extractJob(card, c)
		//jobs = append(jobs, job) //append를 이용하여 jobs라는 배열에 추출한 job을 업데이트함
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}
	mainC <- jobs
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-jk")
	//fmt.Println(id)
	title := CleanString(card.Find(".title>a").Text())
	//fmt.Println(title)
	location := CleanString(card.Find(".sjcl").Text())
	salary := CleanString(card.Find(".salaryText").Text())
	summary := CleanString(card.Find(".summary").Text())
	c <- extractedJob{
		id:       id,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary}
}

//CelanString cleans a string
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ") //TrimSpace 스페이스를 지운다
	//모든공백제거후 택스트로만 된 배열을 만든것
	//Join은 배열을 가져와서 합치는 역할 (seperator 이용)
}

func getPages(url string) int {
	pages := 0
	res, err := http.Get(url)
	checkErr(err)
	checkCode(res) //response(응답)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	// fmt.Println(doc)
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

func checkErr(err error) {
	if err != nil { //request(의뢰)에 에러가 있다면 이렇게!
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status :", res.StatusCode)
	}
}

//read code and study
//read code and study
//read code and study

/*가장먼저 getPages함수가 실행된다. 총 몇페이지인지 알아야하니까
각페이지별로 getPage가 실행된다 (19페이지)
getPage가 실행되는 중에 extractJob도 실행됨 (50개)
(현재 goroutine으로 바꿀 예정) goroutine이 종료되면 getPage로
채널을 전달함 getPage가 실행이 종료되면 main함수로 채널전송


처음에는 총페이지수를 가져옴 그후에는 각페이지별로 goroutine을 생성
getPage는 각 일자리 정보별로 goroutine을 생성
*/
