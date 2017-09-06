package mazenv

import "testing"

func TestGenerators(t *testing.T) {
	gens := map[string]Generator{
		"PrimGenerator":   &PrimGenerator{},
		"IslandGenerator": &IslandGenerator{},
	}
	for name, gen := range gens {
		t.Run(name, func(t *testing.T) {
			var seen []*Maze
			for i := 0; i < 5; i++ {
				m, err := gen.Generate(21, 15)
				if err != nil {
					t.Error(err)
					continue
				}
				if m.Rows != 21 || m.Cols != 15 || len(m.String()) != 21*16-1 {
					t.Error("invalid dimensions")
				}
				if m.Start == m.End {
					t.Error("overlapping start and end")
				}
				if m.Wall(m.Start) || m.Wall(m.End) {
					t.Error("wall at invalid location (start or stop)")
				}
				if err != nil {
					t.Fatal(err)
				}
				if Solve(m) == nil {
					t.Errorf("unsolvable: %#v", m.String())
				}
				for _, m1 := range seen {
					if m1.String() == m.String() {
						t.Errorf("duplicate maze: %#v", m.String())
						break
					}
				}
				seen = append(seen, m)
			}
		})
	}
}
