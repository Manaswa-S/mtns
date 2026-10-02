package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

// {sky, mountain}
var palettes = [][2]color.RGBA{
	// --- Dusk / sunset classics ---
	{{250, 180, 140, 255}, {102, 44, 230, 255}}, // peach sky, violet (matches your image)
	{{255, 190, 120, 255}, {120, 40, 160, 255}}, // apricot, plum
	{{255, 154, 118, 255}, {80, 60, 200, 255}},  // coral, indigo
	{{255, 204, 153, 255}, {200, 60, 120, 255}}, // light orange, raspberry
	{{255, 170, 150, 255}, {140, 50, 190, 255}}, // salmon, amethyst
	{{253, 215, 150, 255}, {220, 90, 70, 255}},  // golden, terracotta
	{{255, 140, 100, 255}, {60, 40, 150, 255}},  // burnt orange, deep indigo
	{{255, 200, 170, 255}, {170, 70, 150, 255}}, // blush, orchid
	{{250, 160, 110, 255}, {40, 90, 180, 255}},  // sunset orange, cobalt
	{{255, 225, 160, 255}, {190, 80, 110, 255}}, // pale gold, rose

	// --- Blues & cool tones ---
	{{200, 230, 255, 255}, {40, 100, 220, 255}}, // pale sky, royal blue
	{{180, 220, 250, 255}, {30, 130, 200, 255}}, // morning blue, azure
	{{220, 240, 255, 255}, {60, 80, 190, 255}},  // ice, periwinkle
	{{170, 210, 240, 255}, {20, 110, 150, 255}}, // powder, teal blue
	{{210, 225, 245, 255}, {70, 70, 170, 255}},  // misty, slate blue
	{{190, 235, 245, 255}, {0, 150, 180, 255}},  // aqua mist, cyan
	{{230, 240, 250, 255}, {25, 60, 140, 255}},  // near white, navy
	{{160, 200, 235, 255}, {90, 60, 200, 255}},  // cool sky, blue violet
	{{200, 215, 240, 255}, {50, 120, 170, 255}}, // haze, steel blue
	{{175, 225, 235, 255}, {20, 80, 160, 255}},  // seafoam sky, deep blue

	// --- Greens & forest ---
	{{220, 245, 200, 255}, {40, 150, 80, 255}},  // light lime sky, forest green
	{{240, 240, 190, 255}, {60, 140, 60, 255}},  // pale yellow, leaf green
	{{200, 235, 220, 255}, {20, 130, 110, 255}}, // mint, jade
	{{235, 245, 210, 255}, {90, 160, 40, 255}},  // cream green, grass
	{{190, 230, 210, 255}, {30, 110, 90, 255}},  // soft mint, pine
	{{250, 235, 180, 255}, {50, 160, 120, 255}}, // butter, emerald
	{{210, 240, 230, 255}, {10, 120, 100, 255}}, // pale aqua, deep teal
	{{245, 245, 200, 255}, {120, 170, 50, 255}}, // lemon cream, olive green

	// --- Reds, pinks, warm ---
	{{255, 220, 200, 255}, {200, 50, 60, 255}},  // peach cream, crimson
	{{255, 200, 210, 255}, {190, 40, 100, 255}}, // pink, magenta red
	{{255, 235, 200, 255}, {210, 80, 40, 255}},  // cream, burnt sienna
	{{255, 190, 190, 255}, {160, 40, 80, 255}},  // light rose, wine
	{{250, 210, 170, 255}, {230, 100, 50, 255}}, // sand, orange
	{{255, 225, 225, 255}, {220, 70, 120, 255}}, // blush white, hot pink
	{{255, 210, 160, 255}, {180, 40, 40, 255}},  // amber, brick red
	{{255, 180, 170, 255}, {230, 60, 90, 255}},  // salmon pink, cherry

	// --- Desert / earth ---
	{{245, 222, 179, 255}, {190, 110, 50, 255}},  // wheat, copper
	{{240, 210, 160, 255}, {160, 90, 60, 255}},   // sand, clay
	{{255, 228, 181, 255}, {205, 133, 63, 255}},  // moccasin, peru
	{{235, 205, 170, 255}, {140, 90, 70, 255}},   // dusty tan, umber
	{{250, 225, 190, 255}, {215, 120, 60, 255}},  // desert light, rust
	{{230, 200, 160, 255}, {120, 100, 130, 255}}, // dust, dusty mauve

	// --- Night / dark skies (mountains still bright enough to gradient) ---
	{{20, 24, 60, 255}, {90, 80, 200, 255}},   // midnight blue, lavender blue
	{{30, 20, 50, 255}, {150, 70, 190, 255}},  // dark violet, purple
	{{15, 30, 50, 255}, {40, 140, 160, 255}},  // deep navy teal, teal
	{{40, 20, 40, 255}, {200, 70, 110, 255}},  // dark plum, rose
	{{10, 20, 40, 255}, {70, 120, 220, 255}},  // night, bright blue
	{{25, 25, 35, 255}, {120, 130, 160, 255}}, // charcoal, steel gray-blue
	{{35, 15, 30, 255}, {220, 100, 60, 255}},  // wine dark, ember orange
	{{10, 35, 40, 255}, {60, 180, 140, 255}},  // dark teal, aurora green

	// --- Fantasy / stylized ---
	{{255, 240, 120, 255}, {255, 90, 150, 255}},  // lemon yellow, bubblegum
	{{255, 120, 200, 255}, {70, 50, 200, 255}},   // synthwave pink, electric indigo
	{{140, 240, 230, 255}, {200, 60, 200, 255}},  // cyan sky, magenta
	{{255, 250, 180, 255}, {80, 200, 170, 255}},  // pale yellow, turquoise
	{{255, 160, 80, 255}, {180, 30, 140, 255}},   // vivid orange, deep magenta
	{{200, 180, 255, 255}, {255, 110, 90, 255}},  // lilac sky, coral mountains
	{{255, 245, 220, 255}, {90, 90, 110, 255}},   // ivory, slate
	{{180, 255, 220, 255}, {240, 100, 130, 255}}, // mint sky, watermelon

	// --- Monochrome-ish & subtle ---
	{{235, 235, 240, 255}, {100, 100, 130, 255}}, // fog, gray violet
	{{225, 230, 235, 255}, {70, 110, 130, 255}},  // overcast, blue gray
	{{245, 240, 235, 255}, {140, 100, 100, 255}}, // warm white, dusty rose
	{{220, 225, 215, 255}, {90, 130, 110, 255}},  // pale sage, sage green

	// --- Golden hour ---
	{{255, 214, 140, 255}, {196, 74, 62, 255}},  // honey, brick
	{{255, 196, 128, 255}, {150, 60, 110, 255}}, // tangerine cream, mulberry
	{{252, 211, 165, 255}, {89, 72, 176, 255}},  // apricot, soft indigo
	{{255, 224, 178, 255}, {222, 110, 75, 255}}, // vanilla, burnt coral
	{{248, 187, 122, 255}, {115, 55, 140, 255}}, // amber, grape
	{{255, 207, 158, 255}, {45, 106, 160, 255}}, // melon, denim
	{{254, 220, 120, 255}, {212, 84, 96, 255}},  // sunflower, watermelon rose
	{{245, 190, 150, 255}, {70, 90, 170, 255}},  // clay pink, cornflower

	// --- Pastel dreams ---
	{{255, 214, 224, 255}, {140, 110, 220, 255}}, // cotton candy, wisteria
	{{214, 230, 255, 255}, {230, 120, 160, 255}}, // baby blue, flamingo
	{{255, 236, 214, 255}, {120, 150, 220, 255}}, // cream, periwinkle
	{{225, 214, 255, 255}, {90, 160, 200, 255}},  // lavender mist, sky teal
	{{214, 255, 236, 255}, {150, 110, 200, 255}}, // mint, soft purple
	{{255, 228, 240, 255}, {100, 120, 190, 255}}, // pale pink, dusk blue
	{{240, 230, 255, 255}, {200, 100, 140, 255}}, // pearl lilac, dusty rose
	{{255, 245, 214, 255}, {110, 170, 190, 255}}, // butter, glacier teal

	// --- Ocean & coastal ---
	{{204, 240, 250, 255}, {16, 100, 140, 255}}, // lagoon mist, deep sea
	{{225, 245, 245, 255}, {40, 130, 150, 255}}, // sea foam, peacock
	{{255, 238, 210, 255}, {20, 120, 170, 255}}, // sand cream, ocean blue
	{{190, 225, 240, 255}, {30, 70, 120, 255}},  // harbor sky, marine
	{{245, 250, 255, 255}, {60, 150, 190, 255}}, // snow white, cerulean
	{{255, 230, 200, 255}, {30, 150, 160, 255}}, // shell, turquoise
	{{200, 235, 230, 255}, {80, 110, 180, 255}}, // tide pool, bluebell

	// --- Forest & meadow ---
	{{255, 247, 205, 255}, {70, 140, 90, 255}},  // morning light, fern
	{{215, 235, 190, 255}, {30, 120, 80, 255}},  // pale moss, evergreen
	{{250, 240, 200, 255}, {40, 110, 110, 255}}, // parchment, spruce teal
	{{230, 245, 220, 255}, {110, 150, 70, 255}}, // dew, olive
	{{255, 232, 190, 255}, {85, 135, 75, 255}},  // warm cream, meadow
	{{200, 225, 205, 255}, {50, 100, 70, 255}},  // lichen, pine
	{{245, 235, 170, 255}, {30, 140, 120, 255}}, // wheat gold, viridian

	// --- Autumn ---
	{{250, 214, 170, 255}, {170, 70, 40, 255}},  // pumpkin cream, rust
	{{240, 200, 140, 255}, {135, 55, 55, 255}},  // harvest, cranberry
	{{255, 225, 170, 255}, {195, 105, 40, 255}}, // maple glow, amber brown
	{{235, 205, 180, 255}, {120, 70, 90, 255}},  // fawn, plum brown
	{{248, 215, 150, 255}, {160, 85, 35, 255}},  // straw, cinnamon
	{{225, 190, 160, 255}, {110, 85, 120, 255}}, // dusty peach, heather

	// --- Twilight & night ---
	{{25, 30, 70, 255}, {200, 100, 170, 255}}, // midnight, orchid pink
	{{15, 15, 45, 255}, {80, 140, 220, 255}},  // ink blue, moonlit blue
	{{40, 25, 70, 255}, {240, 130, 100, 255}}, // royal night, ember coral
	{{20, 40, 60, 255}, {120, 200, 190, 255}}, // deep ocean sky, aqua glow
	{{30, 10, 40, 255}, {230, 90, 140, 255}},  // blackberry, neon rose
	{{12, 24, 38, 255}, {170, 190, 230, 255}}, // abyss, moon silver blue
	{{45, 25, 55, 255}, {255, 160, 90, 255}},  // eggplant, amber
	{{18, 18, 40, 255}, {110, 220, 160, 255}}, // starless navy, aurora mint

	// --- Synthwave & vivid ---
	{{255, 110, 180, 255}, {60, 40, 160, 255}},   // hot pink, ultramarine
	{{255, 170, 60, 255}, {200, 40, 120, 255}},   // marigold, fuchsia
	{{120, 230, 255, 255}, {170, 70, 220, 255}},  // electric cyan, violet
	{{255, 90, 120, 255}, {40, 60, 150, 255}},    // watermelon, sapphire
	{{255, 200, 80, 255}, {230, 70, 70, 255}},    // sunburst, poppy
	{{180, 130, 255, 255}, {255, 130, 110, 255}}, // amethyst sky, peach mountains
	{{255, 130, 90, 255}, {90, 40, 170, 255}},    // sunset orange, royal purple
	{{100, 255, 200, 255}, {220, 70, 150, 255}},  // neon mint, raspberry

	// --- Moody & muted ---
	{{210, 205, 215, 255}, {95, 85, 140, 255}},  // storm cloud, muted indigo
	{{225, 215, 205, 255}, {120, 95, 110, 255}}, // stone, mauve brown
	{{200, 210, 220, 255}, {65, 95, 125, 255}},  // slate haze, steel blue
	{{235, 225, 215, 255}, {150, 100, 80, 255}}, // linen, cocoa clay
	{{215, 220, 205, 255}, {85, 115, 105, 255}}, // gray sage, eucalyptus
	{{230, 210, 215, 255}, {110, 90, 130, 255}}, // dusty pink, smoky purple

	// --- Arctic & icy ---
	{{235, 248, 255, 255}, {110, 140, 220, 255}}, // frost, glacier blue
	{{210, 235, 255, 255}, {150, 120, 210, 255}}, // ice blue, pale violet
	{{240, 245, 250, 255}, {80, 120, 160, 255}},  // snowfield, cold steel
	{{255, 240, 245, 255}, {120, 130, 200, 255}}, // rose-tinted snow, periwinkle
}

