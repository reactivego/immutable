package immutable_test

import (
	"os"
	"testing"

	"github.com/reactivego/immutable"
)

type Country struct {
	Code string
	Name string
}

var Countries []Country

func TestCountries(t *testing.T) {
	m := immutable.Map
	for k, v := range Countries {
		m = m.Set(k, v)
	}
	if len(Countries) != 175 {
		t.Errorf("len(Countries) got %d expected %d", len(Countries), 175)
	}
	if m.Len() != len(Countries) {
		t.Errorf("m.Len() got %d expected %d", m.Len(), len(Countries))
	}
	if m.Depth() != 2 {
		t.Errorf("m.Depth() got %d expected %d", m.Depth(), 2)
	}
	const arch = int(2 - uint64(^uint(0))>>63)
	if m.Size() != 10992/arch {
		t.Errorf("m.Size() got %d expected %d", m.Size(), 10992/arch)
	}
}

func BenchmarkGoMapGet(b *testing.B) {
	m := map[string]string{}
	func() {
		b.Helper()
		for _, c := range Countries {
			m[c.Name] = c.Code
		}
	}()
	for i := 0; i < b.N; i++ {
		c := Countries[i%len(Countries)]
		if _, present := m[c.Name]; !present {
			b.Errorf("_, present := m[%q]; present expected true got false", c.Name)
		}
	}
	b.ReportAllocs()
}

func BenchmarkImmutableMapGet(b *testing.B) {
	m := immutable.Map
	func() {
		b.Helper()
		for _, c := range Countries {
			m = m.Set(c.Name, c.Code)
		}
	}()
	for i := 0; i < b.N; i++ {
		c := Countries[i%len(Countries)]
		if !m.Has(c.Name) {
			b.Errorf("m.Has(%q) expected true got false", c.Name)
		}
	}
	b.ReportAllocs()
}

func BenchmarkGoMapSet(b *testing.B) {
	var m map[string]string
	count := len(Countries)
	for i := 0; i < b.N; i++ {
		if i%count == 0 {
			m = map[string]string{}
		}
		c := Countries[i%count]
		m[c.Name] = c.Code
	}
	b.ReportAllocs()
}

func BenchmarkImmutableMapSet(b *testing.B) {
	// var ss, es runtime.MemStats
	// runtime.ReadMemStats(&ss)
	// MapAlloc := uint64(0)
	m := immutable.Map
	count := len(Countries)
	for i := 0; i < b.N; i++ {
		if i%count == 0 {
			m = immutable.Map
		}
		c := Countries[i%count]
		m = m.Set(c.Name, c.Code)
		// MapAlloc += uint64(m.Size())
	}
	// runtime.ReadMemStats(&es)
	// b.ReportMetric(float64(MapAlloc)/float64(es.TotalAlloc-ss.TotalAlloc), "MapAlloc/TotalAlloc")
	// b.ReportMetric(float64(es.HeapAlloc-ss.HeapAlloc)/float64(b.N), "HeapAllocBytes/op")
	b.ReportAllocs()
}

