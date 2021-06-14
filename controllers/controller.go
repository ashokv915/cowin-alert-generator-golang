package controller

import (
	"context"
	"cowin-alert/database"
	"cowin-alert/models"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	TELEGRAM_URL = "https://api.telegram.org/<KEY>/sendMessage?chat_id=@<GRP_NAME>&text="
)

func ExtractDetails(Sessions []models.Session) {
	//TelegramMessage("Starting To Extract")
	var center models.Centers
	fmt.Println(Sessions[0].Date)
	for i,v := range Sessions {
		//message := fmt.Sprintf("Name: %s Date: %s Dose1: %d Dose2: %d",v.Name,v.Date,v.AvailableCapacityDose1,v.AvailableCapacityDose2)
		log.Println(v.Name, i)
		if Find(v.Pincode) {
			log.Println("Picode Matched:",v.Pincode,v.Name)
			if (v.AvailableCapacity > 0 ) {
				err := database.COLLECTION.FindOne(context.TODO(), models.Centers{CenterID: v.Center_id}).Decode(&center)
				if err != nil {
					log.Println(err)
				}
				pastCountDose1 := center.Dose1
				pastCountDose2 := center.Dose2
				presentCountDose1 := v.AvailableCapacityDose1
				presentCountDose2 := v.AvailableCapacityDose2

				if (pastCountDose1 != int64(presentCountDose1)) || (pastCountDose2 != int64(presentCountDose2)) {
					log.Println("Count Changed")
					message := fmt.Sprintf("Date: %s\nPlace: %s\nDose1: %d\nDose2: %d\nMinimum Age: %d\nVaccine: %s\nFee: %s\nPincode: %d\n \n \n \n",v.Date,v.Name,v.AvailableCapacityDose1,v.AvailableCapacityDose2,v.MinAge,v.Vaccine,v.FeeType,v.Pincode,"https://selfregistration.cowin.gov.in/")
					TelegramMessage(message)
					log.Println("Message sent to Telegram successfully")
					UpdateDB(v.Center_id,int32(presentCountDose1),int32(presentCountDose2))
				}
			}
		}
	}
}

func Find(val int64) bool {
	NEAREST_PINCODE := [11]int64{679534, 678633, 679313, 679301, 679307, 679503, 679308, 679303, 679309, 679102, 679121}
    for _, item := range NEAREST_PINCODE {
        if item == val {
            return true
        }
    }
    return false
}

func UpdateDB(centerid int64, dose1 int32,dose2 int32) {
	_,err := database.COLLECTION.UpdateMany(
		context.TODO(),
		bson.M{"centerid":centerid},
		bson.D{
			{"$set", bson.D{{
				Key:   "dose1",
				Value: dose1,
			},
			{
				Key: "dose2",
				Value: dose2,	
			},
		}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func TelegramMessage(message string) {
	finalURL := TELEGRAM_URL + message
	_,err := http.Get(finalURL)
	if err != nil {
		log.Fatal(err)
	}

}

func AddCenterToDB(Sessions []models.Session) {
	for _,v := range Sessions {
		centerIDFromResponse := v.Center_id
		result := CheckInDB(centerIDFromResponse)
		log.Println(result)
		if result != true {
			log.Println("Inserting to DB")
			InsertIntoDB(v)
		}

	}
}

func CheckInDB(centerID int64) bool {
	var center models.Centers
	err := database.COLLECTION.FindOne(context.TODO(), models.Centers{CenterID: centerID}).Decode(&center)
	if err != nil {
		log.Println("No Record Found",err)
		return false
	}
	log.Println("Record Found", center.CenterID)
	return true
}

func InsertIntoDB(center models.Session) {
	elem := models.Centers{center.Center_id, int64(center.AvailableCapacityDose1), int64(center.AvailableCapacityDose2)}
	_,err := database.COLLECTION.InsertOne(context.TODO(),elem)
	if err != nil {
		log.Fatal("Could not insert into database",err)
	}
	log.Println("Inserted successfully",center.Center_id)
}