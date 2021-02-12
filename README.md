# elmo ![Go](https://github.com/toubib/elmo/workflows/Go/badge.svg)

Elmo let you fetch a web page with its assets.

It can send download time statistics to [InfluxDB](https://www.influxdata.com/products/influxdb/) and be used as a [Nagios](https://www.nagios.org) check plugin.

! This project is not released yet, use with caution !

## Usage

```
Usage of ./elmo:
  -assets-allowed-domains string
    	List of allowed assets domains to fetch from, comma separated.
  -connect-timeout int
    	Connect timeout in ms. (default 1000)
  -influx-database string
    	The influx database name. (default "elmo")
  -influx-url string
    	The influx database access url. (default "http://localhost:8086")
  -nagios-critical int
    	Nagios critical time in ms. (default 10000)
  -nagios-warning int
    	Nagios warning time in ms. (default 5000)
  -parallel int
    	Number of parallel fetch to launch. 0 means unlimited. (default 8)
  -request-timeout int
    	Request timeout in ms. (default 10000)
  -response-header-timeout int
    	Response header timeout in ms.
  -url string
    	The url to get.
  -use-influx
    	Send data to influxdb.
  -use-nagios
    	Nagios compatible output.
  -verbose
    	Print more informations.
  -version
    	Print version information.
```

## Examples

### Nagios output
```
$ ./elmo -url https://yahoo.com -use-nagios
Downloaded 1946KB in 83/83 files in 2.51824949s.|size=1946KB time=2.51824949s;5000;10000;0;10000
```

### Verbose output
```
$ ./elmo -url https://yahoo.com -verbose
1.34699539s	200 https://yahoo.com 1.072768484s 283604b
1.523315099s	200 https://s.yimg.com/nn/lib/metro/g/sda/sda_flex_0.0.43.css 166.358186ms 2482b
1.547134362s	200 https://s.yimg.com/aaq/fp/css/tdv2-wafer-ntk.NTKDesktop.atomic.ltr.94b956089fc91c2f0a244928a927abc9.min.css 166.46574ms 4657b
1.547801904s	200 https://s.yimg.com/os/yc/css/bundle.c60a6d54.css 190.772067ms 4261b
1.550130938s	200 https://s.yimg.com/nn/lib/metro/g/myy/video_styles_0.0.72.css 190.621987ms 13153b
1.554382918s	200 https://s.yimg.com/os/yc/css/bundle.c60a6d54.css 197.491977ms 4261b
1.557617403s	200 https://s.yimg.com/aaq/fp/css/tdv2-wafer-ntk.custom.desktop.a69916e03ec8f658d9530295bab867ab.css 201.013966ms 585b
1.557979664s	200 https://s.yimg.com/aaq/fp/css/react-wafer-featurebar.FeaturebarApplet.atomic.ltr.8ad35615271ca94216b2ec60e84c9bfe.min.css 201.080705ms 2603b
1.559516295s	200 https://s.yimg.com/nn/lib/metro/g/myy/flex_grid_0.0.50.css 201.226398ms 4290b
1.630626717s	200 https://s.yimg.com/aaq/fp/css/tdv2-wafer-stream.custom.desktop.35b4e59342f8c72801c502afb5933cff.css 78.369608ms 3532b
1.630722712s	200 https://s.yimg.com/aaq/fp/css/tdv2-wafer-trending.Trending.atomic.ltr.d9da5ad162c4e96358f3ac5e527aaa03.min.css 80.373699ms 2947b
1.630938565s	200 https://s.yimg.com/aaq/fp/css/react-wafer-weather.WeatherPreview.atomic.ltr.878f61f52a9883465dc66da289ca740e.min.css 76.196138ms 4834b
1.630950195s	200 https://s.yimg.com/aaq/fp/css/tdv2-wafer-user-intent.rollupDesktop.atomic.ltr.85ffd965befa53ddf87e9ad7356550b7.min.css 82.745676ms 12336b
1.630973576s	200 https://s.yimg.com/aaq/fp/css/react-wafer-topics.TopicsDesktop.atomic.ltr.176779a3df68d7abe76be4a21a3c328d.min.css 72.725237ms 2287b
1.631026193s	200 https://s.yimg.com/aaq/fp/css/react-wafer-weather.common.desktop.62d099be776ca538092fa6ba87d1637b.css 73.115381ms 3042b
1.633182368s	200 https://s.yimg.com/aaq/fp/css/tdv2-wafer-footer.FooterDesktop.atomic.ltr.0dabe32d96d30f44862f1509e652a092.min.css 71.317426ms 964b
1.70615618s	200 https://s.yimg.com/aaq/fp/css/tdv2-wafer-stream.StreamWide.atomic.ltr.0516e4514859bcfa6fe23b4fa493a717.min.css 88.751417ms 58433b
1.720044598s	200 https://s.yimg.com/aaq/fp/css/tdv2-wafer-header.custom.desktop.2ce65662738d6cd781c23fc340c7205c.css 89.260449ms 613b
1.720131373s	200 https://s.yimg.com/aaq/fp/css/tdv2-wafer-header.dialogue.desktop.5bb76716485c7bce78d352c4003537ca.css 89.119652ms 258b
1.720226302s	200 https://s.yimg.com/aaq/fp/css/react-wafer-storyswarm.Stories.atomic.ltr.0af72758a4fb13c84f9d503b04e5f73b.min.css 83.910255ms 19733b
1.721789389s	200 https://s.yimg.com/aaq/fp/css/tdv2-wafer-header.HeaderDesktop.atomic.ltr.439dda9ea65f0d89195e06c9309d2dab.min.css 86.004555ms 33791b
1.722199229s	200 https://s.yimg.com/aaq/fp/css/tdv2-wafer-header.MailPreview.atomic.ltr.fc22a5794c70dd972bd3c4a015bf1df9.min.css 89.047751ms 30565b
1.761393682s	200 https://s.yimg.com/aaq/fp/css/tdv2-wafer-user-dialog.UserDialogLite.atomic.ltr.875d949b676096085b14b09e69c5df59.min.css 130.148495ms 2869b
1.785756891s	200 https://s.yimg.com/aaq/fp/js/tdv2-wafer-stream.custom.2499e589777968e7090b7f293c896dc5.js 60.152516ms 11022b
1.786384637s	200 https://s.yimg.com/aaq/fp/js/tdv2-wafer-header.custom.desktop.e0cc81c4de21a0aee644ee9285d79117.js 65.900072ms 2999b
1.786727672s	200 https://s.yimg.com/aaq/scp/css/viewer.248c609437ca1468d87f39058c3bd2cf.css 147.215903ms 16648b
1.786753171s	200 https://s.yimg.com/aaq/fp/js/tdv2-wafer-utils.customErrorHandler.291a68a0c9926c1608e024ce86d950ec.js 65.720401ms 3462b
1.786776127s	200 https://s.yimg.com/nn/lib/metro/g/sda/sda_modern_0.0.41.js 79.750261ms 10670b
1.832232929s	200 https://s.yimg.com/aaq/fp/js/tdv2-wafer-user-dialog.custom.6696ae78fb98b0026e813d04c0206fd3.js 108.757111ms 2473b
1.839400806s	200 https://s.yimg.com/aaq/yc/js/iframe-1.0.26.js 116.886067ms 5445b
1.901406298s	200 https://s.yimg.com/uu/api/res/1.2/z1y3CFj._f12NqrmTwzj2Q--~B/Zmk9c3RyaW07aD0xNjA7cT04MDt3PTM0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/os/creatr-uploaded-images/2020-10/59b28220-196e-11eb-9ded-b7b78eb0e0ef.cf.jpg 63.714652ms 6944b
1.918991154s	200 https://s.yimg.com/aaq/cmp/version/4.0.0/cmp.js 77.671835ms 78818b
1.95148816s	200 https://s.yimg.com/cv/apiv2/Doodle/LunarNewYear_2021/Yahoo_Doodle_410x116.png 60.302113ms 36850b
1.961017348s	200 https://s.yimg.com/uu/api/res/1.2/plkwq7NMULFIke6cqZ7fWw--~B/Zmk9ZmlsbDtoPTE0NDtweW9mZj0wO3c9MjcyO2FwcGlkPXl0YWNoeW9u/https://s.yimg.com/os/creatr-uploaded-images/2020-03/2d4ac360-67d9-11ea-96ef-f1007eaff3de 106.753864ms 15158b
1.976438305s	200 https://s.yimg.com/uu/api/res/1.2/o4HilXmSewmuM26gKjhalQ--~B/Zmk9c3RyaW07aD0xNjA7cT04MDt3PTM0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/os/creatr-uploaded-images/2021-02/1da4b350-6d26-11eb-adbf-dfcc1ff91820.cf.jpg 72.624075ms 5053b
1.995622302s	200 https://s.yimg.com/uu/api/res/1.2/qmss_gYIrbqx_3Y0hXzcmA--~B/Zmk9c3RyaW07aD0xNjA7cT04MDt3PTM0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/uu/api/res/1.2/yYz5WxteVMStEVj2k.ezAg--~B/aD05MDA7dz0xMjAwO2FwcGlkPXl0YWNoeW9u/https://media.zenfs.com/fr/TeleLoisirs.fr/20c28414d22255cb2411c68f1b10d1f6.cf.jpg 71.086422ms 17879b
2.020160321s	200 https://s.yimg.com/cv/apiv2/default/YM6/936.jpg 67.709817ms 47844b
2.020296406s	200 https://s.yimg.com/uu/api/res/1.2/fQu5hl8clVGDxK8XNI.0Jg--~B/Zmk9c3RyaW07aD0xNjA7cT04MDt3PTM0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/os/creatr-uploaded-images/2021-02/27808040-6d56-11eb-a9fb-f26cd0df876f.cf.jpg 115.946011ms 15128b
2.02319688s	200 https://s.yimg.com/rq/darla/4-6-0/js/g-r-min.js 54.967608ms 207362b
2.063646831s	200 https://s.yimg.com/uu/api/res/1.2/Zz2mQIDcxSQwt69IaZ3F_Q--~B/Zmk9c3RyaW07aD0yNDY7cT04MDt3PTQ0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/os/creatr-uploaded-images/2021-02/fee64de0-6d41-11eb-bb7f-87de6847cd68.cf.jpg 88.111772ms 13508b
2.067443308s	200 https://s.yimg.com/uu/api/res/1.2/0FdAA5P8lEIzbKca0Fphag--~B/Zmk9c3RyaW07aD0xNjA7cT04MDt3PTM0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/os/creatr-uploaded-images/2021-02/45d26e70-6d1d-11eb-9cff-1170894587e5.cf.jpg 71.655523ms 11413b
2.070643013s	200 https://s.yimg.com/uu/api/res/1.2/efqN4tDie3WTIJaHz.vHZA--~B/Zmk9c3RyaW07aD0xNDA7cT05MDt3PTE0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/uu/api/res/1.2/S85AaNGRemIk9DaXT3votg--~B/aD04OTc7dz0xMjAwO2FwcGlkPXl0YWNoeW9u/https://media.zenfs.com/fr/TeleLoisirs.fr/246bcbfb523a2b86302a36e1c772c914.cf.jpg 74.858037ms 5926b
2.086400467s	200 https://s.yimg.com/uu/api/res/1.2/s867hnb087JvxAKqpqebyQ--~B/Zmk9c3RyaW07aD0xNDA7cT05MDt3PTE0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/uu/api/res/1.2/Zq1zf.rmNEQ6wwYSPM8eWQ--~B/aD05MDA7dz0xMjAwO2FwcGlkPXl0YWNoeW9u/https://media.zenfs.com/fr/voici.fr/f2db86ec8b07dc6d5d5530f90befa3fd.cf.jpg 62.872441ms 8501b
2.097563048s	200 https://s.yimg.com/uu/api/res/1.2/ZBx3vCR2WUZWoBKTEXmq1g--~B/Zmk9c3RyaW07aD0xNDA7cT05MDt3PTE0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/os/creatr-uploaded-images/2021-02/8415e940-6d28-11eb-affa-ecb46130fe0a.cf.jpg 71.096933ms 5631b
2.107140002s	200 https://s.yimg.com/uu/api/res/1.2/ruUmtcLHl1KfW4l1aOcocA--~B/Zmk9c3RyaW07aD0zODg7cT05NTt3PTcyMDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/os/creatr-uploaded-images/2021-02/14c94a20-6d1c-11eb-af3d-4903e39382be.cf.jpg 107.243334ms 125320b
2.121641515s	200 https://s.yimg.com/uu/api/res/1.2/0tMb.gAdTldce3_Q_hH85g--~B/Zmk9c3RyaW07aD0zODY7cT04MDt3PTQ0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/uu/api/res/1.2/tfRK3dJwbRcrcECB6LpdjQ--~B/aD05MDA7dz0xMjAwO2FwcGlkPXl0YWNoeW9u/https://media.zenfs.com/fr/voici.fr/f8352a60502e0d1bd4f193bdb7771464.cf.jpg 91.500964ms 31943b
2.123162727s	200 https://s.yimg.com/uu/api/res/1.2/XX0rsNv5EypITeAB3kCSfQ--~B/Zmk9c3RyaW07aD0xNDA7cT05MDt3PTE0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/uu/api/res/1.2/k7vj0ZhtbTc7i.lL34rRAw--~B/aD01MzI7dz04MDA7YXBwaWQ9eXRhY2h5b24-/https://media.zenfs.com/fr/reuters.com/41878fe725ec8d3de7397e960d40ad2b.cf.jpg 58.009314ms 7731b
2.129452802s	200 https://s.yimg.com/uu/api/res/1.2/z9xzGPOIkX.5RqtxfrB2MQ--~B/Zmk9c3RyaW07aD0xNDA7cT05MDt3PTE0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/uu/api/res/1.2/vo_FInyYr.hROjdtE9WpRQ--~B/aD05MDA7dz0xMjAwO2FwcGlkPXl0YWNoeW9u/https://media.zenfs.com/fr/gala.fr/0a361b851bdc110648fd32f2294cf47f.cf.jpg 56.661444ms 7106b
2.133523706s	200 https://s.yimg.com/uu/api/res/1.2/p6kANDt9iFC6aJ86b1MKvw--~B/Zmk9c3RyaW07aD0zODY7cT04MDt3PTQ0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/uu/api/res/1.2/_x5dIt8zotwpkRm3siN2jg--~B/aD05MDA7dz0xMjAwO2FwcGlkPXl0YWNoeW9u/https://media.zenfs.com/fr/Programme.fr/28d90282ddb567ae9eeef415d1e7d29d.cf.jpg 55.61577ms 20503b
2.146211946s	200 https://s.yimg.com/uu/api/res/1.2/9TGwTtPIrVJzGNLKNG4Gkg--~B/Zmk9c3RyaW07aD0zODY7cT04MDt3PTQ0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/uu/api/res/1.2/x8WBx6A.FkxgP31ytLbWOw--~B/aD00OTE7dz04NzI7YXBwaWQ9eXRhY2h5b24-/https://media.zenfs.com/fr/closermag.fr/70066605b6484f2549135d5f6be803da.cf.jpg 61.485818ms 23641b
2.154606042s	200 https://s.yimg.com/uu/api/res/1.2/87OoiAC9q7rjCpARF0hScQ--~B/Zmk9c3RyaW07aD0xNDA7cT05MDt3PTE0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/uu/api/res/1.2/a3kPnu38XgtkmY4swiL1eQ--~B/aD05MDg7dz0xNjE0O2FwcGlkPXl0YWNoeW9u/https://media.zenfs.com/fr/closermag.fr/4ca4ae4dd1db9a896b8e0b2a0f41407d.cf.jpg 63.867698ms 8122b
2.167122724s	200 https://s.yimg.com/uu/api/res/1.2/0Mfi0xAwx1bRz7qn6cStXw--~B/Zmk9c3RyaW07aD0yNDY7cT04MDt3PTQ0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/uu/api/res/1.2/HWAmTEBCmCVP9Rx2tC.6.w--~B/aD01NzY7dz0xMDI0O2FwcGlkPXl0YWNoeW9u/https://media.zenfs.com/fr/the_observers_french_924/a1b91ec9a6f4652a060795341384d2cc.cf.jpg 47.593336ms 15942b
2.172010817s	200 https://s.yimg.com/uu/api/res/1.2/hVgk2De_SLZPAU42h0Bf6A--~B/Zmk9c3RyaW07aD0yNDY7cT04MDt3PTQ0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/uu/api/res/1.2/ueQYId_XOfJpxiF9QPWGmw--~B/aD00MjA7dz02MzA7YXBwaWQ9eXRhY2h5b24-/https://media.zenfs.com/fr/parismatch.com/9ddb99462d040b11a67afbd42f143dd4.cf.jpg 56.971787ms 34277b
2.172870082s	200 https://s.yimg.com/g/images/spaceball.gif 51.152673ms 43b
2.192306262s	200 https://s.yimg.com/uu/api/res/1.2/gqKVQYT4_szApT4OSbJuqg--~B/Zmk9c3RyaW07aD0xNDA7cT05MDt3PTE0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/uu/api/res/1.2/Q_00kh2l0xMSOq24R78kyg--~B/aD04OTk7dz0xMjAwO2FwcGlkPXl0YWNoeW9u/https://media.zenfs.com/fr/Programme.fr/de4c6ca8ef4cae4ac1ee62a3103b3b63.cf.jpg 65.86285ms 5956b
2.197455741s	200 https://s.yimg.com/uu/api/res/1.2/T8W08OLTwXvmLENF08PPHQ--~B/Zmk9c3RyaW07aD0xNDA7cT05MDt3PTE0MDthcHBpZD15dGFjaHlvbg--/https://s.yimg.com/uu/api/res/1.2/vO1_sI21KWBMTwAx0iPkBg--~B/aD00NDA7dz02NzA7YXBwaWQ9eXRhY2h5b24-/https://media.zenfs.com/fr/programme_tv_fr_account_146/be3e801cb585e184af744ca9c1625cbe.cf.jpg 67.875843ms 5759b
2.221898127s	200 https://s.yimg.com/ss/rapid-3.53.17.js 70.979408ms 48857b
2.225416794s	200 https://s.yimg.com/aaq/wf/wf-core-1.43.11.js 62.625248ms 140303b
2.235765116s	200 https://s.yimg.com/aaq/hp-viewer/desktop_1.9.22.js 58.227791ms 122912b
2.240730046s	200 https://s.yimg.com/aaq/wf/wf-image-1.1.6.js 67.664859ms 5487b
2.240747852s	200 https://s.yimg.com/aaq/wf/wf-darla-1.0.21.js 68.530488ms 4398b
2.258876627s	200 https://s.yimg.com/aaq/wf/wf-caas-1.14.1.js 61.867839ms 21557b
2.262650343s	200 https://s.yimg.com/aaq/wf/wf-video-2.14.6.js 61.417406ms 25955b
2.28454385s	200 https://s.yimg.com/aaq/wf/wf-bind-1.1.2.js 62.162884ms 4077b
2.284972167s	200 https://s.yimg.com/aaq/wf/wf-text-1.1.3.js 58.865477ms 3859b
2.307603288s	200 https://s.yimg.com/cv/apiv2/default/20191018/FR_FR_Yellow_300x250.png 58.704825ms 113686b
2.308269038s	200 https://s.yimg.com/aaq/wf/wf-toggle-1.13.4.js 71.931377ms 12177b
2.30965847s	200 https://s.yimg.com/aaq/wf/wf-beacon-1.3.1.js 67.03027ms 11142b
2.310213924s	200 https://s.yimg.com/aaq/wf/wf-fetch-1.16.7.js 68.892176ms 16419b
2.321076447s	200 https://s.yimg.com/aaq/wf/wf-countdown-1.2.5.js 58.15804ms 4856b
2.322575425s	200 https://s.yimg.com/aaq/wf/wf-scrollview-2.10.0.js 58.667898ms 25861b
2.335909841s	200 https://s.yimg.com/aaq/wf/wf-drawer-1.0.10.js 48.477184ms 7855b
2.336195877s	200 https://s.yimg.com/aaq/wf/wf-template-1.3.5.js 49.237638ms 7698b
2.361255183s	200 https://s.yimg.com/aaq/wf/wf-form-1.25.7.js 50.348603ms 14312b
2.361384022s	200 https://s.yimg.com/aaq/wf/wf-sticky-1.0.9.js 52.809717ms 7003b
2.362881725s	200 https://s.yimg.com/aaq/wf/wf-rapid-1.5.0.js 50.992709ms 8245b
2.365114235s	200 https://s.yimg.com/aaq/wf/wf-tooltip-1.0.14.js 51.557366ms 9200b
2.378143092s	200 https://s.yimg.com/aaq/wf/wf-tabs-1.10.9.js 53.880444ms 13849b
2.387407536s	200 https://s.yimg.com/aaq/wf/wf-autocomplete-1.20.0.js 54.939662ms 26763b
2.394812491s	200 https://s.yimg.com/aaq/wf/wf-geolocation-1.2.9.js 55.947751ms 8008b
2.397307602s	200 https://s.yimg.com/aaq/wf/wf-menu-1.0.0.js 58.448394ms 5248b
2.414303472s	200 https://s.yimg.com/aaq/wf/wf-account-switch-1.1.2.js 49.645002ms 9897b
2.415895977s	200 https://s.yimg.com/os/yaft/yaft-0.3.27.min.js 50.90786ms 17701b
2.425078925s	200 https://s.yimg.com/os/yaft/yaft-plugin-aftnoad-0.1.5.min.js 61.08537ms 1144b
Downloaded assets: 83/83.
Total time: 2.425085058s.
Total size: 1962kb.
```
