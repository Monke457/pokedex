package pokedex

import "iter"

type Pokedex struct {
	pokemon map[string]Pokemon
}

func New() Pokedex {
	return Pokedex{ pokemon: map[string]Pokemon{} }
}

func (p *Pokedex) Add(pokemon Pokemon) {
	p.pokemon[pokemon.Name] = pokemon
}

func (p Pokedex) Get(name string) (Pokemon, bool) {
	poke, ok := p.pokemon[name] 
	return poke, ok
}

func (p Pokedex) List() iter.Seq[string] {
	return func(yield func(string) bool) {
		for k := range p.pokemon {
			if !yield(k) {
				return
			}
		}
	}
}
