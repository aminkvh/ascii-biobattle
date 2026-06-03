package main

// SpriteFrame holds one frame of multi-line ASCII art.
type SpriteFrame struct {
	Lines []string
}

func (f *SpriteFrame) Width() int {
	w := 0
	for _, l := range f.Lines {
		if len(l) > w {
			w = len(l)
		}
	}
	return w
}

func (f *SpriteFrame) Height() int { return len(f.Lines) }

// ---------------------------------------------------------------------------
// KNIGHT SPRITES  (10 wide × 10 tall)
// Faces right toward the robot side. Includes visored helmet, armor, Cape,
// shield, and a long magic sword.
// ---------------------------------------------------------------------------
var KnightSprites = []SpriteFrame{
	// 0 — Walk A (Sword ready, legs moving)
	{Lines: []string{
		`   _/\_   `,
		`  / \ / \ `,
		`  \  V  / `,
		`  /|===|\ `,
		` / |###| \  /`,
		`[  |###|  ]/ `,
		`  \ ___ / /  `,
		`   /   \ /   `,
		`  /|   |\     `,
		`  \_| |_/     `,
	}},
	// 1 — Walk B (slight stride alteration)
	{Lines: []string{
		`   _/\_   `,
		`  / \ / \ `,
		`  \  V  / `,
		`  /|===|\ `,
		` / |###| \  /`,
		`[  |###|  ]/ `,
		`  \ ___ / /  `,
		`   /   \ /   `,
		`  /_| |_/     `,
		`  \_\ \_\     `,
	}},
	// 2 — Attack (Sword thrust forward)
	{Lines: []string{
		`   _/\_   `,
		`  / \ / \ `,
		`  \  V  / `,
		`  /|===|\=======>`,
		` / |###| \`,
		`[  |###|  ]`,
		`  \ ___ / `,
		`   /   \  `,
		`  /|   |\ `,
		`  \_| |_/ `,
	}},
	// 3 — Defend (Shield raised in front, sword held back)
	{Lines: []string{
		`   _/\_   `,
		`  / \ / \ `,
		`  \  V  / `,
		`  /|===|\ `,
		` /|[###]|\ \`,
		`[ |[###]| ] \`,
		`  \ ___ /    \`,
		`   /   \      \`,
		`  /|   |\      `,
		`  \_| |_/      `,
	}},
	// 4 — Hit (visor skewed, sword dropped back)
	{Lines: []string{
		`   _/\_   `,
		`  / * * \ `,
		`  \  ?  / `,
		`  /|===|\ `,
		` / |###| \ \`,
		`[  |###|  ] \`,
		`  \ ___ /    \`,
		`   /   \     (`,
		`  /|   |\     `,
		`  \_| |_/     `,
	}},
	// 5 — Dead (Fallen knight, sword broken)
	{Lines: []string{
		`          `,
		`          `,
		`          `,
		`          `,
		`          `,
		`   _/\_   `,
		`  / x x \ `,
		`  \  -  / `,
		`  /|===|\==/`,
		`(________)`,
	}},
}

// ---------------------------------------------------------------------------
// ROBOT SPRITES  (11 wide × 11 tall)
// Faces left toward the knight side. Sci-fi mech look with visor lenses,
// shoulder rails, energy core, and left-arm blaster.
// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------
var RobotSprites = []SpriteFrame{
	// 0 — Walk A (Sword ready, legs walking)
	{Lines: []string{
		`  [o---o]  `,
		`  /| X |\  `,
		`  \|===|/  `,
		`  _|_|__   `,
		`//|#####|\ `,
		`/ |#####| ]`,
		`/  \_____/  `,
		`/   /   \   `,
		`   /|   |\  `,
		`   \|   |/  `,
		`   /_] [_/  `,
	}},
	// 1 — Walk B (slight stride alteration)
	{Lines: []string{
		`  [o---o]  `,
		`  /| X |\  `,
		`  \|===|/  `,
		`  _|_|__   `,
		`//|#####|\ `,
		`/ |#####| ]`,
		`/  \_____/  `,
		`/   /   \   `,
		`   /|   |\  `,
		`  /_/   \_\ `,
		`  \_]   [_/ `,
	}},
	// 2 — Attack (Lightsaber thrust left)
	{Lines: []string{
		`  [o---o]  `,
		`  /| X |\  `,
		`  \|===|/  `,
		`  _|_|__   `,
		`<=======#####|\ `,
		`[ |#####| ]`,
		`  \_____/  `,
		`   /   \   `,
		`  /|   |\  `,
		`  \|   |/  `,
		`  /_] [_/  `,
	}},
	// 3 — Defend (Shield plate covers active, sword held back)
	{Lines: []string{
		`  [o---o]  `,
		`  /| X |\  `,
		`  \|===|/  `,
		`  _|_|__   `,
		`//|#####|\ `,
		`/|[|[#####]|]`,
		`/  \_____/  `,
		`/   /   \   `,
		`   /|   |\  `,
		`   \|   |/  `,
		`   /_] [_/  `,
	}},
	// 4 — Hit (chassis sparks, visor offline)
	{Lines: []string{
		`  [o-x-o]  `,
		`  /| * |\  `,
		`  \|===|/  `,
		`  _|_|__   `,
		` /|#*###|\ `,
		`[ |#####| ]`,
		`  \_____/  `,
		`   /   \   `,
		`  /|   |\  `,
		`  \|   |/  `,
		`  /_] [_/  `,
	}},
	// 5 — Dead (Collapsed scrap pile)
	{Lines: []string{
		`           `,
		`           `,
		`           `,
		`           `,
		`           `,
		`           `,
		`  [o-x-o]  `,
		`  /| * |\  `,
		`  \__*__/  `,
		`===========`,
		` \_______/ `,
	}},
}
