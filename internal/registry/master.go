package registry

type MasterRegistry struct {
	Item     ItemRegistry
	Tag      TagRegistry
	Location LocationRegistry
}

func NewMaster(item ItemRegistry, tag TagRegistry, location LocationRegistry) MasterRegistry {
	return MasterRegistry{
		Item:     item,
		Tag:      tag,
		Location: location,
	}
}
