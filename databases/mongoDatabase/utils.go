package mongoDatabase

import "go.mongodb.org/mongo-driver/bson"

func toBsonDoc(object any) (doc *bson.D, err error) {
	data, err := bson.Marshal(object)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}
