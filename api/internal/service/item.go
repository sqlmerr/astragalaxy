package service

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schemas"
	"astragalaxy/internal/util"

	"github.com/google/uuid"
)

func (s *Service) AddItem(userID uuid.UUID, itemCode string, dataTags map[string]string) (*schemas.ItemSchema, error) {
	item := model.Item{
		UserID: userID,
		Code:   itemCode,
	}
	err := s.i.Create(&item)
	if err != nil {
		return nil, err
	}

	for k, v := range dataTags {
		tag := model.ItemDataTag{
			ItemID: item.ID,
			Key:    k,
			Value:  v,
		}
		err = s.idt.Create(&tag)
		if err != nil {
			return nil, err
		}
	}
	schema := schemas.ItemSchemaFromItem(&item)
	return schema, err
}

func (s *Service) FindOneItem(ID uuid.UUID) (*schemas.ItemSchema, error) {
	item, err := s.i.FindOne(ID)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, util.ErrItemNotFound
	}

	schema := schemas.ItemSchemaFromItem(item)
	return schema, nil
}

func (s *Service) FindOneItemByCode(code string) (*schemas.ItemSchema, error) {
	item, err := s.i.FindOneByCode(code)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, util.ErrItemNotFound
	}

	schema := schemas.ItemSchemaFromItem(item)
	return schema, nil
}

func (s *Service) GetItemDataTags(itemID uuid.UUID) map[string]string {
	tags, err := s.idt.FindAll(&model.ItemDataTag{ItemID: itemID})
	if err != nil {
		return map[string]string{}
	}

	res := make(map[string]string)
	for _, t := range tags {
		res[t.Key] = t.Value
	}
	return res
}

func (s *Service) GetItemDataTag(itemID uuid.UUID, key string) (*string, error) {
	tag, err := s.idt.FindOneByFilter(&model.ItemDataTag{ItemID: itemID, Key: key})
	if err != nil {
		return nil, err
	} else if tag == nil {
		return nil, util.ErrItemDataTagNotFound
	}

	val := tag.Value
	return &val, nil
}

func (s *Service) GetUserItems(userID uuid.UUID) []schemas.ItemSchema {
	items, err := s.i.FindAll(&model.Item{UserID: userID})
	if err != nil {
		return nil
	}

	var itemSchemas []schemas.ItemSchema
	for _, i := range items {
		itemSchemas = append(itemSchemas, *schemas.ItemSchemaFromItem(&i))
	}
	return itemSchemas
}
