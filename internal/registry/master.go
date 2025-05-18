package registry

type MasterRegistry struct {
	Item     ItemRegistry
	Tag      TagRegistry
	Location LocationRegistry
	Resource ResourceRegistry
}

func NewMaster(item ItemRegistry, tag TagRegistry, location LocationRegistry, resource ResourceRegistry) MasterRegistry {
	return MasterRegistry{
		Item:     item,
		Tag:      tag,
		Location: location,
		Resource: resource,
	}
}
