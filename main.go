package main

import (
	"fmt"
	"os"
	"encoding/csv"
	"io"
	"strings"
	"strconv"
	"github.com/goJtlParser/round"
	"io/ioutil"
	"github.com/montanaflynn/stats"
)


var responseTime 	 	= [...] int { 10,20,30,40,50,60,70,80,90,100,200,300,400,500,600,700,800,900,1000,2000,3000,4000,5000 }
var reponseTimecount 	= [...] float64 {  0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0 }
var responseCodes 		= [...] string {"1xx","2xx", "3xx","4xx","5xx","6xx"}
var responseCodescount 	= [...] int {0,0,0,0,0,0}
var yAxisBarChat 		= [...] string {"0-100", "100-200", "200-300", "300-400", "400-500", "500-600", "600-700", "700-800", "800-900", "900-1000", "1000-1200", "1200-1400", "1400-1600", "1600-1800", "1800-2000", "2000 <"}
var xAxisBarChat 		= [...] int {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var elapsedTime			[]float64


var barChartCsv = "BarChart.csv";
var lineChart = "lineChart.csv";
var pieChartCsv = "pieChart.csv";
var infoTestCsv = "TestInfo.csv";
var percentileFile = "percentile.csv";
var resTimeDistribution = "resTimeDistribution.csv";
var summaryReport = "summaryReport.csv";

func main() {
	/*
	if len(os.Args) < 2 {
		fmt.Printf("Error: Source file name is required\n")
		fmt.Println("Usage:", os.Args[0], "<filename> \n")
		return
	}
	*/

	//file, err := os.Open(os.Args[1])
	file, err := os.Open("/Users/abhijit.p/Downloads/apache-jmeter-3.2/bin/temp.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','          //field delimiter
	reader.Comment = '#'        //Comment character
	reader.FieldsPerRecord = -1 //Number of records per record. Set to Negative value for variable
	reader.TrimLeadingSpace = true

	transactionsCount := 0.0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if strings.Contains(strings.Join(record, ""),"elapsed") {
			continue;
		}
		if len(record) == 0 {
			continue;
		}
		transactionsCount = transactionsCount+1
		rt,_ := strconv.ParseFloat(record[1], 64)
		elapsedTime = append(elapsedTime, rt)
		setResponseTimeDistribution(rt)
		setResponseCodeDistribution(record[3])


	}
	createCSVForResponseTimeDistribution(transactionsCount)
	createCSVForPercentileDistribution()
	createCSVForResponseCode()



}
func createCSVForResponseCode() {
	var strs []string
	strs = append(strs, "responseCode,count")
	for i:=0;i<len(responseCodes)-1 ;i++  {
		strs = append(strs,	responseCodes[i] +","+ strconv.Itoa(responseCodescount[i]))
	}

	err := ioutil.WriteFile("/tmp/"+pieChartCsv, []byte(strings.Join(strs, "\n")), 0644)
	if(err != nil){
		panic(err)
	}
}


func createCSVForPercentileDistribution() {
	var strs []string
	strs = append(strs, "Percentile Level,Value")

	percentile25,_ :=stats.Percentile([]float64(elapsedTime), 25.0)
	strs = append(strs,	"25 percentile value:,"+strconv.FormatFloat(round.RoundDown(percentile25, 0), 'f', -1, 64)+" ms")

	percentile50,_ :=stats.Percentile([]float64(elapsedTime), 50.0)
	strs = append(strs,	"50 percentile value:,"+strconv.FormatFloat(round.RoundDown(percentile50,0), 'f', -1, 64)+" ms")

	percentile75,_ :=stats.Percentile([]float64(elapsedTime), 75.0)
	strs = append(strs,	"75 percentile value:,"+strconv.FormatFloat(round.RoundDown(percentile75,0), 'f', -1, 64)+" ms")

	percentile80,_ :=stats.Percentile([]float64(elapsedTime), 80.0)
	strs = append(strs,	"80 percentile value:,"+strconv.FormatFloat(round.RoundDown(percentile80,0), 'f', -1, 64)+" ms")

	percentile90,_ :=stats.Percentile([]float64(elapsedTime), 90.0)
	strs = append(strs,	"90 percentile value:,"+strconv.FormatFloat(round.RoundDown(percentile90, 0), 'f', -1, 64)+" ms")

	percentile95,_ :=stats.Percentile([]float64(elapsedTime), 95.0)
	strs = append(strs,	"95 percentile value:,"+strconv.FormatFloat(round.RoundDown(percentile95, 0), 'f', -1, 64)+" ms")

	percentile98,_ :=stats.Percentile([]float64(elapsedTime), 98.0)
	strs = append(strs,	"98 percentile value:,"+strconv.FormatFloat(round.RoundDown(percentile98, 0), 'f', -1, 64)+" ms")

	percentile99,_ :=stats.Percentile([]float64(elapsedTime), 99.0)
	strs = append(strs,	"99 percentile value:,"+strconv.FormatFloat(round.RoundDown(percentile99, 0), 'f', -1, 64)+" ms")

	percentile100,_ :=stats.Percentile([]float64(elapsedTime), 100.0)
	strs = append(strs,	"100 percentile value:,"+strconv.FormatFloat(round.RoundDown(percentile100,0), 'f', -1, 64)+" ms")

	err := ioutil.WriteFile("/tmp/"+percentileFile, []byte(strings.Join(strs, "\n")), 0644)
	if(err != nil){
		panic(err)
	}

}


func createCSVForResponseTimeDistribution(transactionsCount float64){
	var strs []string
	strs = append(strs, "Timing(ms), Count, Fraction, Rolling Fraction")
	fraction := round.RoundDown((reponseTimecount[0]/transactionsCount) * 100, 3);
	rolloingFrac := round.RoundDown((reponseTimecount[0]/transactionsCount) * 100, 3);

	j:=0
	for i:=0;i<len(responseTime)-1 ;i++  {
		if(i>=1){
			strs = append(strs, strconv.Itoa(responseTime[j-1])+"-"+ strconv.Itoa(responseTime[i])+ ","+strconv.FormatFloat(reponseTimecount[i],'f', -1, 64)+ ","+ strconv.FormatFloat(fraction,'f', -1, 64)+"%,"+strconv.FormatFloat(rolloingFrac,'f', -1, 64)+"%")
		}else {
			strs = append(strs, strconv.Itoa(0)+"-"+ strconv.Itoa(responseTime[i])+ ","+ strconv.FormatFloat(reponseTimecount[i],'f', -1, 64)+ ","+ strconv.FormatFloat(fraction,'f', -1, 64)+"%,"+strconv.FormatFloat(rolloingFrac,'f', -1, 64)+"%")
		}

		temp :=rolloingFrac
		fraction = round.RoundDown((reponseTimecount[i+1]/transactionsCount) * 100, 3);
		rolloingFrac = round.RoundDown(temp+fraction, 3);
		j++
	}

	strs = append(strs, "5000+"+ ","+strconv.FormatFloat(reponseTimecount[len(reponseTimecount)-1],'f', -1, 64)+ ","+ strconv.FormatFloat(fraction,'f', -1, 64)+"%,"+strconv.FormatFloat(rolloingFrac,'f', -1, 64)+"%")


	err := ioutil.WriteFile("/tmp/"+resTimeDistribution, []byte(strings.Join(strs, "\n")), 0644)
	if(err != nil){
		panic(err)
	}
}

func setResponseTimeDistribution(d float64){

	if (d > 0) && (d <= 10){
		reponseTimecount[0]++
	}else if (d > 10) && (d <= 20){
		reponseTimecount[1]++
	}else if (d > 20) && (d <= 30){
		reponseTimecount[2]++
	} else if (d > 30) && (d <= 40){
		reponseTimecount[3]++
	}else if (d > 40) && (d <= 50){
		reponseTimecount[4]++
	} else if (d > 50) && (d <= 60){
		reponseTimecount[5]++
	} else if (d > 60) && (d <= 70) {
		reponseTimecount[6]++
	} else if(d > 70 && d <= 80){
		reponseTimecount[7]++
	} else if(d > 80 && d <= 90){
		reponseTimecount[8]++
	} else if(d > 90 && d <= 100){
		reponseTimecount[9]++
	} else if(d > 100 && d <= 200){
		reponseTimecount[10]++
	} else if(d > 200 && d <= 300){
		reponseTimecount[11]++
	} else if(d > 300 && d <= 400){
		reponseTimecount[12]++
	} else if(d > 400 && d <= 500){
		reponseTimecount[13]++
	} else if(d > 500 && d <= 600){
		reponseTimecount[14]++
	} else if(d > 600 && d <= 700){
		reponseTimecount[15]++
	} else if(d > 700 && d <= 800){
		reponseTimecount[16]++
	} else if(d > 800 && d <= 900){
		reponseTimecount[17]++
	} else if(d > 900 && d <= 1000){
		reponseTimecount[18]++
	} else if(d > 1000 && d <= 2000){
		reponseTimecount[19]++
	} else if(d > 2000 && d <= 3000){
		reponseTimecount[20]++
	} else if(d > 3000 && d <= 4000){
		reponseTimecount[21]++
	} else if(d > 4000 && d <= 5000){
		reponseTimecount[22]++
	} else if(d > 5000){
		reponseTimecount[23]++
	}
}

func setResponseCodeDistribution(d string) {

	if (strings.Contains(d, "java.net.SocketException")) {
		responseCodescount[5]++;
		return
	}

	code, _ := strconv.Atoi(d)
	if (code > 100) && (code < 200) {
		responseCodescount[0]++
	} else if (code >= 200) && (code < 300) {
		responseCodescount[1]++
	} else if (code >= 300) && (code < 400) {
		responseCodescount[2]++
	} else if (code >= 400) && (code < 500) {
		responseCodescount[3]++
	} else if (code >= 500) && (code < 600) {
		responseCodescount[4]++
	}
}