type Input struct {
	// the total number of ranges to be generated
	NoOfRanges int
	// color of the ranges
	RangeColor color.RGBA
	// color of the sky
	SkyColor color.RGBA
	// width of the image
	Width int
	// height of the image
	Height int
	// output path of the generated image
	ResultPath string
	FileName   string

	// randomly derived heights of the ranges
	heights []int
	// seed for random number generation
	Seed int64
}

type DebugInfo struct {
	// flag to enable debug output
	DebugFlag bool
	// flag to save the generated image to an archive folder
	ArchiveSaveFlag bool
	// path to the archive folder
	ArchivePath string
}

type RangePoint struct {
	// x-coordinate of the range point
	X float64
	// y-coordinate of the range point
	Y float64
}

type TreePoint struct {
	// x-coordinate of the tree point
	X float64
	// y-coordinate of the tree point
	Y float64
	// height of the tree
	Height float64
}

var debugInfo = DebugInfo{}

var rng *rand.Rand

func main() {
	fmt.Println("Welcome to the Basque Mountain Generator!")

	inp, err := GetInput()
	if err != nil {
		panic(fmt.Sprintln("Error:", err))
	}

	err = GenerateBasqueMountain(inp)
	if err != nil {
		panic(fmt.Sprintln("Error:", err))
	}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

func GetInput() (*Input, error) {

	input := new(Input)

	ranges := flag.Int("ranges", 6, "number of ranges to be generated")
	colorStr := flag.String("color", "random", "color of the ranges&sky in RGBA&RGBA format (e.g. 255,0,0,255&255,0,0,255 for red&red for ranges&sky)")
	width := flag.Int("width", 1920, "width of the image, max 7680, min 2")
	height := flag.Int("height", 1080, "height of the image, max 4320, min 2")
	outPath := flag.String("out", "/home/mnswa/Pictures/Wallpapers", "output path of the generated image")
	fileName := flag.String("file", "basque_mountains", "name of the output file")
	seed := flag.Int64("seed", time.Now().UnixNano(), "seed for random number generation (default: current time in nanoseconds)")

	debugFlag := flag.Bool("debug", false, "enable debug output")
	archiveSaveFlag := flag.Bool("archive", true, "save the generated image to an archive folder")
	archivePath := flag.String("archivePath", "/home/mnswa/Pictures/Wallpapers/.archive/", "path to the archive folder")

	flag.Parse()

	fmt.Println("SEED: ", *seed)

	// IMP: seed the random number generator
	rng = rand.New(rand.NewSource(*seed))

	var palette [2]color.RGBA
	{

		if *colorStr == "random" {
			palette = palettes[getRandIntNInc(rng, 0, len(palettes)-1)]
		} else {
			colors := strings.Split(*colorStr, "&")
			if len(colors) != 2 {
				return nil, fmt.Errorf("Invalid color format. Please provide colors in RGBA format (e.g. 255,0,0,255&255,0,0,255 for red&red for range&sky)")
			}
			colors1 := strings.Split(colors[0], ",")
			colors2 := strings.Split(colors[1], ",")
			if len(colors1) != 4 || len(colors2) != 4 {
				return nil, fmt.Errorf("Invalid color format. Please provide color in RGBA format (e.g. 255,0,0,255 for red)")

			}

			// Parse the color values
			r1, err1 := strconv.Atoi(colors1[0])
			g1, err2 := strconv.Atoi(colors1[1])
			b1, err3 := strconv.Atoi(colors1[2])
			a1, err4 := strconv.Atoi(colors1[3])

			if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
				return nil, fmt.Errorf("Invalid color values. Please provide valid integer values for RGBA.")
			}

			r2, err1 := strconv.Atoi(colors2[0])
			g2, err2 := strconv.Atoi(colors2[1])
			b2, err3 := strconv.Atoi(colors2[2])
			a2, err4 := strconv.Atoi(colors2[3])

			if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
				return nil, fmt.Errorf("Invalid color values. Please provide valid integer values for RGBA.")
			}

			palette[0] = color.RGBA{uint8(r1), uint8(g1), uint8(b1), uint8(a1)}
			palette[1] = color.RGBA{uint8(r2), uint8(g2), uint8(b2), uint8(a2)}
		}

	}

	{
		if (*width > 7680 || *width < 2) || (*height > 4320 || *height < 2) {
			return nil, fmt.Errorf("Invalid width or height. Please provide valid values for width and height.")
		}
	}

	{
		if *ranges < 1 || *ranges > 16 {
			return nil, fmt.Errorf("Invalid number of ranges. Please provide a value between 1 and 16.")
		}
	}

	{
		if err := os.MkdirAll(*outPath, 0755); err != nil {
			return nil, fmt.Errorf("Failed to create output directory.")
		}

		if *fileName == "" {
			return nil, fmt.Errorf("Invalid file name. Please provide a valid file name.")
		}
	}

	input.NoOfRanges = *ranges
	input.RangeColor = palette[1]
	input.SkyColor = palette[0]
	input.Width = *width
	input.Height = *height
	input.ResultPath = *outPath
	input.FileName = *fileName
	input.Seed = *seed

	debugInfo.DebugFlag = *debugFlag
	debugInfo.ArchiveSaveFlag = *archiveSaveFlag
	debugInfo.ArchivePath = *archivePath

	debug("INPUT:", input)

	return input, nil
}

