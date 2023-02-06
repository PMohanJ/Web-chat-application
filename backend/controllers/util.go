package controllers

import (
	"go.mongodb.org/mongo-driver/bson"
)

func LookUpStage(from, localField, foreignField, as string) bson.D {
	return bson.D{
		{
			"$lookup", bson.D{
				{"from", from},
				{"localField", localField},
				{"foreignField", foreignField},
				{"as", as},
			},
		},
	}
}
