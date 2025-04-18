package service

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"github.com/samber/lo"

	"github.com/google/uuid"
)

func (s *Service) AddItem(astralID uuid.UUID, itemCode string, dataTags map[string]string) (*schema.Item, error) {
	item := model.Item{
		AstralID: astralID,
		Code:     itemCode,
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
	itemSchema := schema.ItemSchemaFromItem(&item)
	return itemSchema, err
}

func (s *Service) FindOneItem(ID uuid.UUID) (*schema.Item, error) {
	item, err := s.i.FindOne(ID)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, util.ErrItemNotFound
	}

	itemSchema := schema.ItemSchemaFromItem(item)
	return itemSchema, nil
}

func (s *Service) FindOneItemByCode(code string) (*schema.Item, error) {
	item, err := s.i.FindOneByCode(code)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, util.ErrItemNotFound
	}

	itemSchema := schema.ItemSchemaFromItem(item)
	return itemSchema, nil
}

func (s *Service) FindAllItems(filter *model.Item) ([]schema.Item, error) {
	items, err := s.i.FindAll(filter)
	if err != nil {
		return nil, err
	}
	itemSchemas := lo.Map(items, func(item model.Item, index int) schema.Item {
		return *schema.ItemSchemaFromItem(&item)
	})
	return itemSchemas, nil
}

func (s *Service) GetItemDataTags(itemID uuid.UUID) map[string]string {
	tags, err := s.idt.FindAll(&model.ItemDataTag{ItemID: itemID})
	if err != nil {
		return map[string]string{}
	}

	var res = map[string]string{}
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

func (s *Service) GetAstralItems(astralID uuid.UUID) []schema.Item {
	items, err := s.i.FindAll(&model.Item{AstralID: astralID})
	if err != nil {
		return nil
	}

	var itemSchemas []schema.Item
	for _, i := range items {
		itemSchemas = append(itemSchemas, *schema.ItemSchemaFromItem(&i))
	}
	return itemSchemas
}
