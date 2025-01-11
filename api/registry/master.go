package registry

type MasterRegistry struct {
	Item ItemRegistry
	Tag  TagRegistry
}

func NewMaster(item ItemRegistry, tag TagRegistry) MasterRegistry {
	return MasterRegistry{
		Item: item,
		Tag:  tag,
	}
}
