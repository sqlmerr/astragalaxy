package util

import (
	"fmt"
	"math/rand"
)

var LatinPlanetNames = []string{
	"Aether", "Albus", "Altus", "Amicus", "Amor", "Anima",
	"Aqua", "Arbor", "Arcus", "Ars", "Astra", "Aurora",
	"Aurum", "Alpha", "Bellum", "Beta", "Caelum", "Caligo", "Calor",
	"Campus", "Candor", "Carmen", "Carpe", "Celer",
	"Cernuo", "Cibus", "Cinis", "Civis", "Clarus",
	"Clemens", "Coelum", "Cogito", "Conexus", "Consilium",
	"Corona", "Corpus", "Cresco", "Cura", "Cursus",
	"Decus", "Deus", "Dies", "Dignitas", "Discipulus",
	"Divinus", "Dominus", "Donum", "Dulcis", "Dux",
	"Elementum", "Elysium", "Emendo", "Enigma", "Ensis",
	"Equus", "Erebus", "Espero", "Fama", "Fatum",
	"Fides", "Finis", "Flamma", "Flos", "Flumen",
	"Fortis", "Frater", "Fuga", "Fulmen", "Furor",
	"Gaia", "Genius", "Glacies", "Gloria", "Gratia",
	"Gravis", "Helios", "Hiems", "Honor", "Hora",
	"Humus", "Ignis", "Imago", "Imperium", "Inceptum",
	"Infinitus", "Ingenium", "Initium", "Ira", "Iunctus",
	"Iustitia", "Labor", "Laurus", "Legio", "Lethum",
	"Libertas", "Limes", "Luceo", "Lumen", "Luna",
	"Lux", "Magister", "Magnus", "Mare", "Mater",
	"Memoria", "Mens", "Mira", "Mors", "Mundus",
	"Natura", "Nebula", "Nex", "Nobilis", "Nomen",
	"Notus", "Novus", "Nox", "Nubes", "Oculus",
	"Omnia", "Opus", "Orbis", "Ordo", "Pallor",
	"Pax", "Perpetuus", "Petra", "Pietas", "Pluvia",
	"Pons", "Populus", "Potens", "Primus", "Proelium",
	"Pulcher", "Purus", "Quaero", "Quies", "Radius",
	"Regnum", "Rex", "Sacer", "Sanguis", "Sapientia",
	"Scintilla", "Sensus", "Serenus", "Sidus", "Signum",
	"Silva", "Sol", "Solus", "Somnus", "Sors",
	"Spes", "Spiritus", "Stella", "Summus", "Tactus",
	"Tellus", "Tempus", "Tenebrae", "Terra", "Thronus",
	"Tigris", "Tonitrus", "Trans", "Tristis", "Ultimus",
	"Umbra", "Unitas", "Universum", "Urbs", "Valens",
	"Vates", "Ventus", "Veritas", "Verus", "Vesper",
	"Via", "Victoria", "Vigor", "Vita", "Vox",
	"Vulnus", "Zephyrus",
}

func GeneratePlanetName() string {
	word1 := LatinPlanetNames[rand.Intn(len(LatinPlanetNames))]
	word2 := LatinPlanetNames[rand.Intn(len(LatinPlanetNames))]

	letters := []string{"A", "B", "C", "D", "E", "F", "G"}
	letter := letters[rand.Intn(len(letters))]

	number := rand.Intn(8) + 1

	return fmt.Sprintf("%s %s %s%d", word1, word2, letter, number)
}