// TODO: fade local is wrong, should also consider if the colors are white/black, also the fade should be gradual not linear.

func GenerateBasqueMountain(inp *Input) error {
	// get the random heights for each range
	inp.heights = genHeights(inp.Height, inp.NoOfRanges)
	debug("HEIGHTS:", inp.heights)

	// generate the range lines
	lines := make([][]RangePoint, 0, inp.NoOfRanges)
	for i, h := range inp.heights {
		octaves := 6
		lacunarity := 2.5
		persistence := float64(i)*0.03 + 0.30
		if persistence > 0.55 {
			persistence = 0.55
		}
		lines = append(lines,
			genMountainLine(RangePoint{X: 0, Y: float64(h)}, RangePoint{X: float64(inp.Width), Y: float64(h)}, inp.Width, inp.Height, (inp.Seed+int64(i+1)), octaves, lacunarity, persistence))
	}

	// get mountain colors
	white := getRandIntNInc(rng, 0, 1)
	var mountainColors []color.RGBA
	if white == 1 {
		mountainColors = fadeToWhite(inp.RangeColor, inp.NoOfRanges, 0.75)
	} else {
		mountainColors = fadeToBlack(inp.RangeColor, inp.NoOfRanges, 0.9)
	}
	debug("MOUNTAIN COLORS:", mountainColors)

	// get random trees
	trees := genTrees(lines, -1)
	debug("TREES:", trees)

	// Create image.
	img := image.NewRGBA(image.Rect(0, 0, inp.Width, inp.Height))

	// paint background with sky color
	for y := 0; y < inp.Height; y++ {
		for x := 0; x < inp.Width; x++ {
			img.Set(x, y, inp.SkyColor)
		}
	}

	// draw the mountains
	for i := inp.NoOfRanges - 1; i >= 0; i-- {
		drawMountain(img, lines[i], mountainColors[i], i != inp.NoOfRanges-1)
	}

	// draw the trees
	for _, tree := range trees {
		drawTree(img, tree, rng)
	}

	// Save PNG.
	ext := ".png"
	finalName := fmt.Sprintf("%s%s", inp.FileName, ext)
	finalPath := filepath.Join(inp.ResultPath, finalName)
	file, err := os.Create(finalPath)
	if err != nil {
		return fmt.Errorf("failed to create the result file: %v", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode the result file: %v", err)
	}

	fmt.Println("Generated:", finalPath)
	debug("GENERATED:", finalPath)

	if !debugInfo.ArchiveSaveFlag {
		return nil
	}

	// save to the archive
	arcDir := filepath.Dir(debugInfo.ArchivePath)
	archiveName := fmt.Sprintf("%s_%d_%d%s", inp.FileName, time.Now().UnixMilli(), inp.Seed, ext)
	err = os.MkdirAll(arcDir, 0755)
	if err != nil {
		panic(err)
	}
	arcfile, err := os.Create(filepath.Join(arcDir, archiveName))
	if err != nil {
		panic(err)
	}
	defer arcfile.Close()

	if err := png.Encode(arcfile, img); err != nil {
		panic(err)
	}

	fmt.Println("Generated:", filepath.Join(arcDir, archiveName))

	return nil
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

func debug(msg string, args ...interface{}) {
	if debugInfo.DebugFlag {
		fmt.Print("DEBUG: ", msg)
		fmt.Println(args...)
	}
}

func getRandIntNInc(rng *rand.Rand, min, max int) int {
	return rng.Intn(max-min+1) + min
}

func genHeights(height int, noOfRanges int) []int {
	padding := (height / (10 - getRandIntNInc(rng, 1, 5)))
	heights := make([]int, noOfRanges)
	for i := range noOfRanges {
		heights[i] = getRandIntNInc(rng, padding, height-padding)
	}
	slices.Sort(heights)
	return heights
}

func genTrees(lines [][]RangePoint, count int) []TreePoint {
	// randomly generate trees

	if count < 0 {
		count = getRandIntNInc(rng, 1, 15)
	}

	slices.Reverse(lines)
	for i := 0; i < len(lines[0]); i++ {
		for j := 1; j < len(lines); j++ {
			if lines[j-1][i].Y < lines[j][i].Y {
				lines[j][i].Y = lines[j-1][i].Y
			}
		}
	}

	var trees []TreePoint

	for range count {
		// Select a random line
		line := lines[getRandIntNInc(rng, 0, len(lines)-1)]

		// Select a random point on the line
		point := line[getRandIntNInc(rng, 0, len(line)-1)]

		trees = append(trees, TreePoint{
			X:      point.X,
			Y:      point.Y,
			Height: float64(50 + getRandIntNInc(rng, 0, 50)),
		})
	}

	return trees
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// ------------------------------------------------------------
// Deterministic 1D value noise
// Returns roughly [-1, 1].
// ------------------------------------------------------------

func hash(x int, seed int64) float64 {
	n := int64(x)*374761393 + seed*668265263
	n = (n ^ (n >> 13)) * 1274126177
	n ^= n >> 16

	// Convert to [0, 1]
	return float64(uint64(n)&0x7fffffff) / float64(0x7fffffff)
}

func smoothstep(t float64) float64 {
	return t * t * (3 - 2*t)
}

func noise1D(x float64, seed int64) float64 {
	x0 := int(math.Floor(x))
	x1 := x0 + 1

	t := x - float64(x0)
	t = smoothstep(t)

	v0 := hash(x0, seed)
	v1 := hash(x1, seed)

	return (v0+(v1-v0)*t)*2 - 1
}

// ------------------------------------------------------------
// Fractional Brownian Motion
// ------------------------------------------------------------

func fbm(x float64, seed int64, octaves int, lacunarity, persistence float64) float64 {
	value := 0.0
	amplitude := 1.0
	frequency := 1.0
	maxAmplitude := 0.0

	for i := 0; i < octaves; i++ {
		value += noise1D(x*frequency, seed+int64(i)*1000) * amplitude

		maxAmplitude += amplitude
		amplitude *= persistence
		frequency *= lacunarity
	}

	// Normalize approximately to [-1, 1]
	return value / maxAmplitude
}

// ------------------------------------------------------------
// Generate mountain range line
//
// start  - beginning of range
// end    - ending of range
// width  - number of points
// seed   - controls the generated shape
// ------------------------------------------------------------

func genMountainLine(start, end RangePoint, width, height int, seed int64, octaves int, lacunarity, persistence float64) []RangePoint {
	if width < 2 {
		return []RangePoint{start}
	}

	points := make([]RangePoint, width)

	dx := end.X - start.X
	dy := end.Y - start.Y

	length := math.Hypot(dx, dy)

	// Unit perpendicular vector.
	// This is the direction in which we move the
	// straight line to create the mountain shape.
	nx := -dy / length
	ny := dx / length

	// Controls how large the mountain line can deviate
	// from the original start -> end line.
	amplitude := length * 0.1
	// if amplitude > start.Y || amplitude > (float64(height)-start.Y) {
	// 	amplitude = 0.80 * math.Min(start.Y, length-start.Y)
	// }

	for i := range width {
		t := float64(i) / float64(width-1)

		// Original point on the straight line.
		baseX := start.X + dx*t
		baseY := start.Y + dy*t

		// fBM gives us smooth multi-scale randomness.
		noise := fbm(t*4.0, seed, octaves, lacunarity, persistence)

		// Force displacement to zero at both endpoints.
		envelope := math.Sin(math.Pi * t)

		offset := noise * amplitude * envelope

		points[i] = RangePoint{
			X: baseX + nx*offset,
			Y: baseY + ny*offset,
		}
	}

	return points
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

func fadeToBlack(base color.RGBA, blocks int, strength float64) []color.RGBA {
	colors := make([]color.RGBA, blocks)

	for i := 0; i < blocks; i++ {
		t := float64(i) / float64(blocks-1)
		t *= strength // 0.9
		colors[i] = color.RGBA{
			R: uint8(float64(base.R) * t),
			G: uint8(float64(base.G) * t),
			B: uint8(float64(base.B) * t),
			A: uint8(255),
		}
	}

	return colors
}

func fadeToWhite(base color.RGBA, blocks int, strength float64) []color.RGBA {
	colors := make([]color.RGBA, blocks)

	for i := 0; i < blocks; i++ {
		t := float64(i) / float64(blocks-1)
		t *= strength // 0.75

		colors[i] = color.RGBA{
			R: uint8(float64(base.R) + (255-float64(base.R))*t),
			G: uint8(float64(base.G) + (255-float64(base.G))*t),
			B: uint8(float64(base.B) + (255-float64(base.B))*t),
			A: uint8(255),
		}
	}

	return colors
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

func drawMountain(img *image.RGBA, points []RangePoint, colorRGBA color.RGBA, fadeLocal bool) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// For every segment between two points,
	// interpolate the Y value for each X.
	for i := 0; i < len(points)-1; i++ {
		p1 := points[i]
		p2 := points[i+1]

		dx := p2.X - p1.X
		if dx == 0 {
			continue
		}

		for x := p1.X; x <= p2.X && x < float64(width); x++ {
			t := float64(x-p1.X) / float64(dx)

			y := int(float64(p1.Y) +
				t*float64(p2.Y-p1.Y))

			img.Set(int(x), y, colorRGBA)

			// Fill everything below the mountain line.
			if fadeLocal {
				localCols := fadeToWhite(colorRGBA, height-y+1, 0.15)
				for py := y; py < height; py++ {
					img.Set(int(x), py, localCols[py-y])
				}
			} else {
				for py := y; py < height; py++ {
					img.Set(int(x), py, colorRGBA)
				}
			}

		}
	}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

func drawTree(img *image.RGBA, tree TreePoint, rng *rand.Rand) {
	x := tree.X
	baseY := tree.Y
	h := tree.Height

	// ---------------------------------------------------------
	// Colors
	// ---------------------------------------------------------

	trunkDark := color.RGBA{0, 0, 0, 255}

	greens := []color.RGBA{
		{0, 0, 0, 255},
	}

	// ---------------------------------------------------------
	// Trunk
	// ---------------------------------------------------------

	trunkWidth := h * 0.045

	drawLine(
		img,
		RangePoint{x - trunkWidth, baseY},
		RangePoint{x + trunkWidth, baseY - h*0.92},
		trunkWidth*1.5,
		trunkDark,
	)

	// ---------------------------------------------------------
	// Branches
	// ---------------------------------------------------------

	branchCount := int(12 + h/25)

	for i := 0; i < branchCount; i++ {

		// 0 = bottom, 1 = top
		t := float64(i) / float64(branchCount-1)

		branchY := baseY - h*(0.18+t*0.70)

		// Trees become narrower toward the top.
		widthFactor := 1.0 - t*0.78

		maxLength := h * 0.40 * widthFactor

		// Random branch length.
		length := maxLength *
			(0.65 + rng.Float64()*0.35)

		// Branches point slightly upward.
		angle := 0.12 + rng.Float64()*0.18

		// Left branch.
		drawBranch(
			img,
			RangePoint{x, branchY},
			length,
			math.Pi-angle,
			h*0.045*widthFactor,
			greens[rng.Intn(len(greens))],
			rng,
		)

		// Right branch.
		drawBranch(
			img,
			RangePoint{x, branchY - rng.Float64()*h*0.025},
			length*(0.9+rng.Float64()*0.1),
			angle,
			h*0.045*widthFactor,
			greens[rng.Intn(len(greens))],
			rng,
		)
	}

	// ---------------------------------------------------------
	// Needle clusters / fine foliage
	// ---------------------------------------------------------

	needleCount := int(h * 1.5)

	for i := 0; i < needleCount; i++ {

		t := rng.Float64()

		y := baseY - h*(0.12+t*0.78)

		widthFactor := 1.0 - t*0.8

		xOffset := (rng.Float64()*2 - 1) *
			h * 0.38 * widthFactor

		length := h *
			(0.015 + rng.Float64()*0.025)

		angle := -0.6 + rng.Float64()*1.2

		p1 := RangePoint{x + xOffset, y}
		p2 := RangePoint{
			X: p1.X + math.Cos(angle)*length,
			Y: p1.Y + math.Sin(angle)*length,
		}

		drawLine(
			img,
			p1,
			p2,
			h*0.012,
			greens[rng.Intn(len(greens))],
		)
	}
}

func drawBranch(
	img *image.RGBA,
	start RangePoint,
	length float64,
	angle float64,
	width float64,
	c color.Color,
	rng *rand.Rand,
) {
	segments := 8

	prev := start

	for i := 1; i <= segments; i++ {

		t := float64(i) / float64(segments)

		// Branch curves upward slightly.
		curve := -0.20 * t * t

		a := angle + curve

		// Random small perturbation.
		a += (rng.Float64() - 0.5) * 0.08

		next := RangePoint{
			X: start.X + math.Cos(a)*length*t,
			Y: start.Y + math.Sin(a)*length*t,
		}

		// Branch tapers toward its tip.
		w := width * (1.0 - t*0.75)

		drawLine(img, prev, next, w, c)

		prev = next
	}
}

func drawLine(
	img *image.RGBA,
	a, b RangePoint,
	width float64,
	c color.Color,
) {
	minX := int(math.Min(a.X, b.X) - width)
	maxX := int(math.Max(a.X, b.X) + width)

	minY := int(math.Min(a.Y, b.Y) - width)
	maxY := int(math.Max(a.Y, b.Y) + width)

	r := width / 2

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {

			d := distanceToSegment(
				RangePoint{float64(x), float64(y)},
				a,
				b,
			)

			if d <= r {
				img.Set(x, y, c)
			}
		}
	}
}

func distanceToSegment(p, a, b RangePoint) float64 {
	dx := b.X - a.X
	dy := b.Y - a.Y

	if dx == 0 && dy == 0 {
		return math.Hypot(p.X-a.X, p.Y-a.Y)
	}

	t := ((p.X-a.X)*dx + (p.Y-a.Y)*dy) /
		(dx*dx + dy*dy)

	t = math.Max(0, math.Min(1, t))

	qx := a.X + t*dx
	qy := a.Y + t*dy

	return math.Hypot(p.X-qx, p.Y-qy)
}
