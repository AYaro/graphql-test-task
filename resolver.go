package graphql_test_task

import (
	"context"
	"log"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

var rateUpdatedUSDChannel map[string]chan *Rate
var rateUpdatedEURChannel map[string]chan *Rate

func init() {
	rateUpdatedUSDChannel = map[string]chan *Rate{}
	rateUpdatedEURChannel = map[string]chan *Rate{}
}

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) UpdateRate(ctx context.Context, input RateInput) (bool, error) {

	filter := bson.D{{"currency", input.Currency}}
	update := bson.D{
		{"$set", bson.D{
			{"exchangeRate", input.ExchangeRate},
		}},
	}

	//обновляем курс в бд
	res, err := collection.UpdateOne(context.TODO(), filter, update)

	//если вдруг такая валюта не найдена то возвращаем false
	if res.MatchedCount == 0 || err != nil {
		err = initiateCurrencies()
		return false, err
	}

	//передаем новый курс в каналы с валютой input.Currency
	if input.Currency == CurrencyEur {
		for _, elem := range rateUpdatedEURChannel {
			elem <- &Rate{input.Currency, input.ExchangeRate}
		}
	} else {
		for _, elem := range rateUpdatedUSDChannel {
			elem <- &Rate{input.Currency, input.ExchangeRate}
		}
	}
	return true, err
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) GetRates(ctx context.Context, currency *Currency) ([]*Rate, error) {

	var result []*Rate

	//если указана валюта то выполняем
	if currency != nil {
		filter := bson.D{{"currency", currency}}

		var singleResult *Rate

		err := collection.FindOne(ctx, filter).Decode(&singleResult)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, singleResult)
		return result, err

	}

	//если не указана валюта, то выполняем для всех
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(ctx) {

		var elem Rate
		err = cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, &elem)
	}

	cur.Close(ctx)

	return result, err
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) ObserveRate(ctx context.Context, currency Currency) (<-chan *Rate, error) {
	id := randomString(8)
	subEvents := make(chan *Rate, 1)
	var channels map[string]chan *Rate

	//определяем куда добавлять канал
	switch currency {
	case CurrencyEur:
		channels = rateUpdatedEURChannel
	default:
		channels = rateUpdatedUSDChannel
	}

	go func() {
		<-ctx.Done()
		delete(channels, id)
	}()

	//добовляем канал
	channels[id] = subEvents

	return subEvents, nil

}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
