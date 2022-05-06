package main

import "testing"

func TestCollison(t *testing.T) {
	simpleGame := Game{}
	simpleGame.player = Sprite{yloc: 100, xloc: 100, pict: loadPNGImageFromEmbedded("blueBox.png")}
	simpleGame.enemy = Sprite{yloc: 100, xloc: 100, pict: loadPNGImageFromEmbedded("blueBox.png")}
	tests := []struct {
		sprite1  Sprite
		sprite2  Sprite
		game     Game
		expected bool
	}{{
		sprite1:  simpleGame.player,
		sprite2:  simpleGame.enemy,
		game:     simpleGame,
		expected: true,
	},
		{
			sprite1:  Sprite{xloc: 1, yloc: 1, pict: loadPNGImageFromEmbedded("blueBox.png")},
			sprite2:  Sprite{xloc: 100, yloc: 100, pict: loadPNGImageFromEmbedded("blueBox.png")},
			game:     simpleGame,
			expected: false,
		},
	}

	for _, test := range tests {
		actual := appleCollison(test.sprite1, test.sprite2, &test.game)
		if actual != test.expected {
			t.Errorf("The tests did not meet the expected result")
		}
	}
}
func TestSegListLength(t *testing.T) {
	simpleGame := Game{}
	simpleGame.player = Sprite{yloc: 100, xloc: 100, pict: loadPNGImageFromEmbedded("blueBox.png")}
	simpleGame.enemyList = make([]Sprite, 1001)
	tests := []struct {
		amount   int
		game     Game
		expected int
	}{{
		amount: 1000, game: simpleGame, expected: 1001,
	}, {
		amount: 100, game: simpleGame, expected: 101,
	}}
	for _, test := range tests {
		actual := fillList(&test.game, test.amount)
		if actual != test.expected {
			t.Errorf("The tests did not meet the expected result")
		}
	}
}
func TestYcheck(t *testing.T) {
	simpleGame := Game{}
	simpleGame.player = Sprite{yloc: 100, xloc: 100, pict: loadPNGImageFromEmbedded("blueBox.png")}
	simpleGame.enemy = Sprite{yloc: 100, xloc: 100, pict: loadPNGImageFromEmbedded("blueBox.png")}
	sprite1 := Sprite{xloc: 1, yloc: 1, pict: loadPNGImageFromEmbedded("blueBox.png")}
	sprite2 := Sprite{xloc: 100, yloc: 100, pict: loadPNGImageFromEmbedded("blueBox.png")}
	tests := []struct {
		y1       int
		y2       int
		game     Game
		expected int
	}{{
		y1:       simpleGame.player.yloc,
		y2:       simpleGame.enemy.yloc,
		game:     simpleGame,
		expected: 0,
	},
		{
			y1:       sprite1.yloc,
			y2:       sprite2.yloc,
			game:     simpleGame,
			expected: 99,
		},
	}

	for _, test := range tests {
		actual := yCheck(test.y1, test.y2)
		if actual != test.expected {
			t.Errorf("The tests did not meet the expected result")
		}
	}
}
