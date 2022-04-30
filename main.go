package main

import ( // {{{
	"encoding/json"
	"fmt"
	"math"
	"os"

	"github.com/spf13/pflag"
	"gonum.org/v1/gonum/mat"
) // }}}

// type definition {{{
type Coord struct {
	X      float64 `json:"X"`
	Y      float64 `json:"Y"`
	Z      float64 `json:"Z"`
	ID     int     `json:"ID"`
	Active bool    `json:"Active"`
}

func NewCoord(coord []float64, id int, active bool) Coord {
	return Coord{
		X:      coord[0],
		Y:      0,
		Z:      coord[1],
		ID:     id,
		Active: active,
	}
}

type WayMark struct {
	Name  string `json:"Name"`
	MapID int    `json:"MapID"`
	A____ Coord  `json:"A"`
	B____ Coord  `json:"B"`
	C____ Coord  `json:"C"`
	D____ Coord  `json:"D"`
	One__ Coord  `json:"One"`
	Two__ Coord  `json:"Two"`
	Three Coord  `json:"Three"`
	Four_ Coord  `json:"Four"`
}

type Args struct {
	radius      float64
	angle       float64
	angleOffset float64
	centerCoord []float64
	name        string
	id          int
	mode        int
	order       string
}

/*
{
	"Name":"Imported","MapID":788,
	"A":{"X":106.169,"Y":0.0,"Z":89.725,"ID":0,"Active":true},
	"B":{"X":109.72,"Y":0.0,"Z":94.398,"ID":1,"Active":true},
	"C":{"X":106.186,"Y":0.0,"Z":110.02,"ID":2,"Active":true},
	"D":{"X":90.473,"Y":0.0,"Z":105.334,"ID":3,"Active":true},
	"One":{"X":94.036,"Y":0.0,"Z":89.797,"ID":4,"Active":true},
	"Two":{"X":109.869,"Y":0.0,"Z":105.691,"ID":5,"Active":true},
	"Three":{"X":93.94,"Y":0.0,"Z":109.851,"ID":6,"Active":true},
	"Four":{"X":90.102,"Y":0.0,"Z":94.297,"ID":7,"Active":true}
}
*/
// }}}

var args Args

func init() {
	pRadius := pflag.Float64P("radius", "r", 15, "半径を指定します。単位はメートルです。")
	pAngle := pflag.Float64P("angle", "a", 45, "角度間隔を指定します。正の値が時計回りです。単位は度です。")
	pAngleOffset := pflag.Float64P("offset-angle", "o", 0, "最初に配置するマーカーの位置を指定値角度オフセットします。\n既定でマーカーはABCD1234の順に並べられます。\nAは中心から北の方向に半径分離れた位置に配置されます。\nこのAの中心から見た角度をオフセットします。正の値が時計回りです。")
	pX := pflag.Float64P("center-x-coord", "x", 0, "中心座標のX座標を指定します。")
	pY := pflag.Float64P("center-z-coord", "z", 0, "中心座標のZ座標を指定します。")
	pName := pflag.StringP("name", "n", "", "生成したマーカー群に名前を付けられます。(既定)半径や角度を用いて自動生成します。")
	pId := pflag.IntP("id", "i", 0, "マーカーを使用するマップのIDを指定します。")
	pMode := pflag.IntP("mode", "m", 1, "マーカーを配置するパターンを指定します。\n1: 等間隔(既定)\n2: 12分割(i%3==0除外)\n3: 12分割(i%3==1除外)\n4: 12分割(i%3==2除外) ")
	pOrder := pflag.String("order", "ABCD1234", "配置するマーカーの順番を指定します。ABCD1234以外の文字はエラーとします。")
	pflag.Parse()
	if *pName == "" {
		if *pOrder != "ABCD1234" {
			*pName = fmt.Sprintf("r%.1f_a%.1f_o%.f_%s", *pRadius, *pAngle, *pAngleOffset, *pOrder)
		} else {
			*pName = fmt.Sprintf("r%.1f_a%.1f_o%.f", *pRadius, *pAngle, *pAngleOffset)
		}
	}
	args = Args{
		radius:      *pRadius,
		angle:       *pAngle,
		angleOffset: *pAngleOffset,
		centerCoord: []float64{*pX, *pY},
		name:        *pName,
		id:          *pId,
		mode:        *pMode,
		order:       *pOrder,
	}

	pflag.Usage = func() {
		usageText := `WaymarkPresetPlugin用のjsonを生成します。
https://github.com/PunishedPineapple/WaymarkPresetPlugin

Usage:
	genWaymark [options]

The commands are:
	中心を(100,100)として半径11m時計回りに30度毎にABCD1234の順番で0,3,6,9時方向にマーカーを配置せず北にA1,東にB2、南にC3、西にD4を配置する
	genWaymark --id 788 --name 竜詩戦争P1 -x 100 -y 100 --angle 30 --angle-offset -30 --radius 11 --mode 3 --order A1B2C3D4
	genWaymark --id 788 -x 100 -y 100 -a 45 -r 20

Use "genWaymark -h" for more infomation about a command`
		fmt.Fprintf(os.Stderr, "%s\n\n", usageText)
	}
}

