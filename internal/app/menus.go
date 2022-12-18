package app

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"unicode"
)

type App struct {
	client     *mongo.Client
	ctx        context.Context
	db         *mongo.Database
	collection *mongo.Collection
	projection interface{}
	filter     interface{}
	sorting    interface{}
}

func NewApp(ctx context.Context, client *mongo.Client) App {
	return App{
		client: client,
		ctx:    ctx,
	}
}

func (a *App) choice(list []string) int {
	for i, val := range list {
		fmt.Printf("%d. %s\n", i+1, val)
	}

	var curDB int
	_, err := fmt.Scanf("%d", &curDB)
	if err != nil {
		log.Fatal(err)
	}
	if curDB <= 0 || curDB > len(list) {
		log.Fatal("Неккоректный номер")
	}
	return curDB - 1
}

func (a *App) Start() {
	allDatabases, err := a.client.ListDatabaseNames(a.ctx, struct{}{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Выберите базу данных:")
	choice := a.choice(allDatabases)

	a.db = a.client.Database(allDatabases[choice])
	a.DBMenu()
}

func (a *App) DBMenu() {
	allCollections, err := a.db.ListCollectionNames(a.ctx, struct{}{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Выберите коллекцию:")
	choice := a.choice(allCollections)

	a.collection = a.db.Collection(allCollections[choice])
	a.CollectionMenu()
}

func (a *App) getDecision() (bool, error) {
	var decision rune
	_, err := fmt.Scanf("%c\n", &decision)
	if err != nil {
		log.Fatal(err)
	}
	if unicode.ToLower(decision) == 'y' {
		return true, nil
	} else if unicode.ToLower(decision) == 'n' {
		return false, nil
	}
	return false, fmt.Errorf("неверный символ")
}

func (a *App) mustDecision() bool {
	var decision bool
	var err error
	for {
		decision, err = a.getDecision()
		if err == nil {
			break
		} else {
			fmt.Printf(err.Error())
		}
	}
	return decision
}

func (a *App) getProjection() {
	decision := a.mustDecision()
	a.projection = struct{}{}

	if decision {
		fmt.Printf("Введите проекцию:\n")
		err := json.NewDecoder(os.Stdin).Decode(&a.projection)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (a *App) getFilters() {
	decision := a.mustDecision()
	a.filter = struct{}{}

	if decision {
		fmt.Printf("Введите фильтры:\n")
		err := json.NewDecoder(os.Stdin).Decode(&a.filter)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (a *App) getSorting() {
	decision := a.mustDecision()
	a.sorting = struct{}{}

	if decision {
		fmt.Printf("Введите сортировку:\n")
		err := json.NewDecoder(os.Stdin).Decode(&a.sorting)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (a *App) getCertainDocument() {
	decision := a.mustDecision()
	a.sorting = struct{}{}

	var id string
	if decision {
		fmt.Printf("Введите _id документа:\n")
		_, err := fmt.Scanf("%s", &id)
		if err != nil {
			log.Fatal(err)
		}

		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Fatal(err)
		}

		a.sorting = struct{}{}
		a.projection = struct{}{}
		a.filter = bson.D{{"_id", objectId}}

		a.showCollection()
	}
}

func (a *App) showCollection() {
	opts := options.Find().SetProjection(a.projection).SetSort(a.sorting)
	cursor, err := a.collection.Find(a.ctx, a.filter, opts)
	if err != nil {
		log.Fatal(err)
	}
	var all []bson.M

	if err = cursor.All(a.ctx, &all); err != nil {
		log.Fatal(err)
	}
	for _, val := range all {
		fmt.Printf("%v\n", val)
	}
}

func (a *App) CollectionMenu() {
	fmt.Printf("Нужно ли задавать проекцию: y/n\n")
	a.getProjection()
	fmt.Printf("Нужно ли задавать фильтры: y/n\n")
	a.getFilters()
	fmt.Printf("Нужно ли задавать сортировку: y/n\n")
	a.getSorting()
	a.showCollection()
	fmt.Printf("Хотите выбрать конкретный документ: y/n\n")
	a.getCertainDocument()
	fmt.Printf("Хотите еще раз?: y/n\n")
	a.reload()
}

func (a *App) reload() {
	decision := a.mustDecision()

	if decision {
		a.sorting = struct{}{}
		a.filter = struct{}{}
		a.projection = struct{}{}
		a.Start()
	}
}
