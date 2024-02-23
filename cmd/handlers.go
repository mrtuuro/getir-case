package cmd

import (
	"context"
	"encoding/json"
	"getir-arac/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

// request struct is created to unmarshal request json body to Go's struct type.
type request struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	MinCount  int64  `json:"minCount"`
	MaxCount  int64  `json:"maxCount"`
}

// response struct is to return any response to client. It is designed according to client's needs.
type response struct {
	Code    int             `json:"code"`
	Message string          `json:"msg"`
	Records []models.Record `json:"records"`
}

var (
	// Specified error messages for response body/error codes.
	Success          = 0
	ErrInternalError = 1
	ErrNotFound      = 2
	ErrBadRequest    = 3
)

func (a *App) HandleGetRecords(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Check for HTTP Request Method
	if r.Method != "POST" {
		respondWithJSON(w, http.StatusMethodNotAllowed, response{
			Code:    ErrBadRequest,
			Message: "Method Not Allowed",
			Records: nil,
		})
		return
	}

	// Creating a request variable to store request body parameters.
	var req request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithJSON(w, http.StatusInternalServerError, response{
			Code:    ErrInternalError,
			Message: "errors.InternalError",
			Records: nil,
		})
		return
	}

	// Parsing request body parameter date types.
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, response{
			Code:    ErrBadRequest,
			Message: "errors.BadDateFormat",
			Records: nil,
		})

		return
	}

	// Parsing request body parameter date types.
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, response{
			Code:    ErrBadRequest,
			Message: "errors.BadDateFormat",
			Records: nil,
		})
		return
	}

	// Creating a pipeline to filter documents that is stored in MongoDB.
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.M{
			// We filter to get document that is created between start and end date.
			"createdAt": bson.M{"$gte": startDate, "$lte": endDate},
		}}},
		bson.D{{"$project", bson.M{
			// Here we sum up all the count fields and check for if it is between the wanted scope.
			"totalCount": bson.M{
				"$sum": "$counts"},
			"createdAt": 1,
			"key":       1,
		}}},
		bson.D{{"$match", bson.M{
			"totalCount": bson.M{"$gte": req.MinCount, "$lte": req.MaxCount},
		}}},
	}

	coll := a.Config.MongoClient.Database("getircase-study").Collection("records")

	// Filtering documents according to our pipeline
	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, response{
			Code:    ErrInternalError,
			Message: "errors.InternalError",
			Records: nil,
		})
		return
	}

	// We need to close the cursor at the end of the process in case data leaks.
	defer cursor.Close(ctx)

	// We create an array of Record to store the incoming datas from filter.
	var records []models.Record

	// Iterating and decoding documents to records array.
	if err = cursor.All(ctx, &records); err != nil {
		respondWithJSON(w, http.StatusInternalServerError, response{
			Code:    ErrInternalError,
			Message: "errors.InternalError",
			Records: nil,
		})
		return
	}

	// Check if there is no documents returned.
	if len(records) == 0 {
		respondWithJSON(w, http.StatusNotFound, response{
			Code:    ErrNotFound,
			Message: "errors.NotFound",
			Records: nil,
		})
		return
	}

	// Respond with success message and retrieved data.
	respondWithJSON(w, http.StatusOK, response{
		Code:    Success,
		Message: "Success",
		Records: records,
	})
}

func (a *App) HandleGetInMemory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	key := r.URL.Query().Get("key")
	value, ok := a.Config.MemoryDB[key]
	if !ok {
		respondWithError(w, http.StatusNotFound, "not found!")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"key":   key,
		"value": value,
	})
}

func (a *App) HandleInsertInMemory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	type request struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	var req request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "can not decode request body!")
		return
	}
	a.Config.MemoryDB[req.Key] = req.Value
	respondWithJSON(w, http.StatusCreated, map[string]string{
		"key":   req.Key,
		"value": a.Config.MemoryDB[req.Key],
	})
}
