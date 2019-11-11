const col = "rates";

db.createCollection(col, {
	validator: {
		$jsonSchema: {
			bsonType: "object",
			required: [ "currency", "exchangeRate" ],
			properties: {
				currency: {
					bsonType: "string",
					enum: [ "EUR", "USD" ],
					description: "must be a string, is required and is ENUM - [ 'USD', 'EUR' ]"
				},
				exchangeRate: {
					bsonType: "double",
					description: "must be a double is required"
				}
			}
		}
	}
});

db[col].createIndex("currency", {unique: true});
db[col].createIndex("exchangeRate");