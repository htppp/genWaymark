# genWaymark

[WaymarkPresetPlugin](https://github.com/PunishedPineapple/WaymarkPresetPlugin)用のjsonを生成します。

## 前提

* ff14の座標系は左手座標系
* +x軸方向が西から東、+z方向が北から南
* Z座標はフィールドマーカーで使わないので(x,y,z)の(x,z)だけ使用する
* ~~このプログラムの作者はOpenGL畑なので右手系Z高さな座標でないと混乱する~~

## 使い方

	$ ./genWaymark -h
	Usage of ./genWaymark:
	  -a, --angle float            角度間隔を指定します。正の値が時計回りです。単位は度です。 (default 45)
	  -x, --center-x-coord float   中心座標のX座標を指定します。
	  -z, --center-z-coord float   中心座標のZ座標を指定します。
	  -i, --id int                 マーカーを使用するマップのIDを指定します。
	  -m, --mode int               マーカーを配置するパターンを指定します。
								   1: 等間隔(既定)
								   2: 12分割(i%3==0除外)
								   3: 12分割(i%3==1除外)
								   4: 12分割(i%3==2除外)  (default 1)
	  -n, --name string            生成したマーカー群に名前を付けられます。(既定)半径や角度を用いて自動生成します。
	  -o, --offset-angle float     最初に配置するマーカーの位置を指定値角度オフセットします。
								   既定でマーカーはABCD1234の順に並べられます。
								   Aは中心から北の方向に半径分離れた位置に配置されます。
								   このAの中心から見た角度をオフセットします。正の値が時計回りです。
		  --order string           配置するマーカーの順番を指定します。ABCD1234以外の文字はエラーとします。 (default "ABCD1234")
	  -r, --radius float           半径を指定します。単位はメートルです。 (default 15)
	pflag: help requested

## 例

* イフリート討滅戦

	MapIDは56 マップ中心の座標が(0,0) ボスは東、プレーヤーは西からスタート

	1. 中心(0,0) 半径10m 角度45度毎にABCD1234を配置

		`./genWaymark -i56 -x0 -z0 -r10 -a45`

		![](https://i.gyazo.com/b40a7c2d8ac87b4bef34a8b4f2e33235.png)

	2. 上のマーカーの並びをA1B2C3D4に変更したもの

		`./genWaymark -i56 -x0 -z0 -r10 -a45 --order A1B2C3D4`

		![](https://i.gyazo.com/2e229d391e75bd3fe771ffb2847ee86e.png)

	3. 2のマーカーを時計回りに90度オフセットしたもの

		`./genWaymark -i56 -x0 -z0 -r10 -a45 --order A1B2C3D4 --angle-offset 90`

		![](https://i.gyazo.com/868c6d65408aa6b2ffd76eed51d4ef7b.png)

* 絶竜詩戦争
	* MapIDは788 マップ中心の座標が(100,100)。
	* P1ハイパーディメンションの立ち位置にマーカーを置くとする。
	* 0,3,6,9時にハイパーディメンションを捨てるとワイプなので等間隔ではなく隙間をあけて配置したい。

	* 中心(100,100)、半径11m、時計回りに30度毎にA1B2C3D4の順番で、  
	何も考えずマーカーを配置すると0時から8時までに等間隔で8個並ぶ。

	`./genWaymark --id 788 --name 竜詩戦争P1 -x 100 -z 100 --angle 30 --radius 11 --order A1B2C3D4`

	![](https://i.gyazo.com/427edba81d89579fc7b2e1e60b88fcf1.png)

	* 0,3,6,9時方向にマーカーを配置せず、北にA1,東にB2、南にC3、西にD4を配置したい。
	* Aを0時からではなく11時から配置するのは`--angle-offset -30`で角度を-30度分オフセットすればよい。
	`./genWaymark --id 788 --name 竜詩戦争P1 -x 100 -z 100 --angle 30 --radius 11 --order A1B2C3D4 --angle-offset -30`

	![](https://i.gyazo.com/9cef718681bc76d8d44d2334f4b7fd2f.png)

	* また等間隔ではなく、Aを11時に配置し、12時に配置せず、1時に1を配置し、  
	2時にBを配置し、3時に配置せず、のように置いていく必要がある。  

	原点が指定座標、フィールドの北を角度0とし、時計回りが正な極座標系を考える。つまり

		v = (radius, i * angle + offsetAngle)

	な位置ベクトルを考える。A1B2C3D4の順番でマーカーを配置するとき、  
	Aの位置ベクトルは11時に置きたいので`(11, 0*30-30) = (11, -30)`、  
	Aと1の間にはマーカーを置かず、  
	1の位置ベクトルは1時に置きたいので`(11, 2*30-30) = (11, +30)`、  
	Bの位置ベクトルは2時に置きたいので`(11, 3*30-30) = (11, +60)`、  
	Bと2の間にはマーカーを置かず、  
	2の位置ベクトルは4時に置きたいので`(11, 5*30-30) = (11, +120)`、  
	のようにになる。

	このために`./genWaymark -h` で出力されるヘルプの`mode`オプションの3番を使う。  
	モード3は0から11の12個の連番を考えた時、3で割った余りが1である1,4,7,10にマーカーを配置せず、  
	0,2,3,5,6,8,9,11にマーカーを置いていく。

	`./genWaymark --id 788 --name 竜詩戦争P1 -x 100 -z 100 --angle 30 --angle-offset -30 --radius 11 --mode 3 --order A1B2C3D4`

	![](https://i.gyazo.com/375bc1a2a11c2669d43a252df6820c76.png)

	これでP1のハイパーディメンションの散会位置にマーカーが配置で来た。

	参考に角度オフセットをせずに配置するとmode2,3,4は以下のようになる。
	* mode == 2  
	i % 3 == 0 の0時にマーカーが置かれず、1,2時に置かれ、3時に置かれない  
	![角度オフセットなしmode2](https://i.gyazo.com/ed5600a390aba18dafece46ebfe6524c.png)
	* mode == 3  
	i % 3 == 1 の1時にマーカーが置かれず、0,2,3時に置かれ、4時に置かれない  
	![角度オフセットなしmode3](https://i.gyazo.com/7844aa0cf623366c596efaf1212d5b1f.png)
	* mode == 4  
	i % 3 == 2 の時にマーカーが置かれず、0,1,3,4,6時に置かれ、5時に置かれない  
	![角度オフセットなしmode4](https://i.gyazo.com/b412bd6c13fbe867cc7c249cb3395b54.png)

	上記が時計の時針の位置と一致するのは角度指定が30度だからであり、modeが12分割だからではないことに注意。  
	参考に角度15度のmode3の場合はこうなる。
	`./genWaymark --id 788 -x 100 -z 100 --angle 15 --radius 15 --order A1B2C3D4 --mode 3`
	![](https://i.gyazo.com/38740cbd966a1950ff162c06b67626de.png)

	* idをイフリート討滅戦に変えて配置するとこうなる

	`./genWaymark --id 56 -x 100 -z 100 --angle 30 --angle-offset -30 --radius 11 --mode 3 --order A1B2C3D4`

	![](https://i.gyazo.com/5b4df446df95e4af30b32061dbe030ba.png)