func main() {
	if len(os.Args) == 1 {
		pflag.PrintDefaults()
		os.Exit(1)
	}
	GenWaymark(args)
}

func Rotate(vec *mat.Dense, rad float64) *mat.Dense {
	rot := mat.NewDense(2, 2, []float64{math.Cos(rad), -1 * math.Sin(rad), math.Sin(rad), math.Cos(rad)})
	out := mat.NewDense(2, 1, nil)
	out.Mul(rot, vec)
	return out
}

func GenWaymark(args Args) {

	fmt.Printf("Parsed input args :\n")
	fmt.Printf("    name          : %s\n", args.name)
	fmt.Printf("    map id        : %d\n", args.id)
	fmt.Printf("    center coord  : %f, %f\n", args.centerCoord[0], args.centerCoord[1])
	fmt.Printf("    radius        : %f\n", args.radius)
	fmt.Printf("    angle         : %f\n", args.angle)
	fmt.Printf("    angle offset  : %f\n", args.angleOffset)
	fmt.Printf("    mode          : %d\n", args.mode)
	fmt.Printf("    order         : %s\n", args.order)

	v := mat.NewDense(2, 1, []float64{0, -args.radius}) // -r で12時方向
	{                                                   // offset
		var angleRad float64 = args.angleOffset / 180 * math.Pi // radに変換
		v = Rotate(v, angleRad)
	}

	calced := []mat.Dense{}
	if args.mode == 1 {
		for i := 0; i < 8; i++ {
			var angleRad float64 = float64(i) * args.angle / 180 * math.Pi // radに変換
			calced = append(calced, *Rotate(v, angleRad))
		}
	} else if args.mode == 2 {
		for i := 0; i < 12; i++ {
			if i%3 != 0 {
				var angleRad float64 = float64(i) * args.angle / 180 * math.Pi // radに変換
				calced = append(calced, *Rotate(v, angleRad))
			}
		}
	} else if args.mode == 3 {
		for i := 0; i < 12; i++ {
			if i%3 != 1 {
				var angleRad float64 = float64(i) * args.angle / 180 * math.Pi // radに変換
				calced = append(calced, *Rotate(v, angleRad))
			}
		}
	} else if args.mode == 4 {
		for i := 0; i < 12; i++ {
			if i%3 != 2 {
				var angleRad float64 = float64(i) * args.angle / 180 * math.Pi // radに変換
				calced = append(calced, *Rotate(v, angleRad))
			}
		}
	}
	coord := [][]float64{}

	center := mat.NewDense(2, 1, args.centerCoord)
	for i := 0; i < len(calced); i++ {
		var o mat.Dense
		o.Add(&calced[i], center)
		coord = append(coord, []float64{o.At(0, 0), o.At(1, 0)})
	}
	out := WayMark{}
	{ // 12時から時計まわりに座標が入っている 北にA1と並べるには1からA
		//fmt.Println(len(name))
		out.Name = args.name
		out.MapID = args.id
		for i, c := range args.order {
			switch c {
			case 'A':
				out.A____ = NewCoord(coord[i], 0, true)
			case 'B':
				out.B____ = NewCoord(coord[i], 1, true)
			case 'C':
				out.C____ = NewCoord(coord[i], 2, true)
			case 'D':
				out.D____ = NewCoord(coord[i], 3, true)
			case '1':
				out.One__ = NewCoord(coord[i], 4, true)
			case '2':
				out.Two__ = NewCoord(coord[i], 5, true)
			case '3':
				out.Three = NewCoord(coord[i], 6, true)
			case '4':
				out.Four_ = NewCoord(coord[i], 7, true)
			default:
				fmt.Printf("未知のフィールドマーカー文字が指定されました。\n")
				fmt.Printf("エラー : %v はフィールドマーカー指定の文字として予期していない文字です。\n", c)
				fmt.Printf("ABCD1234のいずれかのみを空白を開けずに指定してください\n")
				os.Exit(1)
			}
		}
	}

	b, _ := json.Marshal(out)
	fmt.Println("")
	fmt.Println(string(b))
}
