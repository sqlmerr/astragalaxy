package services

import (
	"astragalaxy/internal/models"
	"astragalaxy/internal/repositories"
	"astragalaxy/internal/schemas"
	"astragalaxy/internal/utils"

	"github.com/google/uuid"
)

type ItemService struct {
	i repositories.ItemRepository
	d repositories.ItemDataTagRepository
	// r registry.ItemRegistry
}

func NewItemService(i repositories.ItemRepository, d repositories.ItemDataTagRepository) ItemService {
	return ItemService{i: i, d: d}
}

func (s *ItemService) AddItem(userID uuid.UUID, itemCode string, dataTags map[string]string) (*schemas.ItemSchema, error) {
	item := models.Item{
		UserID: userID,
		Code:   itemCode,
	}
	err := s.i.Create(&item)
	if err != nil {
		return nil, err
	}

	for k, v := range dataTags {
		tag := models.ItemDataTag{
			ItemID: item.ID,
			Key:    k,
			Value:  v,
		}
		err = s.d.Create(&tag)
		if err != nil {
			return nil, err
		}
	}
	schema := schemas.ItemSchemaFromItem(&item)
	return schema, err
}

func (s *ItemService) FindOne(ID uuid.UUID) (*schemas.ItemSchema, error) {
	item, err := s.i.FindOne(ID)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, utils.ErrItemNotFound
	}

	schema := schemas.ItemSchemaFromItem(item)
	return schema, nil
}

func (s *ItemService) FindOneByCode(code string) (*schemas.ItemSchema, error) {
	item, err := s.i.FindOneByCode(code)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, utils.ErrItemNotFound
	}

	schema := schemas.ItemSchemaFromItem(item)
	return schema, nil
}

func (s *ItemService) GetItemDataTags(itemID uuid.UUID) map[string]string {
	tags, err := s.d.FindAll(&models.ItemDataTag{ItemID: itemID})
	if err != nil {
		return map[string]string{}
	}

	res := make(map[string]string)
	for _, t := range tags {
		res[t.Key] = t.Value
	}
	return res
}

func (s *ItemService) GetItemTag(itemID uuid.UUID, key string) (*string, error) {
	tag, err := s.d.FindOneByFilter(&models.ItemDataTag{ItemID: itemID, Key: key})
	if err != nil {
		return nil, err
	} else if tag == nil {
		return nil, utils.ErrItemDataTagNotFound
	}

	val := tag.Value
	return &val, nil
}

func (s *ItemService) GetUserItems(userID uuid.UUID) []schemas.ItemSchema {
	items, err := s.i.FindAll(&models.Item{UserID: userID})
	if err != nil {
		return nil
	}

	var itemSchemas []schemas.ItemSchema
	for _, i := range items {
		itemSchemas = append(itemSchemas, *schemas.ItemSchemaFromItem(&i))
	}
	return itemSchemas
}
