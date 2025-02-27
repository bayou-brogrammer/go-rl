package constants

import (
	"image/color"

	"codeberg.org/anaseto/gruid"
)

const Opaque = uint8(255)

// Terminal Color Schemes
// 0: Black
// 1: Red
// 2: Green
// 3: Yellow
// 4: Blue
// 5: Magenta
// 6: Cyan
// 7: White
// 8: Bright Black
// 9: Bright Red
// 10: Bright Green
// 11: Bright Yellow
// 12: Bright Blue
// 13: Bright Magenta
// 14: Bright Cyan
// 15: Bright White

//////////////////////
// GIMP Palette
// Name: db16.gpl
// Columns: 8
//////////////////////

// Colors slightly changed and named by 'http://mkweb.bcgsc.ca/colornames/namethatcolor/?rgb={red},{green},{blue}' (Decimal Values)
// #17111D     ; very_dark_violet
// #4e4a4e     ; shadowy_lavender
// #716E61     ; flint
// #86949F     ; regent_grey
// #D7E7D0     ; peppermint
// #462428     ; red_earth
// #814D30     ; root_beer
// #D3494E     ; faded_red
// #CD7F32     ; bronze
// #D4A798     ; birthday_suit
// #E3CF57     ; banana
// #333366     ; deep_koamaru
// #5D76CB     ; indigo
// #7AC5CD     ; cadet_blue
// #215E21     ; hunter_green
// #71AA34     ; leaf

// (20,12,28) -> #140C1C Dark Purple
// (68,36,52) -> #442434 Burgundy
// (48,52,109) -> #30346D Navy Blue
// (78,74,78) -> #4E4A4E Dark Gray
// (133,76,48) -> #854C30 Brown
// (52,101,36) -> #346524 Forest Green
// (208,70,72) -> #D04648 Red
// (117,113,97) -> #757161 Gray
// (89,125,206) -> #597DCE Blue
// (210,125,44) -> #D27D2C Orange
// (133,149,161) -> #8595A1 Light Gray
// (109,170,44) -> #6DAA2C Lime Green
// (210,170,153) -> #D2AA99 Peach
// (109,194,202) -> #6DC2CA Light Blue
// (218,212,94) -> #DAD45E Yellow
// (222,238,214) -> #DEEED6 Off White

const (
	DB16_PaletteColorBlack       gruid.Color = 0
	DB16_PaletteColorDarkPurple  gruid.Color = 1
	DB16_PaletteColorBurgundy    gruid.Color = 2
	DB16_PaletteColorNavyBlue    gruid.Color = 3
	DB16_PaletteColorDarkGray    gruid.Color = 4
	DB16_PaletteColorBrown       gruid.Color = 5
	DB16_PaletteColorForestGreen gruid.Color = 6
	DB16_PaletteColorRed         gruid.Color = 7
	DB16_PaletteColorGray        gruid.Color = 8
	DB16_PaletteColorBlue        gruid.Color = 9
	DB16_PaletteColorOrange      gruid.Color = 10
	DB16_PaletteColorLightGray   gruid.Color = 11
	DB16_PaletteColorLimeGreen   gruid.Color = 12
	DB16_PaletteColorPeach       gruid.Color = 13
	DB16_PaletteColorLightBlue   gruid.Color = 14
	DB16_PaletteColorYellow      gruid.Color = 15
	DB16_PaletteColorOffWhite    gruid.Color = 16
)

/////////////////////////////////////////////////
// Solarized Colors
/////////////////////////////////////////////////

// $base03:    #002b36;
// $base02:    #073642;
// $base01:    #586e75;
// $base00:    #657b83;
// $base0:     #839496;
// $base1:     #93a1a1;
// $base2:     #eee8d5;
// $base3:     #fdf6e3;
// $yellow:    #b58900;
// $orange:    #cb4b16;
// $red:       #dc322f;
// $magenta:   #d33682;
// $violet:    #6c71c4;
// $blue:      #268bd2;
// $cyan:      #2aa198;
// $green:     #859900;

const (
	Solairzed_Base03  gruid.Color = 0  // #002b36
	Solairzed_Base02  gruid.Color = 1  // #073642
	Solairzed_Base01  gruid.Color = 2  // #586e75
	Solairzed_Base00  gruid.Color = 3  // #657b83
	Solairzed_Base0   gruid.Color = 4  // #839496
	Solairzed_Base1   gruid.Color = 5  // #93a1a1
	Solairzed_Base2   gruid.Color = 6  // #eee8d5
	Solairzed_Base3   gruid.Color = 7  // #fdf6e3
	Solairzed_Yellow  gruid.Color = 8  // #b58900
	Solairzed_Orange  gruid.Color = 9  // #cb4b16
	Solairzed_Red     gruid.Color = 10 // #dc322f
	Solairzed_Magenta gruid.Color = 11 // #d33682
	Solairzed_Violet  gruid.Color = 12 // #6c71c4
	Solairzed_Blue    gruid.Color = 13 // #268bd2
)

var SolarizedColorMap = map[gruid.Color]color.RGBA{
	Solairzed_Base03:  {0, 43, 54, Opaque},     //"#002b36",
	Solairzed_Base02:  {7, 54, 66, Opaque},     //"#073642",
	Solairzed_Base01:  {88, 110, 117, Opaque},  //"#586e75",
	Solairzed_Base00:  {101, 123, 131, Opaque}, //"#657b83",
	Solairzed_Base0:   {131, 148, 150, Opaque}, //"#839496",
	Solairzed_Base1:   {147, 161, 161, Opaque}, //"#93a1a1",
	Solairzed_Base2:   {238, 232, 213, Opaque}, //"#eee8d5",
	Solairzed_Base3:   {253, 246, 227, Opaque}, //"#fdf6e3",
	Solairzed_Yellow:  {181, 137, 0, Opaque},   //"#b58900",
	Solairzed_Orange:  {203, 75, 22, Opaque},   //"#cb4b16",
	Solairzed_Red:     {220, 50, 47, Opaque},   //"#dc322f",
	Solairzed_Magenta: {211, 54, 130, Opaque},  //"#d33682",
	Solairzed_Violet:  {108, 113, 196, Opaque}, //"#6c71c4",
	Solairzed_Blue:    {38, 139, 210, Opaque},  //"#268bd2",
}