func TestMain(m *testing.M) {
	var countries = map[string]string{
		"af": "Afghanistan",
		"al": "Albania",
		"dz": "Algeria",
		"ao": "Angola",
		"ai": "Anguilla",
		"ag": "Antigua and Barbuda",
		"ar": "Argentina",
		"am": "Armenia",
		"au": "Australia",
		"at": "Austria",
		"az": "Azerbaijan",
		"bs": "Bahamas",
		"bh": "Bahrain",
		"bb": "Barbados",
		"by": "Belarus",
		"be": "Belgium",
		"bz": "Belize",
		"bj": "Benin",
		"bm": "Bermuda",
		"bt": "Bhutan",
		"bo": "Bolivia",
		"ba": "Bosnia and Herzegovina",
		"bw": "Botswana",
		"br": "Brazil",
		"vg": "British Virgin Islands",
		"bn": "Brunei",
		"bg": "Bulgaria",
		"bf": "Burkina-Faso",
		"kh": "Cambodia",
		"cm": "Cameroon",
		"ca": "Canada",
		"cv": "Cape Verde",
		"ky": "Cayman Islands",
		"td": "Chad",
		"cl": "Chile",
		"cn": "China",
		"co": "Colombia",
		"cr": "Costa Rica",
		"ci": "Cote d'Ivoire",
		"hr": "Croatia",
		"cy": "Cyprus",
		"cz": "Czech Republic",
		"cd": "Democratic Republic of the Congo",
		"cg": "Republic of the Congo",
		"dk": "Denmark",
		"dm": "Dominica",
		"do": "Dominican Republic",
		"ec": "Ecuador",
		"eg": "Egypt",
		"sv": "El Salvador",
		"ee": "Estonia",
		"fm": "Federated States of Micronesia",
		"fj": "Fiji",
		"fi": "Finland",
		"fr": "France",
		"ga": "Gabon",
		"gm": "Gambia",
		"de": "Germany",
		"ge": "Georgia",
		"gh": "Ghana",
		"gb": "Great Britain",
		"gr": "Greece",
		"gd": "Grenada",
		"gt": "Guatemala",
		"gw": "Guinea Bissau",
		"gy": "Guyana",
		"hn": "Honduras",
		"hk": "Hong Kong",
		"hu": "Hungaria",
		"is": "Iceland",
		"in": "India",
		"id": "Indonesia",
		"iq": "Iraq",
		"ie": "Ireland",
		"il": "Israel",
		"it": "Italy",
		"jm": "Jamaica",
		"jp": "Japan",
		"jo": "Jordan",
		"kz": "Kazakhstan",
		"xk": "Kosovo",
		"ke": "Kenya",
		"kg": "Krygyzstan",
		"kw": "Kuwait",
		"la": "Laos",
		"lv": "Latvia",
		"lb": "Lebanon",
		"lr": "Liberia",
		"ly": "Libya",
		"lt": "Lithuania",
		"lu": "Luxembourg",
		"mo": "Macau",
		"mk": "Macedonia",
		"mg": "Madagascar",
		"mw": "Malawi",
		"my": "Malaysia",
		"ml": "Mali",
		"mv": "Maldives",
		"mt": "Malta",
		"mr": "Mauritania",
		"mu": "Mauritius",
		"mx": "Mexico",
		"md": "Moldova",
		"mn": "Mongolia",
		"me": "Montenegro",
		"ms": "Montserrat",
		"ma": "Morocco",
		"mz": "Mozambique",
		"mm": "Myanmar",
		"nr": "Nauru",
		"na": "Namibia",
		"np": "Nepal",
		"nl": "Netherlands",
		"nz": "New Zealand",
		"ni": "Nicaragua",
		"ne": "Niger",
		"ng": "Nigeria",
		"no": "Norway",
		"om": "Oman",
		"pk": "Pakistan",
		"pw": "Palau",
		"pa": "Panama",
		"pg": "Papua New Guinea",
		"py": "Paraguay",
		"pe": "Peru",
		"ph": "Philippines",
		"pl": "Poland",
		"pt": "Portugal",
		"qa": "Qatar",
		"tt": "Republic of Trinidad and Tobago",
		"ro": "Romania",
		"ru": "Russia",
		"rw": "Rwanda",
		"kn": "Saint Kitts and Nevis",
		"lc": "Saint Lucia",
		"vc": "Saint Vincent and the Grenadines",
		"st": "São Tomé e Príncipe",
		"sa": "Saudi Arabia",
		"sn": "Senegal",
		"rs": "Serbia",
		"sc": "Seychelles",
		"sl": "Sierra Leone",
		"sg": "Singapore",
		"sk": "Slovakia",
		"si": "Slovenia",
		"sb": "Soloman Islands",
		"za": "South Africa",
		"kr": "South Korea",
		"es": "Spain",
		"lk": "Sri Lanka",
		"sr": "Suriname",
		"sz": "Swaziland",
		"se": "Sweden",
		"ch": "Switzerland",
		"tw": "Taiwan",
		"tj": "Tajikistan",
		"tz": "Tanzania",
		"th": "Thailand",
		"to": "Tonga",
		"tn": "Tunisia",
		"tr": "Turkey",
		"tm": "Turkmenistan",
		"tc": "Turks and Caicos Islands",
		"ug": "Uganda",
		"ua": "Ukraine",
		"ae": "United Arab Emirates",
		"us": "United States of America",
		"uy": "Uruguay",
		"uz": "Uzbekistan",
		"vu": "Vanatu",
		"ve": "Venezuela",
		"vn": "Vietnam",
		"ye": "Yemen",
		"zm": "Zambia",
		"zw": "Zimbabwe",
	}
	Countries = make([]Country, 0, len(countries))
	for k, v := range countries {
		Countries = append(Countries, Country{Code: k, Name: v})
	}
	code := m.Run()
	os.Exit(code)
}
