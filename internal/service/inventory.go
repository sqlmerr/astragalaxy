package service

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"

	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

func (s *Service) CreateInventory(holder string, holderID uuid.UUID) (*model.Inventory, error) {
	inventory := &model.Inventory{
		Holder: holder, HolderID: holderID,
	}
	err := s.inv.Create(inventory)
	if err != nil {
		return nil, err
	}

	_, err = s.b.Create(s.txManager.DB(), &model.Bundle{InventoryID: inventory.ID, Resources: map[string]int{}})
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

func (s *Service) GetInventoryByHolder(holder string, holderID uuid.UUID) (*model.Inventory, error) {
	inventory, err := s.inv.FindOneByFilter(&model.Inventory{Holder: holder, HolderID: holderID})
	if err != nil {
		return nil, err
	}
	if inventory == nil {
		return nil, util.ErrNotFound
	}

	return inventory, nil
}

func (s *Service) FindOneInventory(id uuid.UUID) (*model.Inventory, error) {
	return s.inv.FindOne(id)
}

func (s *Service) AddItem(inventoryID uuid.UUID, itemCode string, dataTags map[string]string) (*schema.Item, error) {
	item := model.Item{
		InventoryID: inventoryID,
		Code:        itemCode,
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

func (s *Service) AddItemToAstral(astralID uuid.UUID, itemCode string, dataTags map[string]string) (*schema.Item, error) {
	inventory, err := s.GetInventoryByHolder("astral", astralID)
	if err != nil {
		return nil, err
	}
	return s.AddItem(inventory.ID, itemCode, dataTags)
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

func (s *Service) GetInventoryItems(inventoryID uuid.UUID) ([]schema.Item, error) {
	items, err := s.i.FindAll(&model.Item{InventoryID: inventoryID})
	if err != nil {
		return nil, err
	}
	itemSchemas := lo.Map(items, func(item model.Item, index int) schema.Item {
		return *schema.ItemSchemaFromItem(&item)
	})
	return itemSchemas, nil
}

func (s *Service) GetAstralItems(astralID uuid.UUID) ([]schema.Item, error) {
	astralInventory, err := s.GetInventoryByHolder("astral", astralID)
	if err != nil {
		return nil, err
	}
	return s.GetInventoryItems(astralInventory.ID)
}

func (s *Service) SendItem(itemID uuid.UUID, fromID uuid.UUID, toID uuid.UUID) error {
	inventoryFrom, err := s.inv.FindOne(fromID)
	if err != nil {
		return err
	}
	inventoryTo, err := s.inv.FindOne(toID)
	if err != nil {
		return err
	}

	item, err := s.FindOneItem(itemID)
	if err != nil {
		return err
	}

	if item.InventoryID != inventoryFrom.ID {
		return util.ErrItemNotFound
	}

	return s.i.Update(&model.Item{
		InventoryID: inventoryTo.ID,
	})
}

func (s *Service) UpdateItemDataTags(itemID uuid.UUID, dataTags map[string]string) error {
	for k, v := range dataTags {
		tag, err := s.idt.FindOneByFilter(&model.ItemDataTag{ItemID: itemID, Key: k})
		if err != nil {
			return err
		}

		// TODO: transactions
		if tag == nil {
			err = s.idt.Create(&model.ItemDataTag{ItemID: itemID, Key: k, Value: v})
		} else {
			err = s.idt.Update(&model.ItemDataTag{ID: tag.ID, Value: v})
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) UpdateItem(data *model.Item) error {
	return s.i.Update(data)
}

func (s *Service) UpdateItemRaw(itemID uuid.UUID, data map[string]any) error {
	return s.i.UpdateRaw(itemID, data)
}

func (s *Service) DeleteItem(itemID uuid.UUID) error {
	return s.i.Delete(itemID)
}

func (s *Service) GetAstralFromInventory(inventory *model.Inventory) (*schema.Astral, error) {
	var astralID uuid.UUID
	switch inventory.Holder {
	case "astral":
		astralID = inventory.HolderID
	case "spaceship":
		spaceship, err := s.FindOneSpaceship(inventory.HolderID)
		if err != nil {
			return nil, err
		}
		astralID = spaceship.AstralID
	default:
		return nil, util.ErrNotFound
	}

	return s.FindOneAstral(astralID)
}

// --- Bundle methods ---

func (s *Service) FindOneBundle(id uuid.UUID) (*model.Bundle, error) {
	return s.b.FindOne(s.txManager.DB(), id)
}

func (s *Service) GetBundleByInventory(inventoryID uuid.UUID) (*model.Bundle, error) {
	bundle, err := s.b.FindOneByFilter(s.txManager.DB(), &model.Bundle{InventoryID: inventoryID})
	if err != nil {
		return nil, err
	}
	if bundle == nil {
		return nil, util.ErrNotFound
	}
	return bundle, nil
}

func (s *Service) InsertOrUpdateBundleResources(bundleID uuid.UUID, resources map[string]int) error {
	// Validate resources
	for k := range resources {
		r, err := s.r.Resource.FindOne(k)
		if err != nil || r == nil {
			return util.ErrNotFound
		}
	}

	return s.txManager.WithTx(func(tx *gorm.DB) error {
		bundle, err := s.b.FindOne(tx, bundleID)
		if err != nil {
			return err
		}
		if bundle == nil {
			return util.ErrNotFound
		}

		currentResources := bundle.Resources
		for k, v := range resources {
			currentResources[k] = v
		}
		err = s.b.Update(tx, &model.Bundle{ID: bundleID, Resources: currentResources})
		if err != nil {
			return err
		}
		return nil
	})
}
